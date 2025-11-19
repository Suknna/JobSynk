package database

import (
	"JobSynk/job"
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	var err error
	db, err = gorm.Open(sqlite.Open("jobs.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 自动迁移创建表
	err = db.AutoMigrate(&job.Job{})
	if err != nil {
		return err
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// GetJobsByStatus 根据状态查询所有工作记录
func GetJobsByStatus(status string) ([]job.Job, error) {
	var jobs []job.Job
	result := db.Where("status = ?", status).Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}
	return jobs, nil
}

// CreateJob 创建新的工作记录
func CreateJob(job *job.Job) error {
	result := db.Create(job)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// DeleteJob 根据ID删除工作记录
func DeleteJob(id string) error {
	result := db.Where("id = ?", id).Delete(&job.Job{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("job not found")
	}
	return nil
}

// UpdateJobStatus 更新工作状态（只能更新允许的字段）
func UpdateJobStatus(id string, updates map[string]interface{}) error {
	// 只允许更新指定的字段
	allowedFields := map[string]bool{
		"status":         true,
		"failed_message": true,
		"hr_response":    true,
		"completed_at":   true,
	}

	// 过滤不允许的字段
	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	// 如果没有有效的更新字段，返回错误
	if len(filteredUpdates) == 0 {
		return errors.New("no valid fields to update")
	}

	result := db.Model(&job.Job{}).Where("id = ?", id).Updates(filteredUpdates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("job not found")
	}
	return nil
}
