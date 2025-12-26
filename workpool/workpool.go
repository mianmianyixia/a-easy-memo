package workpool

import (
	"sync"
)

// 创建一个协程池加快读取数据库
type GetTool struct {
	ToolChan chan func()
	Tasks    int
	wg       sync.WaitGroup
}

func CreatTask(worker int) *GetTool { //任务分发
	get := &GetTool{
		ToolChan: make(chan func()),
		Tasks:    worker,
		wg:       sync.WaitGroup{},
	}
	for range get.Tasks {
		go get.Worker()
	}
	return get
}
func (get *GetTool) Worker() { //工作
	for task := range get.ToolChan {
		task()
	}
	get.wg.Done()
}
func (get *GetTool) AddTask(task func()) { //添加任务
	get.wg.Add(1)
	get.ToolChan <- task
}
func (get *GetTool) Wait() { //等待任务完成
	get.wg.Wait()
}
func (get *GetTool) Close() { //关闭通道
	close(get.ToolChan)
}
