package service

import (
	"errors"
	"log"

	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	responseparams "github.com/IRONICBo/QiYin_BE/internal/params/response"
	util "github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/IRONICBo/QiYin_BE/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserServiceImpl struct {
	Service
	FavoriteService
	CollectionService
}

// NewUserService return new service with gin context.
func NewUserService(c *gin.Context) *UserServiceImpl {
	return &UserServiceImpl{
		Service: Service{
			ctx: c,
		},
		FavoriteService:   &FavoriteServiceImpl{},
		CollectionService: &CollectionServiceImpl{},
	}
}

// GetTableUserList 获得全部TableUser对象.
func (usi *UserServiceImpl) GetTableUserList() []dao.ResUser {
	tableUsers, err := dao.GetTableUserList()
	if err != nil {
		log.Println("Err:", err.Error())
		return tableUsers
	}
	return tableUsers
}

// GetTableUserByUsername 根据username获得TableUser对象.
func (usi *UserServiceImpl) GetTableUserByUsername(name string) (dao.User, error) {
	user, err := dao.GetTableUserByUsername(name)
	if err != nil {
		log.Println("Err:", err.Error())
		return user, err
	}
	log.Println("Query User Success")
	return user, nil
}

func (usi *UserServiceImpl) IsUserExistByName(name string) bool {
	user, ok := dao.QueryUserLogin(name, "name")
	if !ok || user.Id == "" {
		return false
	}
	return true
}

// GetTableUserById 根据user_id获得TableUser对象.
func (usi *UserServiceImpl) GetTableUserById(id string) (dao.ResUser, error) {
	user, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("User Not Found", err.Error())
		return user, err
	}
	log.Println("Query User Success")
	return user, nil
}

// InsertTableUser 将tableUser插入表内.
func (usi *UserServiceImpl) InsertTableUser(User *dao.User) bool {
	flag := dao.InsertTableUser(User)
	if flag == false {
		log.Println("insert failed")
		return false
	}
	return true
}

// login with user.
func (usi *UserServiceImpl) Login(param *requestparams.UserParams) (*responseparams.UserResponse, error) {
	resp := &responseparams.UserResponse{}

	// Check user
	user, err := usi.GetTableUserByUsername(param.Name)
	if err != nil {
		return resp, err
	}

	if user.Id == "" {
		return resp, errors.New("the user has not registered")
	}

	// Check password
	if !utils.ComparePassword(param.Password, user.Password) {
		return resp, errors.New("password is not correct")
	}

	// Get token
	token_string, _, err := util.GenerateJwtToken(user.Id)
	if err != nil {
		return resp, err
	}

	resp = &responseparams.UserResponse{
		Token:  token_string,
		UserId: user.Id,
	}
	return resp, nil
}

// register with user.
func (usi *UserServiceImpl) Register(param *requestparams.UserParams) (*responseparams.UserResponse, error) {
	resp := &responseparams.UserResponse{}

	// Check user
	has := usi.IsUserExistByName(param.Name)
	if has {
		return resp, errors.New("the user has registered")
	}

	uuid := utils.GenUUID()
	newUser := dao.User{
		Name:      param.Name,
		Password:  utils.EncryptPassword(param.Password),
		Id:        uuid,
		Avatar:    "",
		Signature: "",
	}
	if usi.InsertTableUser(&newUser) != true {
		return resp, errors.New("insert failed")
	}
	// Get token
	token_string, _, err := util.GenerateJwtToken(newUser.Id)
	if err != nil {
		return resp, err
	}

	resp = &responseparams.UserResponse{
		Token:  token_string,
		UserId: uuid,
	}
	return resp, nil
}

// 获得用户信息  不需要登录  只有点赞操作和评论需要登录.
func (usi *UserServiceImpl) GetUserById(id string) (dao.ResUser, error) {
	user := dao.ResUser{}
	user, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("User Not Found")
		return user, err
	}
	log.Println("Query User Success")

	u := GetFavoriteService() //解决循环依赖
	// 获取点赞以及被点赞的数量
	totalFavorited, _ := u.TotalFavorite(id)
	favoritedCount, _ := u.FavoriteVideoCount(id)
	c := GetCollectionService()
	totalCollection, _ := c.TotalCollection(id)
	collectionCount, _ := c.CollectionVideoCount(id)
	user.TotalFavorited = totalFavorited
	user.FavoriteCount = favoritedCount
	user.TotalCollected = totalCollection
	user.CollectionCount = collectionCount

	return user, nil
}
