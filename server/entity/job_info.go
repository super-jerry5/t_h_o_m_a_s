package entity

import (
	"time"
	"github.com/astaxie/beego/orm"
	"errors"
)

type JobInfo struct {
	Id                  int           `orm:"pk;auto"`
	Name                string        `orm:"notnull"`
	InputUrl            string        `orm:"notnull"`
	InputUrlBackup      string        `orm:"notnull"`
	OutputUrl           string        `orm:"notnull"`
	Scale               string        `orm:"notnull"`
	VcodecType          int           `orm:"notnull"`
	VcodecBitrate       int           `orm:"notnull"`
	WatermarkText       string        `orm:"notnull"`
	WatermarkParam      string        `orm:"notnull"`
	AudioParam          string        `orm:"notnull"`
	GpuNumber           int           `orm:"notnull"`
	GpuDeint            int           `orm:"notnull"`
	Score               int           `orm:"notnull"`
	Worker              string        `orm:"notnull"`
	Op                  int           `orm:"notnull"`
	State               int           `orm:"notnull"`
	CreateTime          time.Time
	ModifyTime          time.Time
	UpdateTime          time.Time
}

func (this *JobInfo)FindAllJobInfo() ([]*JobInfo, error) {

	var jobs []*JobInfo
	o := orm.NewOrm()
	qs := o.QueryTable("job_info")
	//qs = qs.Filter("is_activity",1)
	_, err := qs.OrderBy("-modify_time", "-create_time").All(&jobs)

	//common.PanicIf(err)
	return jobs, err

}

func (this *JobInfo)FindAllJobInfoByPage() ([]*JobInfo, error) {

	var jobs []*JobInfo
	o := orm.NewOrm()
	qs := o.QueryTable("job_info")
	if this.Name != "" {
		qs = qs.Filter("name", this.Name)
	}

	_, err := qs.OrderBy("-modify_time", "-create_time").All(&jobs)

	//common.PanicIf(err)
	return jobs, err

}

func (this *JobInfo)SaveJobInfo() error {

	_, err := orm.NewOrm().Insert(this)

	return err
}

func (this *JobInfo)GetJobInfoById() error {

	o := orm.NewOrm()
	return o.Read(this)
}

func (this *JobInfo)UpdateJobInfo() error {

	o := orm.NewOrm()
	id, err := o.Update(this, "name", "input_url", "input_url_backup", "output_url", "scale", "vcodec_type", "vcodec_bitrate", 
		"watermark_text", "watermark_param", "audio_param", "gpu_number", "gpu_deint", "worker", "op", "state", "score", "modify_time")
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("record no exist")
	}
	return nil
}

func (this *JobInfo)UpdateStateAndOp() error {

	o := orm.NewOrm()
	id, err := o.Update(this, "worker", "op", "state", "update_time", "modify_time")
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("record no exist")
	}
	return nil
}

func (this *JobInfo)OperateJobInfo() error {

	o := orm.NewOrm()
	id, err := o.Update(this,  "op", "state", "modify_time")
	if err != nil {
		return err
	}

	if id == 0 {
		return errors.New("record no exist")
	}
	return nil
}

func (this *JobInfo)GetJobInfo() (error) {

	o := orm.NewOrm()
	return o.Read(this)

}

func (this *JobInfo)DeleteJobInfo() error {

	err := this.GetJobInfo()
	o := orm.NewOrm()
	_, err = o.Delete(this)
	return err
}