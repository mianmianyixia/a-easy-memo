package login

import (
	"a-easy-memo/internal/api/response"
	"a-easy-memo/internal/dao"
	"a-easy-memo/internal/model"
	"a-easy-memo/pkg/utils"

	"github.com/gin-gonic/gin"
)

// 注意修改，使用依赖注入
// 该层重构为用户操作层，api层应只写接收并处理前端数据
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
		if err := db.FindMember(member.UserName, member.PassWord); err != nil {
			response.RequestError(c, err)
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
		err := db.DeleteMemberData(member.UserName)
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
