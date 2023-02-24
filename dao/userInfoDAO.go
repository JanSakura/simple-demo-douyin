package dao

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	//空指针
	errNullPtr = errors.New("Error:Nullptr")
	//空用户列表
	errEmptyUserList = errors.New("UserList:NULL")
)

// UserInfoDAO 抽象的空用户信息结构体，具体的信息
type UserInfoDAO struct {
}

var (
	userInfoDAO  *UserInfoDAO
	userInfoOnce sync.Once
)

func NewUserInfoDao() *UserInfoDAO {
	userInfoOnce.Do(func() {
		userInfoDAO = new(UserInfoDAO)
	})
	return userInfoDAO
}

// InquiryUserInfoByID  根据ID查询用户信息
func (u *UserInfoDAO) InquiryUserInfoByID(userId int64, userInfo *models.UserInfo) error {
	if userInfo == nil {
		return errNullPtr
	}
	//models.GlobalDB.Where("id=?",userId).First(userInfo)
	models.GlobalDB.Where("id=?", userId).Select([]string{"id", "name", "follow_count", "is_follow"}).First(userInfo)
	//SQL执行失败
	if userInfo.Id == 0 {
		return errors.New("不存在该用户")
	}
	return nil
}

// AddUserInfo 向UserInfo表添加用户信息
func (u *UserInfoDAO) AddUserInfo(userInfo *models.UserInfo) error {
	if userInfo == nil {
		return errNullPtr
	}
	return models.GlobalDB.Create(userInfo).Error
}

func (u *UserInfoDAO) IsUserExistByID(id int64) bool {
	var userInfo models.UserInfo
	if err := models.GlobalDB.Where("id=?", id).Select("id").First(&userInfo).Error; err != nil {
		log.Println("IsUserExistByID err:", err)
	}
	if userInfo.Id == 0 {
		return false
	}
	return true
}

// AddUserFollow 添加用户关注者
func (u *UserInfoDAO) AddUserFollow(userId int64, followId int64) error {
	return models.GlobalDB.Transaction(func(tx *gorm.DB) error {
		//添加用户的关注数
		if err := tx.Exec("update user_infos set follow_count=follow_count+1 where id=?", userId).Error; err != nil {
			return err
		}
		//添加被关注者的粉丝数
		if err := tx.Exec("update user_infos set follower_count=follower_count+1 where id=?", followId).Error; err != nil {
			return err
		}
		//向关系表添加关系
		if err := tx.Exec("insert into `user_relations` (`user_info_id`,`follow_id`) values (?,?)", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}

// CancelUserFollow 取消用户关注
func (u *UserInfoDAO) CancelUserFollow(userId int64, followId int64) error {
	return models.GlobalDB.Transaction(func(tx *gorm.DB) error {
		//和添加是相反的逻辑,并且数量>0
		if err := tx.Exec("update user_infos set follow_count=follow_count-1 where follow_count>0 and id=?", userId).Error; err != nil {
			return err
		}
		if err := tx.Exec("update user_infos set follower_count=follower_count-1 where follower_count>0 id=?", followId).Error; err != nil {
			return err
		}
		if err := tx.Exec("delete from `user_relations` where user_info_id=? and follow_id=?", userId, followId).Error; err != nil {
			return err
		}
		return nil
	})
}

// InquiryFollowListByUserId 查询用户的关注者列表
func (u *UserInfoDAO) InquiryFollowListByUserId(userId int64, userList *[]*models.UserInfo) error {
	if userList == nil {
		return errNullPtr
	}
	if err := models.GlobalDB.Raw("select u.* from user_relations r,user_infos u where r.user_info_id =? and r.follow_id =u.id", userId).Scan(userList).Error; err != nil {
		return err
	}
	if len(*userList) == 0 || (*userList)[0].Id == 0 {
		return errEmptyUserList
	}
	return nil
}

// InquiryFollowerListByUserId 查询用的的粉丝列表
func (u *UserInfoDAO) InquiryFollowerListByUserId(userId int64, userList *[]*models.UserInfo) error {
	if userList == nil {
		return errEmptyUserList
	}
	if err := models.GlobalDB.Raw("select u.* from user_relations r,user_infos u where r.user_info_id=u.id and r.followe_id=? ", userId).Scan(userList).Error; err != nil {
		return err
	}
	return nil
}
