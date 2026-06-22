package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

const TableNameJob = "job"

type JobStatus string

const (
	JobStatusNormal    JobStatus = "Normal"
	JobStatusPending   JobStatus = "Pending"
	JobStatusRunning   JobStatus = "Running"
	JobStatusSucceeded JobStatus = "Succeeded"
	JobStatusFailed    JobStatus = "Failed"
	JobStatusUnknown   JobStatus = "Unknown"
)

type Job struct {
	ID        int64     `gorm:"column:id" json:"id"`                 // 任务 ID
	Name      string    `gorm:"column:name" json:"name"`             // 任务名称
	Namespace string    `gorm:"column:namespace" json:"namespace"`   // 任务 Kubernetes Namespace
	Info      JobInfo   `gorm:"column:info" json:"info"`             // 任务 Kubernetes 相关信息
	Status    JobStatus `gorm:"column:status" json:"status"`         // 任务状态
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"` // 修改时间
}

func (*Job) TableName() string {
	return TableNameJob
}

type JobInfo struct {
	Image   string   `json:"image"`
	Command []string `json:"command"`
	Args    []string `json:"args"`
}

// Scan implements the [sql.Scanner] interface.
func (ti *JobInfo) Scan(value any) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), ti)
}

// Value implements the [driver.Valuer] interface.
func (ti JobInfo) Value() (driver.Value, error) {
	bytes, err := json.Marshal(ti)
	return string(bytes), err
}
