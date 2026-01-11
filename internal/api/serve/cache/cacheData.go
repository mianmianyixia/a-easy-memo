package cache

import (
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"
	"a-easy-memo/pkg/utils"
	"strings"
	"time"
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
//	@param	data：仅传入用户名

func GetCache(data dao.Data, redis dao.MemberTask, member dao.MemberData) (bool, error) {
	var key []string
	cursor := uint64(0)
	var uniqueKey []string
	seen := make(map[string]bool)
	for {
		nextCursor, appendKeys, err := redis.GetAll(data, cursor)
		if err != nil {
			return false, err
		}
		for _, appendKey := range appendKeys {
			if !seen[appendKey] {
				seen[appendKey] = true
				uniqueKey = append(uniqueKey, appendKey)
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	key = append(key, uniqueKey...) //获得对应用户备份缓存的所有键
	for _, keys := range key {
		var da dao.Data
		var task model.Task
		da.Name = strings.Replace(keys, "always", "", 1) //获得可过期缓存的键
		ok, err := redis.IsExpired(da)                   //判断是否过期
		if err != nil {
			return false, err
		}
		if ok { //已过期,从备份缓存拿数据，再存入数据库
			user, err := redis.GetData(keys)
			if err != nil {
				return false, err
			}
			mem, err := member.FindMemberData(user) //检查数据库是否存在该用户数据
			if err != nil {
				return false, err
			}
			if mem == nil {
				err, id := member.FindMemberID(user)
				if err != nil {
					return false, err
				}
				task = model.Task{
					TaskName:    user.TaskName,
					TaskContent: user.Value,
					MemberID:    id,
				}
				err = member.CreateTask(&task)
				if err != nil {
					return false, err
				}
				err = redis.SetRedis(da, 24*time.Hour+utils.RandomDuration(6))
				if err != nil {
					return false, err
				}
				err = redis.DeleteWrite(keys)
			} else {
				mem1 := mem[0]
				_, err := member.UpdateTask(&mem1)
				if err != nil {
					return false, err
				}
				err = redis.DeleteWrite(keys)
				if err != nil {
					return false, err
				}
				err = redis.SetRedis(da, 24*time.Hour+utils.RandomDuration(6))
				if err != nil {
					return false, err
				}
			}
		}
	}
	return true, nil
}

func Save(data dao.Data, redis dao.MemberTask, member dao.MemberData) error {
	var key []string
	cursor := uint64(0)
	var uniqueKey []string
	seen := make(map[string]bool)
	for {
		nextCursor, appendKeys, err := redis.GetAll(data, cursor)
		if err != nil {
			return err
		}
		for _, appendKey := range appendKeys {
			if !seen[appendKey] {
				seen[appendKey] = true
				uniqueKey = append(uniqueKey, appendKey)
			}
		}
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}
	key = append(key, uniqueKey...)
	if len(key) == 0 {
		return nil
	}
	for _, keys := range key {
		var da dao.Data
		Name := strings.Replace(keys, "always", "", 1)
		da, err := redis.GetData(keys)
		if err != nil {
			return err
		}
		mem, err := member.FindMemberData(da) //检查数据库是否存在该用户数据
		if err != nil {
			return err
		}
		if mem == nil {
			err, id := member.FindMemberID(da)
			if err != nil {
				return err
			}
			task := model.Task{
				TaskName:    da.TaskName,
				TaskContent: da.Value,
				MemberID:    id,
			}
			err = member.CreateTask(&task)
			if err != nil {
				return err
			}
			err = redis.SetRedis(da, 24*time.Hour+utils.RandomDuration(6))
			if err != nil {
				return err
			}
			err = redis.DeleteWrite(Name)
			if err != nil {
				return err
			}
			err = redis.DeleteWrite(keys)
			if err != nil {
				return err
			}
		} else {
			mem1 := mem[0]
			mem1.TaskContent = da.Value
			_, err := member.UpdateTask(&mem1)
			if err != nil {
				return err
			}
			err = redis.SetRedis(da, 24*time.Hour+utils.RandomDuration(6))
			if err != nil {
				return err
			}
			err = redis.DeleteWrite(Name)
			if err != nil {
				return err
			}
			err = redis.DeleteWrite(keys)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//
//func findcache(data dao.Data, redis dao.MemberTask) error {
//	var key []string
//	cursor := uint64(0)
//	var uniqueKey []string
//	seen := make(map[string]bool)
//	for {
//		nextCursor, appendKeys, err := redis.GetAll(data, cursor)
//		if err != nil {
//			return err
//		}
//		for _, appendKey := range appendKeys {
//			if !seen[appendKey] {
//				seen[appendKey] = true
//				uniqueKey = append(uniqueKey, appendKey)
//			}
//		}
//		if nextCursor == 0 {
//			break
//		}
//		cursor = nextCursor
//	}
//	key = append(key, uniqueKey...)
//	if len(key) == 0 {
//		return nil
//	}
//	for _, keys := range key {
//		var da dao.Data
//		Name := strings.Replace(keys, "always", "", 1)
//		da, err := redis.GetData(keys)
//		if err != nil {
//			return err
//		}
//		fmt.Printf("keys: %v,value:%v\n", da.TaskName, da.Value)
//		fmt.Printf("keys: %v\n", Name)
//	}
//	return nil
//}
