package main

import (
	//"github.com/astaxie/beego"
	"os/exec"
	"fmt"
	"time"
	"bytes"
	"bufio"
	"io"
	"thomas/entity"
	"strconv"
	"strings"
	"github.com/astaxie/beego/logs"
)

type FfmpegInfo struct {
	TaskId int
	Pid int
	Param string
}

func checkErr(err error) {
    if err != nil {
        logs.Error(err)
        panic(err)  
    }
}

func get_cmd(command string) *exec.Cmd {
	return exec.Command("/bin/bash", "-c", command)
}
//不需要执行命令的结果与成功与否，执行命令马上就返回
func exec_shell_no_result(command string)(error) {  
    cmd := get_cmd(command)
    //开始执行c包含的命令，但并不会等待该命令完成即返回
    err := cmd.Start()
    if err != nil {
        logs.Error("%v: exec command:%v error:%v\n", time.Now(), command, err)
        return err
    }
    fmt.Printf("Waiting for command:%v to finish...\n", command)
//    //阻塞等待fork出的子进程执行的结果，和cmd.Start()配合使用[不等待回收资源，会导致fork出执行shell命令的子进程变为僵尸进程]
//    err = cmd.Wait()
//    if err != nil {
//        logs.Error("%v: Command finished with error: %v\n", time.Now(), err)
//        return err
//    }
    return nil
}

//阻塞式的执行外部shell命令的函数,等待执行完毕并返回标准输出
func exec_shell_string(command string) (string, error) {
    cmd := get_cmd(command)

    //读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
    var out bytes.Buffer
    cmd.Stdout = &out

    //Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
    err := cmd.Run()
    checkErr(err)


    return out.String(), err
}

func exec_shell_lines(command string) ([]string, error) {
    cmd := get_cmd(command)
    
    var lines []string = nil  
    //显示运行的命令
    logs.Info(cmd.Args)
    //StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
    stdout, err := cmd.StdoutPipe()

    if err != nil {
        logs.Error(err)
        return nil, err
    }

    err3 := cmd.Start()
    if err3 != nil {
        logs.Error(err3)
        return nil, err3
    }
    //创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
    reader := bufio.NewReader(stdout)

    //实时循环读取输出流中的一行内容
    for {
        line, err2 := reader.ReadString('\n')
        if err2 != nil || io.EOF == err2 {
            break
        }
        lines = append(lines, line)
    }

    //阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
    err4 := cmd.Wait()
    return lines, err4
}

func MakeFfmpegCmd(job *entity.JobInfo)(string, error) {
	var cmd string = "ffmpeg "
	cmd = cmd + "-identifier " + strconv.Itoa(job.Id) + " "
	cmd = cmd + "-i " + job.InputUrl + " "
	cmd = cmd + "-s " + job.Scale + " "
	cmd = cmd + "-acodec copy "
	var vcodec string = ""
	if job.VcodecType == 0 {
		vcodec = "h264_nvenc"
	} else {
		vcodec = "h265_nvenc"
	}
	cmd = cmd + "-vcodec " + vcodec + " "
	cmd = cmd + "-gpu " + strconv.Itoa(job.GpuNumber) + " "
	cmd = cmd + "-b " + strconv.Itoa(job.VcodecBitrate) + " "
	cmd = cmd + "-f flv " + job.OutputUrl
	
	//TODO: remove the debug cmd
	//cmd = "ffmpeg -identifier " + strconv.Itoa(job.Id) + " -re -i ~/Movies/chiji_7.5M_repeat10.mp4 -s 960x540 -acodec copy -vcodec libx264 -b 60000 -f flv rtmp://10.86.0.101:1935/live/720p_540p"
	logs.Info("MakeFfmpegCmd succ, cmd:%s", cmd)
	return cmd, nil
}

func ExecFfmpegCmd(cmd string)(error) {
	return exec_shell_no_result(cmd)
}

func RestartFfmpeg(ff FfmpegInfo, job *entity.JobInfo)(error) {
	err := StopFfmpeg(ff)
	if err != nil {
		return err
	}
	
	err = StartFfmpeg(job)
	if err != nil {
		return err
	}
	
	return err
}

func StartFfmpeg(job *entity.JobInfo)(error) {	
	cmd, err := MakeFfmpegCmd(job)
	if err != nil {
		logs.Error("MakeFfmpegCmd failed, job_id:%d, job_name:%s, err:%s",
			job.Id, job.Name, err)
		return err
	}
	
	err = ExecFfmpegCmd(cmd)
	if err != nil {
		logs.Error("ExecFfmpegCmd failed, job_id:%d, job_name:%s, cmd:%s, err:%s",
			job.Id, job.Name, cmd, err)
		return err
	} 
	
	logs.Info("StartFfmpeg succ. id:%d, cmd:%s, job:%v", job.Id, cmd, job)
	return err
}


func StopFfmpeg(ff FfmpegInfo)(error) {
	var cmd string = "kill -9 " + strconv.Itoa(ff.Pid)
	err := exec_shell_no_result(cmd)
	if err != nil {
		logs.Error("StopFfmpegCmd failed, ff_taskid:%d, ff_pid:%d, ff_cmd:%s, err:%s",
			ff.TaskId, ff.Pid, ff.Param, err)
	}	
	return err
}



func GetFfmpegInfoList()(map[int]FfmpegInfo, error) {
	 var ffmpeg_cmd_prefix string = "ffmpeg -identifier"
	var find_ffmpeg_cmd string = "ps -ef | grep '" + ffmpeg_cmd_prefix +"'"
	outlines, err := exec_shell_lines(find_ffmpeg_cmd)
	if err != nil {
		logs.Error("exec_shell_lines failed, err:%s, find_ffmpeg_cmd:%s", err, find_ffmpeg_cmd)
		return nil, err
	} else {
		logs.Debug("exec_shell_lines succ, find_ffmpeg_cmd:%s, outlines:%s", find_ffmpeg_cmd, outlines)
	}

    cmd_map := make(map[int]FfmpegInfo)
    for idx, line := range outlines {
		if strings.Contains(line, "grep") {
			logs.Debug("%dth, line:%s contains grep, ", idx, line)
			continue
		}
		ffmpeg_begin_index := strings.Index(line, ffmpeg_cmd_prefix)
		if ffmpeg_begin_index == -1 {
			continue
		}
		all_fields := strings.Fields(line)
		pid, err2 := strconv.Atoi(all_fields[1])
		if err2 != nil {
			logs.Error("%dth pid invalid, pid_fields[%s], %v, line:%s, err2:%s", idx, all_fields[1], all_fields, line, err2.Error())
			continue
		}
		ffmpeg_cmd := line[ffmpeg_begin_index:]
		ffmpeg_fields := strings.Fields(ffmpeg_cmd)
		task_id,err3 := strconv.Atoi(ffmpeg_fields[2])
		if err3 != nil {
			logs.Error("%dth task_id invalid, task_fields[%s], %v, line:%s, err3:%s", idx, ffmpeg_fields[2], all_fields, line, err3.Error())
			continue
		}
		logs.Debug("get ff task succ, %dth, task_id:%d, pid:%d, ffmpeg_cmd:%s", idx, task_id, pid, ffmpeg_cmd)
		ffmpegInfo := FfmpegInfo{task_id, pid, ffmpeg_cmd}
		cmd_map[task_id] = ffmpegInfo
		
		logs.Debug("%dth, task_id:%d, pid:%d, ffmpeg_cmd:%s, line:%v", 
			idx, task_id, pid, ffmpeg_cmd, line)
	}
	
	return cmd_map, nil
}
	
	
//	exec_shell_no_result("ps -ef | grep ffmpeg -identifier")
//	out, err := exec_shell_string("ps -ef | grep Chrome | awk '{print $9}'")
//	if err != nil {
//	logs.Errorn("exec_shell_string failed, err:", err)
//	} else {
//		fmt.Println("exec_shell_string succ, out:%s", out)
//	}
//	
//	outlines, err := exec_shell_lines("ps -ef | grep Chrome | awk '{print $9}'")
//	if err != nil {
//		fmt.Println("exec_shell_lines failed, err:", err)
//	} else {
//		fmt.Println("exec_shell_lines succ, outlines:%s", outlines)
//	}

