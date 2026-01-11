package api

import (
	"a-easy-memo/internal/api/response"
	"a-easy-memo/internal/api/serve"
	"a-easy-memo/internal/api/serve/cache"
	"a-easy-memo/internal/api/serve/time"
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"

	"github.com/gin-gonic/gin"
)

// 接受数据并调用

func Find(db dao.MemberData, redis dao.MemberTask) func(c *gin.Context) {
	return func(c *gin.Context) {
		var tasks []model.TaskRequest
		if err := c.ShouldBindJSON(&tasks); err != nil {
			_ = c.Error(err)
			response.RequestError(c, err)
			return
		}
		userName := tasks[0].UserName
		var taskNames []string
		for _, task := range tasks {
			taskNames = append(taskNames, task.TaskName)
		}
		resmem, err := serve.Find(db, redis, userName, taskNames)
		if err != nil {
			_ = c.Error(err)
			response.InternalError(c, err)
			return
		}

		response.Success(c, *resmem)
	}
}

func Delete(db dao.MemberData, redis dao.MemberTask) func(c *gin.Context) {
	return func(c *gin.Context) {
		var tasks []model.TaskRequest
		if err := c.ShouldBindJSON(&tasks); err != nil {
			_ = c.Error(err)
			response.RequestError(c, err)
			return
		}
		userName := tasks[0].UserName
		var taskNames []string
		for _, task := range tasks {
			taskNames = append(taskNames, task.TaskName)
		}
		err, taskErr := serve.Del(db, redis, userName, taskNames)
		if err != nil {
			_ = c.Error(err)
			response.InternalError(c, err)
			return
		}
		if taskErr != nil {
			response.RequestTaskError(c, taskErr)
			return
		}
		response.Success(c, "成功删除")
	}
}

func Change(db dao.MemberData, redis dao.MemberTask) func(c *gin.Context) {
	return func(c *gin.Context) {
		var tasks []model.TaskRequest
		if err := c.ShouldBindJSON(&tasks); err != nil {
			_ = c.Error(err)
			response.RequestError(c, err)
			return
		}
		userName := tasks[0].UserName
		task := make(map[string]string)
		cacheUpdater := time.NewCacheUpdater(userName, redis, db)
		(*cacheUpdater).Start()
		for _, t := range tasks {
			task[t.TaskName] = t.TaskContent
		}
		err := serve.Add(redis, userName, task)
		if err != nil {
			_ = c.Error(err)
			response.InternalError(c, err)
			return
		}
		response.Success(c, "成功修改")
	}
}

func Save(db dao.MemberData, redis dao.MemberTask) func(c *gin.Context) {
	return func(c *gin.Context) {
		var ok model.SaveRequest
		if err := c.ShouldBindJSON(&ok); err != nil {
			_ = c.Error(err)
			response.RequestError(c, err)
			return
		}
		if ok.Save {
			var data dao.Data
			data.Name = ok.UserName
			err := cache.Save(data, redis, db)
			if err != nil {
				_ = c.Error(err)
				response.InternalError(c, err)
				return
			}
			response.Success(c, "成功保存")
		} else {
			response.Success(c, "仍未保存")
		}
	}
}
