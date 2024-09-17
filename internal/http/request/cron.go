package request

type CronCreate struct {
	Name       string `form:"name" json:"name"`
	Type       string `form:"type" json:"type"`
	Time       string `form:"time" json:"time"`
	Script     string `form:"script" json:"script"`
	BackupType string `form:"backup_type" json:"backup_type"`
	BackupPath string `form:"backup_path" json:"backup_path"`
	Target     string `form:"target" json:"target"`
	Save       int    `form:"save" json:"save"`
}

type CronUpdate struct {
	ID     uint   `form:"id" json:"id"`
	Name   string `form:"name" json:"name"`
	Time   string `form:"time" json:"time"`
	Script string `form:"script" json:"script"`
}

type CronStatus struct {
	ID     uint `form:"id" json:"id"`
	Status bool `form:"status" json:"status"`
}
