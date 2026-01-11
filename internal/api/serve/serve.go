package serve

import (
	"a-easy-memo/internal/api/serve/cache"
	"a-easy-memo/internal/api/serve/lock"
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"
	"a-easy-memo/pkg/utils"
	"errors"
	"fmt"
	"time"
)

//加入保存功能，每次编辑必须保存才能查询，set建立写缓存，每个操作先保存后再放入数据库，未保存会定时存入

// @param name: 解析的用户名字
// 		  tasks:解析的任务名

func Find(db dao.MemberData, task dao.MemberTask, name string, tasks []string) (*model.Member, error) {
	var data dao.Data
	member := &model.Member{}
	data.Name = name
	var result []model.Task
	for _, taskName := range tasks {
		data.TaskName = taskName
		mem, err := cache.FindTask(data, task)
		if err != nil {
			return nil, errors.New("查询缓存失败")
		}
		if len((*mem).Tasks) > 1 { //有的缓存过期的无法查到了
			member = mem
			return member, nil
		} else if len((*mem).Tasks) == 1 {
			result = append(result, (*mem).Tasks[0])
		} else {
			ok, err := lock.Locked(data, task) //设置锁，防止缓存击穿
			if err != nil {
				return nil, err
			}
			if !ok {
				time.Sleep(30 * time.Second)
				return Find(db, task, name, tasks)
			}
			tas, err := db.FindMemberData(data)
			if err != nil {
				return nil, err
			}
			if tas == nil {
				data.Value = "null" //缓存空值，防止缓存穿透
				err = cache.AddTask(data, task, 60*time.Minute)
				if err != nil {
					return nil, err
				}
				var null model.Task
				null.TaskName = taskName
				null.TaskContent = data.Value
				result = append(result, null)
				continue
			} else {
				for _, value := range tas {
					data.Value = value.TaskContent
					err = cache.AddTask(data, task, 24*time.Hour+utils.RandomDuration(6)) //设置随机存在时间避免缓存雪崩
					if err != nil {
						return nil, err
					}
				}
			}
			result = append(result, tas[0])
			err = lock.DelLock(data, task)
			if err != nil {
				return nil, err
			}
		}
	}
	(*member).Tasks = result
	(*member).UserName = name
	return member, nil
}

// 直接添加进写缓存

func Add(task dao.MemberTask, name string, tasks map[string]string) error {
	var data dao.Data
	data.Name = name
	for taskName, taskValue := range tasks {
		data.TaskName = taskName
		data.Value = taskValue
		err := cache.UpdateCache(data, task)
		if err != nil {
			return err
		}
	}
	return nil
}

//亦需修改成定时缓存中删除

func Del(db dao.MemberData, task dao.MemberTask, name string, tasks []string) (error, []string) {
	var data dao.Data
	data.Name = name
	m, err := Find(db, task, data.Name, tasks)
	if err != nil {
		return err, nil
	}
	var taskErr []string
	for _, t := range m.Tasks {
		data.TaskName = t.TaskName
		if t.TaskContent != "null" {
			_, err := db.DeleteTask(data)
			if err != nil {
				return err, nil
			}
			_, err = cache.DeleteTask(data, task)
			if err != nil {
				return err, nil
			}
			key := data.Name + ":" + data.TaskName
			alwaysKey := "always" + data.Name + ":" + data.TaskName
			err = task.DeleteWrite(alwaysKey)
			if err != nil {
				return err, nil
			}
			var existData dao.Data
			existData.Name = key
			ok, err := task.IsExpired(existData) //判断是否过期
			if err != nil {
				return err, nil
			}
			if !ok {
				err = task.DeleteWrite(key)
				if err != nil {
					return err, nil
				}
			}
		} else {
			taskErr = append(taskErr, fmt.Sprintf("并未存在该任务%s:%s", data.Name, data.TaskName))
		}
	}
	if len(taskErr) > 0 {
		return nil, taskErr
	}
	return nil, nil
}
