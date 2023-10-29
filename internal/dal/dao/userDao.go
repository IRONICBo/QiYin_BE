package dao

import (
	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"log"
)

//// User 对应数据库User表结构的结构体
//type User struct {
//	Id       string
//	Name     string
//	Password string
//}

type User struct {
	Id              string
	Name            string
	Password        string
	Avatar          string `json:"avatar" gorm:"column:avatar"`
	BackgroundImage string `json:"background_image" gorm:"column:background_image"`
	Signature       string `json:"signature" gorm:"column:signature"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
}

// TableName 修改表名映射
func (user User) TableName() string {
	return "users"
}

// GetTableUserList 获取全部TableUser对象
func GetTableUserList() ([]User, error) {
	tableUsers := []User{}
	if err := db.GetMysqlDB().Table("users").Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

func QueryUserLogin(username string, key string) (User, bool) {
	var user User
	res := db.GetMysqlDB().Table("users").Where(key+" = ?", username).First(&user)
	if res.Error != nil || res.RowsAffected == 0 {
		return User{}, false
	}
	return user, true
}

// GetTableUserByUsername 根据username获得TableUser对象
func GetTableUserByUsername(name string) (User, error) {
	user := User{}
	if err := db.GetMysqlDB().Where("name = ?", name).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// GetTableUserById 根据user_id获得TableUser对象
func GetTableUserById(id int64) (User, error) {
	user := User{}
	if err := db.GetMysqlDB().Table("users").Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// InsertTableUser 将tableUser插入表内
func InsertTableUser(user *User) bool {
	if err := db.GetMysqlDB().Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
