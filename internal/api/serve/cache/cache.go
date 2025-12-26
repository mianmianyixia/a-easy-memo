package cache

import (
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"
	"time"
)

// 在缓存区查找任务

func FindTask(data dao.Data, redis dao.MemberTask) (*model.Member, error) {
	member := &model.Member{}
	ok, err := redis.IsExpired(data) //先检查是否过期
	if err != nil {
		return nil, err
	}
	if ok {
		return &model.Member{}, nil //如果过期返回空结构体
	} else {
		if data.TaskName == "" {
			result, err := redis.GetRedisList(data)
			if err != nil {
				return nil, err
			}
			results := result.(map[string]interface{})
			if len(results) == 0 {
				return member, nil //未能查到数据也返回空结构体，serve会在后面去数据库找
			}
			var tasks []model.Task
			for taskName, taskContext := range results {
				task := model.Task{
					TaskName:    taskName,
					TaskContent: taskContext.(string),
				}
				tasks = append(tasks, task)
			}
			member.Tasks = tasks
			return member, nil
		} else {
			results, err := redis.GetRedis(data)
			if err != nil {
				return nil, err
			}
			result := results.(string)
			member = &model.Member{
				UserName: data.Name,
				Tasks: []model.Task{
					{TaskName: data.TaskName, TaskContent: result},
				},
			}
			return member, nil
		}
	}
}

// 将某个任务存入缓存区

func AddTask(data dao.Data, redis dao.MemberTask, existTime time.Duration) error {
	err := redis.SetRedis(data, existTime)
	if err != nil {
		return err
	}
	return nil
}

// 删除某个任务

func DeleteTask(data dao.Data, redis dao.MemberTask) (bool, error) {
	ok, err := redis.DeleteRedis(data)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}
