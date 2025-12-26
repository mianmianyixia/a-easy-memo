package cache

import (
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"
	"strings"
)

//将更新数据放入缓存

func UpdateCache(data dao.Data, redis dao.MemberTask) error {
	err := redis.UpdateCache(data)
	if err != nil {
		return err
	}
	err = redis.AlwaysCache(data)
	if err != nil {
		return err
	}
	return nil
}

//提取缓存所有数据，使用计时器，在用户无操作后设置的缓存过期后将缓存存入数据库，想要立即存储，必须先保存

func GetCache(data dao.Data, redis dao.MemberTask, member dao.MemberData) (bool, error) {
	var key []string
	var tasks []model.Task
	cursor := uint64(0)
	for {
		nextCursor, appendKeys, err := redis.GetAll(data, cursor)
		if err != nil {
			return false, err
		}
		key = append(key, appendKeys...)
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	for _, keys := range key {
		var da dao.Data
		var task model.Task
		da.Name = keys
		da.Name = strings.Replace(keys, "always", "", 1)
		ok, err := redis.IsExpired(da)
		if err != nil {
			return false, err
		}
		if ok {
			user, err := redis.GetData(keys)
			if err != nil {
				return false, err
			}
			mem, err := member.FindMemberData(user)
			if err != nil {
				return false, err
			}
			if mem == nil {
				task = model.Task{}
				member.CreateTask(&task)
			}
		}
	}
}
