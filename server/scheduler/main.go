package main

import (
	//"github.com/astaxie/beego"
	"thomas/common"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"thomas/entity"
	"time"
	"runtime"
	"github.com/astaxie/beego/logs"
	config "github.com/astaxie/beego/config"
)

var g_max_update_interval_s int64 = 30

func init() {
	ini, err := config.NewConfig("ini", "conf/scheduler.conf")
	if err != nil {
		panic(err)
	}
	mysql_string := ini.String("mysql_string")
	max_update_interval_s, err := ini.Int64("max_update_interval_s")
	if err != nil  {
		panic(err)
	}
	g_max_update_interval_s = max_update_interval_s
        logs.SetLogger(logs.AdapterFile,`{"filename":"/data/logs/scheduler/scheduler.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	logs.SetLogFuncCall(true)
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysql_string)
	orm.RegisterModel(&entity.JobInfo{})
}


func ProcessStop(worker string, mapJobSelf map[int]*entity.JobInfo, mapJobDb map[int]*entity.JobInfo, mapFf map[int]FfmpegInfo)(error) {
	var retErr error = nil
	mapStopJob, err := FindStopTask(worker, mapJobDb)
	if err != nil {
		logs.Error("FindStopTask failed, err:", err)
		return err
	} else {
		logs.Debug("FindStopTask succ, mapStopJob_size:%d", len(mapStopJob))
	}
	for id, job := range mapStopJob {
		if ff, ok := mapFf[id]; ok {
			err := StopFfmpeg(ff)
			if err != nil {
				logs.Error("StopFfmpeg failed, id:%d, err:%s", id, err)
				retErr = err
				continue
			}
			delete(mapFf, id)
		}
		
		if _, ok := mapJobSelf[id]; ok {
			delete(mapJobSelf, id)
		}
		
		err := StopDb(worker, job)
		if err != nil {
			logs.Error("StopDb failed, id:%d, errL%s", id, err)
			retErr = err
			continue
		}
	}
	return retErr
}

func ProcessRestart(worker string, mapJobSelf map[int]*entity.JobInfo, mapJobDb map[int]*entity.JobInfo, mapFf map[int]FfmpegInfo)(error) {
	var retErr error = nil
	mapRestartJob, err := FindRestartTask(worker, mapJobDb)
	if err != nil {
		logs.Error("FindRestartTask failed, err:", err)
		return err
	} else {
		logs.Debug("FindRestartTask succ, mapRestartJob_size:%d", len(mapRestartJob))
	}
	
	
	for id, job := range mapRestartJob {
		if ff, ok := mapFf[id]; ok {
			err := RestartFfmpeg(ff, job)
			if err != nil {
				logs.Error("RestartFfmpeg failed, id:%d, err:%s", id, err)
				retErr = err
				continue
			}
		}
		
		if _, ok := mapJobSelf[id]; ok {
			delete(mapJobSelf, id)
		}
		mapJobSelf[id] = job
		
		err := RestartDb(worker, job)
		if err != nil {
			logs.Error("RestartDb failed, id:%d, errL%s", id, err)
			retErr = err
			continue
		}
	}
	
	return retErr
}

func ProcessStart(worker string, mapJobSelf map[int]*entity.JobInfo, mapJobDb map[int]*entity.JobInfo, mapFf map[int]FfmpegInfo)(error) {
	var retErr error = nil
	mapStartJob, err := FindStartTask(worker, mapJobDb)
	if err != nil {
		logs.Error("FindStartTask failed, err:", err)
		return err
	} else {
		logs.Debug("FindStartTask succ, mapStartJob:%d", len(mapStartJob))
	}
	
	
	for id, job := range mapStartJob {
		
		if ff, ok := mapFf[id]; ok {
			err := StopFfmpeg(ff)
			if err != nil {
				logs.Error("ProcessStart StopFfmpeg failed, id:%d, err:%s", id, err)
				retErr = err
				continue
			}
			delete(mapFf, id)
			
		}
		
		err := StartFfmpeg(job)
		if err != nil {
			logs.Error("StartFfmpeg failed, id:%d, err:%s", id, err)
			retErr = err
			continue
		}
			
		if _, ok := mapJobSelf[id]; ok {
			delete(mapJobSelf, id)
		}
		mapJobSelf[id] = job
		
		err = StartDb(worker, job)
		if err != nil {
			logs.Error("StartDb failed, id:%d, errL%s", id, err)
			retErr = err
			continue
		}
	}
	
	return retErr
}

func ProcessCheckDoing(worker string, mapJobSelf map[int]*entity.JobInfo, mapJobDb map[int]*entity.JobInfo, mapFf map[int]FfmpegInfo)(error) {
	var retErr error = nil
	time.Sleep(1000 * time.Millisecond)
	// in order to make sure the follow updateDb will succ(for the reason the update change at least)
	
	mapDoingJob, err := FindDoingTask(worker, mapJobDb)
	if err != nil {
		logs.Error("FindDoingTask failed, err:", err)
		return err
	} else {
		logs.Debug("FindDoingTask succ, mapDoingJob:%d", len(mapDoingJob))
	}
	
	for id, job := range mapDoingJob {
		if _, ok := mapJobSelf[id]; ok {
			continue
		}
		mapJobSelf[id] = job
	}
	
// step1. 确认目前所有的本地任务的状态是否和远端任务一致。
//        如果本地存在，远端不存在，或者远端状态非法，或者远端的work非法，则删除本地任务，且结束任务
// step2. 如果远端存在，但本地不存在，且没有过期。启动任务，且添加到本地任务
// step3. 如果本地任务列表和实际任务列表不一致。包括本地列表有，实际没有，则启动；本地没有，实际有，则停止；

    now := time.Now().Unix()
    for id, job := range mapJobSelf {
    	//本地任务存在，远端任务不存在，删除本地任务
	    if _, ok := mapJobDb[id]; !ok {
			delete(mapJobSelf, id)
			logs.Info("find no job in Db, delete local, id:%d, worker:%s", id, worker)
			continue
		}
	    //本地任务存在，但远端任务信息非法
	    if (mapJobDb[id].State != common.JOB_STATE_DOING || 
		    ((mapJobDb[id].UpdateTime.Unix() + g_max_update_interval_s) < now) ||
		    mapJobDb[id].Worker != worker) {
	    	logs.Info("job in Db invalid, delete local, id:%d, worker:%s, now:%d, job:%v", id, worker, now, mapJobDb[id])
			delete(mapJobSelf, id)
			continue
	    }
	    
	    if _, ok := mapFf[id]; ok {
			logs.Debug("id:%d, check succ, pid:%d, cmd:%s", id, mapFf[id].Pid, mapFf[id].Param)
		} else {
			err := StartFfmpeg(job)
			if err != nil {
				logs.Error("id:%d check failed, StartFfmpeg failed, err:%s", id, err)
				retErr = err
				continue
			}
			logs.Debug("id:%d, check succ, start it again ", id)
		}
	    
	    err := UpdateDb(worker, job)
	    if err != nil {
			logs.Error("UpdateDb failed, id:%d, errL%s", id, err)
			retErr = err
			continue
		}
    }
	
	for id, ff := range mapFf {
		if _, ok := mapJobSelf[id]; ok {
			continue
		}
		err := StopFfmpeg(ff)
		if err != nil {
			logs.Error("StopFfmpeg failed, id:%d, ff:%v,  errL%s", id, ff, err)
			retErr = err
			continue
		}
		
	}
	
	return retErr	
}

func GetList()(map[int]FfmpegInfo, map[int]*entity.JobInfo, error) {
	mapJob, err := FindJobList()
	if err != nil {
		logs.Error("FindJobList failed, err:", err)
		return nil, nil,err
	} else {
		logs.Debug("FindAllJobInfo succ, jobs_size:%d", len(mapJob))
	}
	
	mapFf, err := GetFfmpegInfoList()
	if err != nil {
		logs.Error("GetFfmpegInfoList failed, err:", err)
		return nil,nil,err
	} else {
		logs.Debug("GetFfmpegInfoList succ, ff_size:%d", len(mapFf))
	}
	return mapFf, mapJob, nil
}

func scheduler(mapJobSelf map[int]*entity.JobInfo, worker string) (map[int]*entity.JobInfo, error){
	logs.Debug("scheduler begin, mapJobSelf_size:%d, worker:%s", len(mapJobSelf), worker)
	
	mapFf, mapJob, err := GetList()
	if err != nil {
		logs.Error("GetList before Stop failed, err:", err)
		return mapJobSelf, err
	}
	err = ProcessStop(worker, mapJobSelf, mapJob, mapFf)
	if err != nil {
		logs.Error("ProcessStop failed, err:", err)
		return mapJobSelf,err
	} else {
		logs.Debug("ProcessStop succ, mapJobSelf_size:%d, worker:%s", len(mapJobSelf), worker)
	}
	
	err = ProcessRestart(worker, mapJobSelf, mapJob, mapFf)
	if err != nil {
		logs.Error("ProcessRetart failed, err:", err)
		return mapJobSelf,err
	} else {
		logs.Debug("ProcessRetart succ, mapJobSelf_size:%d, worker:%s", len(mapJobSelf), worker)
	}
	
	err = ProcessStart(worker, mapJobSelf, mapJob, mapFf)
	if err != nil {
		logs.Error("ProcessStart failed, err:", err)
		return mapJobSelf,err
	} else {
		logs.Debug("ProcessStart succ, mapJobSelf_size:%d, worker:%s", len(mapJobSelf), worker)
	}
	
	mapFf, mapJob, err = GetList()
	if err != nil {
		logs.Error("GetList Before Doing failed, err:", err)
		return mapJobSelf, err
	}
	
	err = ProcessCheckDoing(worker, mapJobSelf, mapJob, mapFf)
	if err != nil {
		logs.Error("ProcessStart failed, err:", err)
		return mapJobSelf,err
	} else {
		logs.Debug("ProcessStart succ, mapJobSelf_size:%d, worker:%s", len(mapJobSelf), worker)
	}
	
		
	
	logs.Debug("scheduler end, mapJobSelf_size:%d, worker:%s", len(mapJobSelf), worker)
	
	return mapJobSelf, err
	
}

func main()  {

	logs.Debug("Thomas Live Transcode System Scheduler start")
	runtime.GOMAXPROCS(runtime.NumCPU())
	orm.Debug = true
	
	
	mapJob := make(map[int]*entity.JobInfo)
	var worker string = "this machine"
	sum := 0
	for {
		sum++
//		if sum > 1 {
//			break
//		}
		
		mapJob, err := scheduler(mapJob, worker)
		if err != nil {
			logs.Error("scheduler failed, %dth err:", sum, err)
		} else {
			logs.Debug("scheduler succ, %dth jobsize_%d:", sum, len(mapJob))
		}
		time.Sleep(10000 * time.Millisecond)	
	}

}

