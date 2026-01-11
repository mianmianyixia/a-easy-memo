package dao

import (
	"a-easy-memo/internal/model"
	"errors"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func NewGorm(db *gorm.DB) *Gorm {
	return &Gorm{db: db}
}

// 建立用户

func (db Gorm) CreateMemberData(member *model.Member) error {
	result := db.db.Create(member)
	return result.Error
}

// 创建数据库任务

func (db Gorm) CreateTask(task *model.Task) error {
	result := db.db.Create(task)
	return result.Error
}

// 更新任务数据库

func (db Gorm) UpdateTask(tasks *model.Task) (*model.Task, error) {
	task := &model.Task{}
	result := db.db.Model(task).Where("member_id = ? AND task_name= ?", tasks.MemberID, tasks.TaskName).Updates(tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return task, nil
}

// 删除任务

func (db Gorm) DeleteTask(data Data) (bool, error) {
	var member model.Member
	result := db.db.Where("user_name= ?", data.Name).First(&member)
	if result.Error != nil {
		return false, result.Error
	}
	if data.TaskName == "" {
		res := db.db.Where("member_id = ?", member.ID).Delete(&model.Task{})
		if res.Error != nil {
			return false, res.Error
		}
		return true, nil
	}
	res := db.db.Where("member_id = ? AND task_name = ?", member.ID, data.TaskName).Delete(&model.Task{})
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

// 数据库查找任务

func (db Gorm) FindMemberData(data Data) ([]model.Task, error) {
	member := &model.Member{}
	var tasks []model.Task
	err, ID := db.FindMemberID(data)
	if err != nil {
		return nil, err
	}
	(*member).ID = ID
	if data.TaskName == "" {
		results := db.db.Where("member_id= ?", member.ID).Find(&tasks)
		if results.Error != nil {
			return nil, results.Error
		}
		if len(tasks) == 0 {
			return nil, nil
		}
		return tasks, nil
	} else {
		results := db.db.Where("member_id = ? AND task_name = ?", member.ID, data.TaskName).Find(&tasks)
		if results.Error != nil {
			return nil, results.Error
		}
		if len(tasks) == 0 {
			return nil, nil
		}
		return tasks, nil
	}
}

// 删除数据库用户的数据

func (db Gorm) DeleteMemberData(data Data) error {
	result := db.db.Where("user_name=?", data.Name).Delete(&model.Member{})
	return result.Error
}

// 更新用户数据

func (db Gorm) UpdateMemberData(members *model.Member) (*model.Member, error) {
	member := &model.Member{}
	result := db.db.Model(member).Where("user_name=?", (*members).UserName).Select("pass_word").Updates(members)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return member, result.Error
}

// 查找数据库中的用户

func (db Gorm) FindMember(members *model.Member) (error, bool) {
	member := &model.Member{}
	if (*members).UserName == "" {
		return errors.New("用户名为空"), false
	}
	result := db.db.Where("user_name", (*members).UserName).First(member)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false
		}
		return result.Error, false
	}
	if (*member).PassWord != (*members).PassWord {
		return nil, false
	}
	return nil, true
}

//查询用户id

func (db Gorm) FindMemberID(data Data) (error, uint) {
	member := &model.Member{}
	result := db.db.Where("user_name= ?", data.Name).First(member)
	return result.Error, (*member).ID
}
