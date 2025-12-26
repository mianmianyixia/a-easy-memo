package model

type Member struct {
	ID       uint   `gorm:"primary_key;auto_increment" json:"id"`
	UserName string `json:"user_name" gorm:"size:255;not null;index"`
	PassWord string `json:"pass_word" gorm:"size:255;not null"`
	Tasks    []Task `json:"tasks" gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
}

func (Member) TableName() string {
	return "member" // 或 "members"，必须一致
}

type Task struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	TaskName    string `json:"task_name" gorm:"size:255;not null;index"`
	TaskContent string `json:"task_content" gorm:"size:6553"`
	MemberID    uint   `json:"member_id" gorm:"index"`
}

func (Task) TableName() string {
	return "task" // 或 "tasks"
}
