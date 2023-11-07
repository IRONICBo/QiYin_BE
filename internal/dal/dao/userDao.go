package dao

import (
	"log"

	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
)

type ResUser struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	//BackgroundImage string `json:"background_image"`
	Style           string `json:"style"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	FavoriteCount   int64  `json:"favorite_count"`
	TotalCollected  int64  `json:"total_collected"`
	CollectionCount int64  `json:"collection_count"`
}

type User struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	//BackgroundImage string `json:"background_image,omitempty"`
	Signature string `json:"signature,omitempty"`
	Password  string `json:"password,omitempty"`
	Style     int64  `json:"style"`
}

// TableName 修改表名映射.
func (user User) TableName() string {
	return "users"
}

// GetTableUserList 获取全部TableUser对象.
func GetTableUserList() ([]ResUser, error) {
	tableUsers := []ResUser{}
	if err := db.GetMysqlDB().Table("users").Find(&tableUsers).Error; err != nil {
		log.Println(err.Error())
		return tableUsers, err
	}
	return tableUsers, nil
}

func QueryUserLogin(username string, key string) (ResUser, bool) {
	var user ResUser
	res := db.GetMysqlDB().Table("users").Where(key+" = ?", username).First(&user)
	if res.Error != nil || res.RowsAffected == 0 {
		return ResUser{}, false
	}
	return user, true
}

// GetTableUserByUsername 根据username获得TableUser对象.
func GetTableUserByUsername(name string) (User, error) {
	user := User{}
	if err := db.GetMysqlDB().Where("name = ?", name).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// GetTableUserById 根据user_id获得TableUser对象.
func GetTableUserById(id string) (ResUser, error) {
	user := ResUser{}
	if err := db.GetMysqlDB().Table("users").Where("id = ?", id).First(&user).Error; err != nil {
		log.Println(err.Error())
		return user, err
	}
	return user, nil
}

// InsertTableUser 将tableUser插入表内.
func InsertTableUser(user *User) bool {
	if err := db.GetMysqlDB().Create(&user).Error; err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// GetUserIdByName
// 通过关键字搜索用户  name.
func GetUserIdByName(value string) ([]ResUser, error) {
	var userList []ResUser
	result := db.GetMysqlDB().Table("users").Where("name LIKE ?", "%"+value+"%").Find(&userList)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return []ResUser{}, result.Error
	}
	return userList, nil
}

func SetStyle(userId string, style string) error {
	result := db.GetMysqlDB().Table("users").Where("id = ?", userId).Update("style", style)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateUser(userId string, userInfo User) error {
	// 使用 struct 更新多个属性，只会更新其中有变化且为非零值的字段
	result := db.GetMysqlDB().Table("users").Where("id = ?", userId).Updates(userInfo)
	// 如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return result.Error
	}
	return nil
}
