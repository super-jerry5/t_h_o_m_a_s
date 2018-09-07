package main

import (
	"thomas/entity"
	"thomas/common"
	"time"
	"github.com/astaxie/beego/logs"
)

func FindDoingTask(worker string,  mapJob map[int]*entity.JobInfo)(map[int]*entity.JobInfo, error) {
	mapDoingJob := make(map[int]*entity.JobInfo)
	now := time.Now().Unix()
	for id, job := range mapJob {
		if (job.State == common.JOB_STATE_DOING && job.Worker == worker && 
			(job.UpdateTime.Unix() + g_max_update_interval_s) >= now) {
			mapDoingJob[id] = job
			logs.Debug("FindDoingTask, id:%d, now:%d, update_time:%d, job:%v", id, now, job.UpdateTime.Unix(), job)
		}
	}
	return mapDoingJob, nil
}

func FindStartTask(worker string,  mapJob map[int]*entity.JobInfo)(map[int]*entity.JobInfo, error) {
	mapStartJob := make(map[int]*entity.JobInfo)
	now := time.Now().Unix()
	for id, job := range mapJob {
		if ((job.State == common.JOB_STATE_TODO && (job.Worker == "" ||
			(job.Worker != "" && (job.UpdateTime.Unix() + g_max_update_interval_s) < now))) ||
			(job.State == common.JOB_STATE_DOING && (job.UpdateTime.Unix() + g_max_update_interval_s) < now)) {
			mapStartJob[id] = job
			logs.Debug("FindStartTask, id:%d, now:%d, update_time:%d, job:%v", id, now, job.UpdateTime.Unix(), job)
		}
	}
	return mapStartJob, nil
}

func FindStopTask(worker string,  mapJob map[int]*entity.JobInfo)(map[int]*entity.JobInfo, error) {
	mapStopJob := make(map[int]*entity.JobInfo)
	for id, job := range mapJob {
		if job.State == common.JOB_STATE_TOSTOP {
				mapStopJob[id] = job
				logs.Debug("FindStopTask, id:%d, update_time:%d, job:%v", id, job.UpdateTime.Unix(), job)
			}	
	}
	
	return mapStopJob, nil
}

func FindRestartTask(worker string,  mapJob map[int]*entity.JobInfo)(map[int]*entity.JobInfo, error) {
	mapRestartJob := make(map[int]*entity.JobInfo)
	//now := time.Now().Unix()
	for id, job := range mapJob {
		if job.State == common.JOB_STATE_TORESTART && (job.Worker == "" ||
			(job.Worker == worker)) {
			mapRestartJob[id] = job
			logs.Debug("FindRestartTask, id:%d, update_time:%d, job:%v", id, job.UpdateTime.Unix(), job)
		}
	}
	return mapRestartJob, nil
}

func FindJobList()(map[int]*entity.JobInfo, error) {
	var name = ""
	jobInfo := &entity.JobInfo{Name:name}
	jobs, err := jobInfo.FindAllJobInfo()
	if err != nil {
		return nil, err
	}
	
	mapJob := make(map[int]*entity.JobInfo)
	for _, job := range jobs {
		mapJob[job.Id] = job
	}
	return mapJob, nil
}

func StartDb(worker string,  job *entity.JobInfo) (error) {
	job.Worker = worker
	now := time.Now()
	job.ModifyTime = now
	job.UpdateTime = now
	job.Op = common.OP_TYPE_NONE
	job.State = common.JOB_STATE_DOING
	logs.Info("StartDb id:%d, ModifyTime:%v, UpdateTime:%v, op:%d, state:%d",
				job.Id, job.ModifyTime, job.UpdateTime, job.Op, job.State)
	return job.UpdateStateAndOp()
}

func StopDb(worker string,  job *entity.JobInfo) (error) {
	job.Worker = ""
	now := time.Now()
	job.ModifyTime = now
	job.UpdateTime = now
	job.Op = common.OP_TYPE_NONE
	job.State = common.JOB_STATE_STOP
	logs.Info("StopDb id:%d, ModifyTime:%v, UpdateTime:%v, op:%d, state:%d",
				job.Id, job.ModifyTime, job.UpdateTime, job.Op, job.State)
	return job.UpdateStateAndOp()
}

func RestartDb(worker string,  job *entity.JobInfo) (error) {
	job.Worker = worker
	now := time.Now()
	job.ModifyTime = now
	job.UpdateTime = now
	job.Op = common.OP_TYPE_NONE
	job.State = common.JOB_STATE_DOING
	logs.Info("RestartDb id:%d, ModifyTime:%v, UpdateTime:%v, op:%d, state:%d",
				job.Id, job.ModifyTime, job.UpdateTime, job.Op, job.State)
	return job.UpdateStateAndOp()
}

func UpdateDb(worker string,  job *entity.JobInfo) (error) {
	now := time.Now()
	job.ModifyTime = now
	job.UpdateTime = now
	logs.Info("UpdateDb id:%d, ModifyTime:%v, UpdateTime:%v, op:%d, state:%d",
				job.Id, job.ModifyTime, job.UpdateTime, job.Op, job.State)
	return job.UpdateStateAndOp()
}
