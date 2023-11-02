package dao

import (
	"errors"
	"github.com/IRONICBo/QiYin_BE/internal/conn/db"
	"github.com/IRONICBo/QiYin_BE/internal/utils"
	"log"
)

// Collection 表的结构。
type Collection struct {
	Id      int64  //自增主键
	UserId  string //点赞用户id
	VideoId int64  //视频id
	Cancel  int8   //是否收藏
}

// TableName 修改表名映射
func (Collection) TableName() string {
	return "collections"
}

// GetCollectionUserIdList 根据videoId获取点赞userId
func GetCollectionUserIdList(videoId int64) ([]string, error) {
	var collectionUserIdList []string //存所有该视频点赞用户id；
	//查询collections表对应视频id点赞用户，返回查询结果
	err := db.GetMysqlDB().Model(Collection{}).Where(map[string]interface{}{"video_id": videoId, "cancel": utils.IsCollection}).
		Pluck("user_id", &collectionUserIdList).Error
	//查询过程出现错误，返回默认值0，并输出错误信息
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("get collectionUserIdList failed")
	} else {
		//没查询到或者查询到结果，返回数量以及无报错
		return collectionUserIdList, nil
	}
}

// UpdateCollection 根据userId，videoId,actionType点赞或者取消赞
func UpdateCollection(userId string, videoId int64, actionType int32) error {
	//更新当前用户观看视频的点赞状态“cancel”，返回错误结果
	err := db.GetMysqlDB().Model(Collection{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", actionType).Error
	//如果出现错误，返回更新数据库失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	//更新操作成功
	return nil
}

// InsertCollection 插入点赞数据
func InsertCollection(collectionData Collection) error {
	//创建点赞数据，默认为点赞，cancel为0，返回错误结果
	err := db.GetMysqlDB().Model(Collection{}).Create(&collectionData).Error
	//如果有错误结果，返回插入失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert data fail")
	}
	return nil
}

// GetCollectionInfo 根据userId,videoId查询点赞信息
func GetCollectionInfo(userId string, videoId int64) (Collection, error) {
	//创建一条空collection结构体，用来存储查询到的信息
	var collectionInfo Collection
	//根据userid,videoId查询是否有该条信息，如果有，存储在collectionInfo,返回查询结果
	err := db.GetMysqlDB().Model(Collection{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&collectionInfo).Error
	if err != nil {
		//查询数据为0，打印"can't find data"，返回空结构体，这时候就应该要考虑是否插入这条数据了
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return Collection{}, nil
		} else {
			//如果查询数据库失败，返回获取collectionInfo信息失败
			log.Println(err.Error())
			return collectionInfo, errors.New("get collectionInfo failed")
		}
	}
	return collectionInfo, nil
}

// GetCollectionVideoIdList 根据userId查询所属点赞全部videoId
func GetCollectionVideoIdList(userId string) ([]string, error) {
	var collectionVideoIdList []string
	err := db.GetMysqlDB().Model(Collection{}).Where(map[string]interface{}{"user_id": userId, "cancel": utils.IsCollection}).
		Pluck("video_id", &collectionVideoIdList).Error
	if err != nil {
		//查询数据为0，返回空collectionVideoIdList切片，以及返回无错误
		if "record not found" == err.Error() {
			log.Println("there are no collectionVideoId")
			return collectionVideoIdList, nil
		} else {
			//如果查询数据库失败，返回获取collectionVideoIdList失败
			log.Println(err.Error())
			return collectionVideoIdList, errors.New("get collectionVideoIdList failed")
		}
	}
	return collectionVideoIdList, nil
}
