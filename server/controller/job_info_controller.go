package controller

import (
	"thomas/entity"
	"thomas/common"
	"time"
	"github.com/astaxie/beego/logs"
)

// job info
type JobInfoManagerController struct  {
	BaseController

}

func (this *JobInfoManagerController) List() {
	name := this.GetString("Name")
	jobInfo := &entity.JobInfo{Name:name}
	jobs,err := jobInfo.FindAllJobInfoByPage()
	common.PanicIf(err)
	this.Data["jobs"] = jobs
	this.Data["name"] = name
	this.TplName = "jobinfo/list.html"
	this.Render()
}

func (this *JobInfoManagerController) ToAdd()  {

	this.TplName = "jobinfo/add.html"
	this.Render()
}

func (this *JobInfoManagerController) Add()  {
	logs.Debug("Add Begin")
	jsonResult := &entity.JsonResult{}
	var worker = ""
	var op = common.OP_TYPE_START
	var state = common.JOB_STATE_TODO
	name := this.GetString("Name")
	input_url := this.GetString("InputUrl")
	input_url_backup := this.GetString("InputUrlBackup")
	output_url := this.GetString("OutputUrl")
	scale := this.GetString("Scale")
	vcodec_type, err := this.GetInt("VcodecType")
	if err != nil {
		jsonResult.Message = "vcodec_type invalid"
	}
	vcodec_bitrate, err := this.GetInt("VcodecBitrate")
	if err != nil {
		jsonResult.Message = "vcodec_bitrate invalid"
	}
	watermark_text := this.GetString("WatermarkText")
	watermark_param := this.GetString("WatermarkParam")
	audio_param := this.GetString("AudioParam")
	gpu_number, err := this.GetInt("GpuNumber")
	if err != nil {
		jsonResult.Message = "gpu number invalid"
	}
	gpu_deint, err := this.GetInt("GpuDeint")
	if err != nil {
		jsonResult.Message = "gpu deint invalid"
	}
	
	//worker := this.GetString("Worker")
	//op := this.GetInt("Op")
	//state := this.GetInt("State")
	//score := this.GetInt("Score")
	
	if name == ""{
		jsonResult.Message ="name couldn't be empty"
	} else if  input_url == "" {
		jsonResult.Message ="input url couldn't be empty"
	}  else if output_url == "" {
		jsonResult.Message ="input url couldn't be empty"
	} else if vcodec_bitrate == 0 {
		jsonResult.Message ="vcodec bitrate couldn't be zero"
	} else if scale == "" {
		jsonResult.Message = "scale info couldn't be empty"
	}
	if jsonResult.Message == ""{
		jobInfo := &entity.JobInfo{}
		jobInfo.CreateTime = time.Now()
		jobInfo.UpdateTime = jobInfo.CreateTime
		jobInfo.ModifyTime = jobInfo.CreateTime
		jobInfo.Name = name
		jobInfo.InputUrl = input_url
		jobInfo.InputUrlBackup = input_url_backup
		jobInfo.OutputUrl = output_url
		jobInfo.InputUrl = input_url
		jobInfo.Scale    = scale
		jobInfo.VcodecType = vcodec_type
		jobInfo.VcodecBitrate = vcodec_bitrate
		jobInfo.WatermarkText = watermark_text
		jobInfo.WatermarkParam = watermark_param
		jobInfo.AudioParam = audio_param
		jobInfo.GpuNumber = gpu_number
		jobInfo.GpuDeint = gpu_deint
		jobInfo.Worker = worker
		jobInfo.Op = op
		jobInfo.State = state
		
		err := jobInfo.SaveJobInfo()
		if err != nil {
			jsonResult.Message = "save jobinfo failed"
		} else {

			jsonResult.Message = "succ"
			jsonResult.Success = true
		}
	}
	logs.Debug("Add End", jsonResult)
	this.Data["json"] = jsonResult
	this.ServeJSON()
}

func (this *JobInfoManagerController)ToEdit()  {

	id,err := this.GetInt("Id")
	common.PanicIf(err)

	jobInfo := &entity.JobInfo{Id:id,}
	err = jobInfo.GetJobInfoById()
	common.PanicIf(err)
	this.Data["jobInfo"] = jobInfo
	this.TplName = "jobinfo/edit.html"
	this.Render()
}

func (this *JobInfoManagerController)Info()  {

	logs.Debug("Info Begin")
	id,err := this.GetInt("Id")
	common.PanicIf(err)

	jobInfo := &entity.JobInfo{Id:id,}
	err = jobInfo.GetJobInfoById()
	common.PanicIf(err)
	
	logs.Debug("Info End")
	this.Data["jobInfo"] = jobInfo
	this.TplName = "jobinfo/info.html"
	this.Render()
}
func (this *JobInfoManagerController) Edit()  {
	jsonResult := &entity.JsonResult{}
	for {
		id,err:= this.GetInt("Id")
		name := this.GetString("Name")
		input_url := this.GetString("InputUrl")
		input_url_backup := this.GetString("InputUrlBackup")
		output_url := this.GetString("OutputUrl")
		scale := this.GetString("Scale")
		vcodec_type, err := this.GetInt("VcodecType")
		if err != nil {
			jsonResult.Message = "vcodec_type invalid"
		}
		vcodec_bitrate, err := this.GetInt("VcodecBitrate")
		if err != nil {
			jsonResult.Message = "vcodec_bitrate invalid"
		}
		watermark_text := this.GetString("WatermarkText")
		watermark_param := this.GetString("WatermarkParam")
		audio_param := this.GetString("AudioParam")
		gpu_number, err := this.GetInt("GpuNumber")
		if err != nil {
			jsonResult.Message = "gpu number invalid"
		}
		gpu_deint, err := this.GetInt("GpuDeint")
		if err != nil {
			jsonResult.Message = "gpu deint invalid"
		}
		
		if err != nil {
			jsonResult.Message = "job state invalid"
		}
		
		if name == ""{
			jsonResult.Message ="job name couldn't be empty"
		}  else if err != nil  {
			jsonResult.Message = "record doesn't exist"
		}else if  input_url == "" {
			jsonResult.Message ="input url couldn't be empty"
		}  else if output_url == "" {
			jsonResult.Message ="input url couldn't be empty"
		} else if vcodec_bitrate == 0 {
			jsonResult.Message ="vcodec bitrate couldn't be zero"
		} else if scale == "" {
			jsonResult.Message = "scale info couldn't be empty"
		}
	
	
		jobInfoOld := &entity.JobInfo{}
		jobInfoOld.Id = id
		err = jobInfoOld.GetJobInfoById()
		if err != nil {
			jsonResult.Message = "get jobInfo old failed"
		}
		if jobInfoOld.State != common.JOB_STATE_DOING && jobInfoOld.State != common.JOB_STATE_STOP {
			jsonResult.Message = "job state is not in(TRANSCODING, STOPED), couldn't stop"
		}
		
		if jsonResult.Message != ""{
			logs.Error("Edit failed, id:%d, msg:%s", id, jsonResult.Message)
			break
		}
		
		jobInfo := &entity.JobInfo{}
		//jobInfo.CreateTime = time.Now()
		//jobInfo.UpdateTime = jobInfo.CreateTime
		jobInfo.Id = id
		jobInfo.ModifyTime = time.Now()
		jobInfo.Name = name
		jobInfo.InputUrl = input_url
		jobInfo.InputUrlBackup = input_url_backup
		jobInfo.OutputUrl = output_url
		jobInfo.Scale    = scale
		jobInfo.VcodecType = vcodec_type
		jobInfo.VcodecBitrate = vcodec_bitrate
		jobInfo.WatermarkText = watermark_text
		jobInfo.WatermarkParam = watermark_param
		jobInfo.AudioParam = audio_param
		jobInfo.GpuNumber = gpu_number
		jobInfo.GpuDeint = gpu_deint
		
		jobInfo.State = jobInfoOld.State
		jobInfo.Op = jobInfoOld.Op
			
		if jobInfoOld.InputUrl != jobInfo.InputUrl ||
			jobInfoOld.InputUrlBackup != jobInfo.InputUrlBackup ||
			jobInfoOld.OutputUrl != jobInfo.OutputUrl ||
			jobInfoOld.Scale != jobInfo.Scale ||
			jobInfoOld.VcodecType != jobInfo.VcodecType ||
			jobInfoOld.VcodecBitrate != jobInfo.VcodecBitrate ||
			jobInfoOld.WatermarkText != jobInfo.WatermarkText ||
			jobInfoOld.WatermarkParam != jobInfo.WatermarkParam ||
			jobInfoOld.AudioParam != jobInfo.AudioParam ||
			jobInfoOld.GpuNumber != jobInfo.GpuNumber ||
			jobInfoOld.GpuDeint != jobInfo.GpuDeint {
				
			if jobInfoOld.State == common.JOB_STATE_DOING {
				jobInfo.State = common.JOB_STATE_TORESTART
				jobInfo.Op = common.OP_TYPE_RESTART
			}
		}
		
		err = jobInfo.UpdateJobInfo()
		if err != nil {
			jsonResult.Message = err.Error()
		} else {
			jsonResult.Message = "update jobinfo succ"
			jsonResult.Success = true
		}
		break
	}
	
	this.Data["json"] = jsonResult
	logs.Info("Edit end, result:%s", jsonResult)
	this.ServeJSON()

}



// delete jobinfo


func (this *JobInfoManagerController) Delete() {
	jsonResult := &entity.JsonResult{}
	id,err := this.GetInt("Id")

	if err != nil {
		jsonResult.Message = "this record doesn't exist"
	} else {
		jobInfo := entity.JobInfo{Id:id}
		err = jobInfo.DeleteJobInfo()
		if err != nil {
			jsonResult.Message = "delete failed, pls retry"
		} else  {
			jsonResult.Message = "delete succ"
			jsonResult.Success = true
		}
	}
	this.Data["json"] = jsonResult
	this.ServeJSON()
}

func (this *JobInfoManagerController) Operate()  {
	logs.Debug("Operate Begin")
	jsonResult := &entity.JsonResult{}
	for {
		id,err := this.GetInt("Id")
		if err != nil {
			jsonResult.Message = "input param failed"
		}
		op,err := this.GetInt("Op")
		if err != nil {
			jsonResult.Message = "input param failed"
		} 
		logs.Debug("Operate Get Op&Id, op:%d, id:%d", op, id)
		
		jobInfoOld := &entity.JobInfo{}
		jobInfoOld.Id = id
		err = jobInfoOld.GetJobInfoById()
		if err != nil {
			jsonResult.Message = "get jobInfo old failed"
		}
			
		if jsonResult.Message != "" {
			break
		}
		
		if (op == common.OP_TYPE_START && jobInfoOld.State != common.JOB_STATE_STOP) ||
			(op == common.OP_TYPE_STOP && jobInfoOld.State == common.JOB_STATE_STOP) ||
			(op == common.OP_TYPE_RESTART && (jobInfoOld.State != common.JOB_STATE_DOING)) {
			logs.Error("Operate failed, id:%d, op:%d invalid, old state:%d", id, op, jobInfoOld.State)
			jsonResult.Message = "op invalid"
			break
		}
		logs.Debug("Operate xxx_1, result:%s", jsonResult)
		jobInfo := &entity.JobInfo{Id:id}
		err = jobInfo.GetJobInfoById()
		if err != nil {
			
			logs.Debug("Operate xxx_2, result:%s", jsonResult)
			jsonResult.Message = "this job record doesn't exist"
		} else {
			
			logs.Debug("Operate xxx_3, result:%s", jsonResult)
			if (op == common.OP_TYPE_START) {
				jobInfo.State = common.JOB_STATE_TODO
			} else if (op == common.OP_TYPE_STOP) {
				jobInfo.State = common.JOB_STATE_TOSTOP
			} else if (op == common.OP_TYPE_RESTART) {
				jobInfo.State = common.JOB_STATE_TORESTART
			}
			jobInfo.Op = op
			jobInfo.ModifyTime = time.Now()
			err = jobInfo.OperateJobInfo()
			if err != nil {
				logs.Debug("Operate xxx_4, result:%s", jsonResult)
				jsonResult.Message = "failed"
			} else {
				
				logs.Debug("Operate xxx_5, id:%d, op:%d, state:%d, result:%s", jobInfo.Id, jobInfo.Op, jobInfo.State, jsonResult)
				jsonResult.Message = "succ"
				jsonResult.Success = true
			}
		}
		break
	}

	logs.Debug("Operate End, result:%s", jsonResult)
	this.Data["json"] = jsonResult
	this.ServeJSON()
}




