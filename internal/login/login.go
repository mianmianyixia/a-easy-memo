package login

import (
	"a-easy-memo/internal/api/response"
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"
	"a-easy-memo/pkg/utils"
	"errors"

	"github.com/gin-gonic/gin"
)

func Register(db dao.MemberData) func(c *gin.Context) {
	return func(c *gin.Context) {
		var member model.Member
		if err := c.ShouldBindJSON(&member); err != nil {
			response.RequestError(c, err)
			return
		}
		err := db.CreateMemberData(&member)
		if err != nil {
			response.RequestError(c, err)
			return
		}
		response.Success(c, "注册成功")
	}
}
func Login(db dao.MemberData) func(c *gin.Context) {
	return func(c *gin.Context) {
		var member model.Member
		if err := c.ShouldBindJSON(&member); err != nil {
			response.RequestError(c, err)
			return
		}
		if err, ok := db.FindMember(&member); err != nil || ok == false {
			response.RequestError(c, errors.New("密码或其他错误"))
			return
		}
		token, err := utils.MakeToken(member.UserName)
		if err != nil {
			response.InternalError(c, err)
			return
		}
		response.Success(c, token)
	}
}
func Del(db dao.MemberData) func(c *gin.Context) {
	return func(c *gin.Context) {
		var member model.Member
		if err := c.ShouldBindJSON(&member); err != nil {
			response.RequestError(c, err)
			return
		}
		var data dao.Data
		data.Name = member.UserName
		err := db.DeleteMemberData(data)
		if err != nil {
			response.InternalError(c, err)
			return
		}
		response.Success(c, "成功删除")
	}
}
func Change(db dao.MemberData) func(c *gin.Context) {
	return func(c *gin.Context) {
		var member model.Member
		if err := c.ShouldBindJSON(&member); err != nil {
			response.RequestError(c, err)
			return
		}
		_, err := db.UpdateMemberData(&member)
		if err != nil {
			response.InternalError(c, err)
			return
		}
		response.Success(c, "成功修改密码")
	}
}
