package service

import (
	"errors"
	"github.com/IRONICBo/QiYin_BE/internal/dal/dao"
	requestparams "github.com/IRONICBo/QiYin_BE/internal/params/request"
	responseparams "github.com/IRONICBo/QiYin_BE/internal/params/response"
	util "github.com/IRONICBo/QiYin_BE/internal/utils"
	"github.com/IRONICBo/QiYin_BE/pkg/utils"
	"github.com/gin-gonic/gin"
	"log"
)

type UserServiceImpl struct {
	Service
	//FollowService
	//LikeService
}

// NewCommonService return new service with gin context.
func NewUserService(c *gin.Context) *UserServiceImpl {
	return &UserServiceImpl{
		Service: Service{
			ctx: c,
		},
	}
}

// GetTableUserList 获得全部TableUser对象
func (usi *UserServiceImpl) GetTableUserList() []dao.User {
	tableUsers, err := dao.GetTableUserList()
	if err != nil {
		log.Println("Err:", err.Error())
		return tableUsers
	}
	return tableUsers
}

// GetTableUserByUsername 根据username获得TableUser对象
func (usi *UserServiceImpl) GetTableUserByUsername(name string) (dao.User, error) {
	user, err := dao.GetTableUserByUsername(name)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
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

// GetTableUserById 根据user_id获得TableUser对象
func (usi *UserServiceImpl) GetTableUserById(id int64) dao.User {
	User, err := dao.GetTableUserById(id)
	if err != nil {
		log.Println("Err:", err.Error())
		log.Println("User Not Found")
		return User
	}
	log.Println("Query User Success")
	return User
}

// InsertTableUser 将tableUser插入表内
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

//
//// GetUserById 未登录情况下,根据user_id获得User对象
//func (usi *UserServiceImpl) GetUserById(id int64) (User, error) {
//	user := User{
//		Id:             0,
//		Name:           "",
//		FollowCount:    0,
//		FollowerCount:  0,
//		IsFollow:       false,
//		TotalFavorited: 0,
//		FavoriteCount:  0,
//	}
//	User, err := dao.GetTableUserById(id)
//	if err != nil {
//		log.Println("Err:", err.Error())
//		log.Println("User Not Found")
//		return user, err
//	}
//	log.Println("Query User Success")
//	followCount, _ := usi.GetFollowingCnt(id)
//	if err != nil {
//		log.Println("Err:", err.Error())
//	}
//	followerCount, _ := usi.GetFollowerCnt(id)
//	if err != nil {
//		log.Println("Err:", err.Error())
//	}
//	u := GetLikeService() //解决循环依赖
//	totalFavorited, _ := u.TotalFavourite(id)
//	favoritedCount, _ := u.FavouriteVideoCount(id)
//	user = User{
//		Id:             id,
//		Name:           User.Name,
//		FollowCount:    followCount,
//		FollowerCount:  followerCount,
//		IsFollow:       false,
//		TotalFavorited: totalFavorited,
//		FavoriteCount:  favoritedCount,
//	}
//	return user, nil
//}
//
//// GetUserByIdWithCurId 已登录(curID)情况下,根据user_id获得User对象
//func (usi *UserServiceImpl) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
//	user := User{
//		Id:             0,
//		Name:           "",
//		FollowCount:    0,
//		FollowerCount:  0,
//		IsFollow:       false,
//		TotalFavorited: 0,
//		FavoriteCount:  0,
//	}
//	User, err := dao.GetTableUserById(id)
//	if err != nil {
//		log.Println("Err:", err.Error())
//		log.Println("User Not Found")
//		return user, err
//	}
//	log.Println("Query User Success")
//	followCount, err := usi.GetFollowingCnt(id)
//	if err != nil {
//		log.Println("Err:", err.Error())
//	}
//	followerCount, err := usi.GetFollowerCnt(id)
//	if err != nil {
//		log.Println("Err:", err.Error())
//	}
//	isfollow, err := usi.IsFollowing(curId, id)
//	if err != nil {
//		log.Println("Err:", err.Error())
//	}
//	u := GetLikeService() //解决循环依赖
//	totalFavorited, _ := u.TotalFavourite(id)
//	favoritedCount, _ := u.FavouriteVideoCount(id)
//	user = User{
//		Id:             id,
//		Name:           User.Name,
//		FollowCount:    followCount,
//		FollowerCount:  followerCount,
//		IsFollow:       isfollow,
//		TotalFavorited: totalFavorited,
//		FavoriteCount:  favoritedCount,
//	}
//	return user, nil
//}
//
//// GenerateToken 根据username生成一个token
//func GenerateToken(username string) string {
//	u := UserService.GetTableUserByUsername(new(UserServiceImpl), username)
//	fmt.Printf("generatetoken: %v\n", u)
//	token := NewToken(u)
//	println(token)
//	return token
//}
//
//// NewToken 根据信息创建token
//func NewToken(u dao.User) string {
//	expiresTime := time.Now().Unix() + int64(config.OneDayOfHours)
//	fmt.Printf("expiresTime: %v\n", expiresTime)
//	id64 := u.Id
//	fmt.Printf("id: %v\n", strconv.FormatInt(id64, 10))
//	claims := jwt.StandardClaims{
//		Audience:  u.Name,
//		ExpiresAt: expiresTime,
//		Id:        strconv.FormatInt(id64, 10),
//		IssuedAt:  time.Now().Unix(),
//		Issuer:    "tiktok",
//		NotBefore: time.Now().Unix(),
//		Subject:   "token",
//	}
//	var jwtSecret = []byte(config.Secret)
//	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
//		token = "Bearer " + token
//		println("generate token success!\n")
//		return token
//	} else {
//		println("generate token fail\n")
//		return "fail"
//	}
//}
//
//// EnCoder 密码加密
//func EnCoder(password string) string {
//	h := hmac.New(sha256.New, []byte(password))
//	sha := hex.EncodeToString(h.Sum(nil))
//	fmt.Println("Result: " + sha)
//	return sha
//}
