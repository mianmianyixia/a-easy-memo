package dao

import (
	"a-easy-memo/internal/model"
	"time"
)

type Data struct {
	Name     string
	TaskName string
	Value    string
}

type MemberData interface {
	CreateMemberData(member *model.Member) error
	FindMemberData(data Data) ([]model.Task, error)
	DeleteMemberData(data Data) error
	UpdateMemberData(members *model.Member) (*model.Member, error)
	FindMember(members *model.Member) (error, bool)
	FindMemberID(data Data) (error, uint)
	CreateTask(task *model.Task) error
	UpdateTask(tasks *model.Task) (*model.Task, error)
	DeleteTask(data Data) (bool, error)
}
type MemberTask interface {
	SetRedis(data Data, existTime time.Duration) error
	GetRedis(data Data) (interface{}, error)
	GetRedisList(data Data) (interface{}, error)
	DeleteRedis(data Data) (bool, error)
	IsExpired(data Data) (bool, error)
	Lock(data Data) (bool, error)
	Unlock(data Data) error
	UpdateCache(data Data) error
	AlwaysCache(data Data) error
	GetAll(data Data, cursor uint64) (uint64, []string, error)
	GetData(data string) (Data, error)
	DeleteWrite(data string) error
}
