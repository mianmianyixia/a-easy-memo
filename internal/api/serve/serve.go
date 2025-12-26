package serve

import (
	"a-easy-memo/internal/api/serve/cache"
	"a-easy-memo/internal/api/serve/lock"
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"
	"a-easy-memo/pkg/utils"
	"a-easy-memo/workpool"
	"errors"
	"time"
)

// 把服务层接收数据改在api层
// 查找对应任务,目前这个仅图简易，能跑，但问题还很多

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
		if len((*mem).Tasks) > 1 {
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
			members, err := db.FindMemberData(data)
			if err != nil {
				return nil, err
			}
			if members == nil {
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
				for _, value := range (*members).Tasks {
					data.Value = value.TaskContent
					err = cache.AddTask(data, task, 24*time.Hour+utils.RandomDuration(6)) //设置随机存在时间避免缓存雪崩
					if err != nil {
						return nil, err
					}
				}
			}
			result = append(result, (*members).Tasks[0])
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

func Add(db dao.MemberData, task dao.MemberTask, name string, tasks []string) error {
	var data dao.Data
	var tas *model.Task
	data.Name = name
	tas = &model.Task{}
	db.UpdateTask(tas)
	for _, taskName := range tasks {
		data.TaskName = taskName
		err := cache.AddTask(data, task, 60*time.Minute)
		if err != nil {
			return err
		}
	}
}

func Del(db dao.MemberData, task dao.MemberTask, name string, tasks []string) error {
	var data dao.Data
	var finerr error
	work := workpool.CreatTask(100)
	data.Name = name
	for _, taskName := range tasks {
		mem, err := Find(db, task, name, tasks)
		if err != nil {
			return err
		}
		for _, t := range (*mem).Tasks {
			if t.TaskContent != "null" {
				work.AddTask(func() {
					data.TaskName = taskName
					_, err := cache.DeleteTask(data, task)
					if err != nil {
						finerr = err
						return
					}
					_, err = db.DeleteTask(data)
					if err != nil {
						finerr = err
						return
					}
				})
				if finerr != nil {
					return finerr
				}
			}
		}
	}
	work.Close()
	work.Wait()
	return nil
}

func Change(db dao.MemberData, task dao.MemberTask, name string, tasks []string) error {

}

func PutSql() {

}
