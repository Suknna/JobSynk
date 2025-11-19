package job

import "time"

type Job struct {
	ID            string    `json:"id" gorm:"primaryKey"` // 主键
	Name          string    `json:"name"`                 // 岗位名称
	Company       string    `json:"company"`              // 岗位所属公司
	URL           string    `json:"url"`                  // 岗位沟通界面的url
	Salary        string    `json:"salary"`               // 薪资情况
	Location      string    `json:"location"`             // 工作位置
	Description   string    `json:"description"`          // 岗位描述
	PublishedAt   time.Time `json:"published_at"`         // 岗位发布时间
	Rating        int       `json:"rating"`               // 评分 1-5分，越低越垃圾
	Greeting      string    `json:"greeting"`             // 问候词
	CompletedAt   time.Time `json:"completed_at"`         // 打招呼语发送完成且获取到了hr的回复
	HRResponse    string    `json:"hr_response"`          // HR的回复内容
	Status        string    `json:"status"`               // 完成状态 success、failed、wait
	FailedMessage string    `json:"failed_message"`       // 失败原因
}
