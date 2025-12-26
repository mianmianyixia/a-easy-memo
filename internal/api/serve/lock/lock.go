package lock

import "a-easy-memo/internal/dao"

//设置互斥锁

func Locked(data dao.Data, redis dao.MemberTask) (bool, error) {
	lock, err := redis.Lock(data)
	if err != nil {
		return false, err
	}
	if !lock {
		return false, nil
	}
	return true, nil
}

//解锁

func DelLock(data dao.Data, redis dao.MemberTask) error {
	err := redis.Unlock(data)
	return err
}
