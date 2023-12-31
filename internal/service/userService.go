package service

import (
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	responseparams "github.com/IRONICBo/QiYin_BE/internal/params/response"
)

type UserService interface {
	/*
		个人使用
	*/
	// GetTableUserList 获得全部TableUser对象
	GetTableUserList() []dao.ResUser

	// GetTableUserByUsername 根据username获得TableUser对象
	GetTableUserByUsername(name string) (dao.User, error)

	// GetTableUserById 根据user_id获得TableUser对象
	GetTableUserById(id string) (dao.ResUser, error)

	// InsertTableUser 将tableUser插入表内
	InsertTableUser(User *dao.User) bool

	// login with user
	Login(param *requestparams.UserParams) (*responseparams.UserResponse, error)

	/*
		他人使用
	*/
	// GetUserById 未登录情况下,根据user_id获得User对象
	//GetUserById(id int64) (User, error)
	//
	//// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
	//GetUserByIdWithCurId(id int64, curId int64) (User, error)

	// 根据token返回id
	// 接口:auth中间件,解析完token,将userid放入context
	//(调用方法:直接在context内拿参数"userId"的值)	fmt.Printf("userInfo: %v\n", c.GetString("userId"))
}
