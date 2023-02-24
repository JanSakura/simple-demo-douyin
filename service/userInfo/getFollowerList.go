package userInfo

import (
	"github.com/JanSakura/simple-demo-douyin/cache"
	"github.com/JanSakura/simple-demo-douyin/dao"
	"github.com/JanSakura/simple-demo-douyin/models"
)

type FollowerList struct {
	UserList []*models.UserInfo `json:"user_list"`
}

func QueryFollowerList(userId int64) (*FollowerList, error) {
	return NewQueryFollowerListFlow(userId).Do()
}

type QueryFollowerListFlow struct {
	userId int64

	userList []*models.UserInfo

	*FollowerList
}

func NewQueryFollowerListFlow(userId int64) *QueryFollowerListFlow {
	return &QueryFollowerListFlow{userId: userId}
}

func (q *QueryFollowerListFlow) Do() (*FollowerList, error) {
	var err error
	if err = q.checkNum(); err != nil {
		return nil, err
	}
	if err = q.prepareData(); err != nil {
		return nil, err
	}
	if err = q.packData(); err != nil {
		return nil, err
	}
	return q.FollowerList, nil
}

func (q *QueryFollowerListFlow) checkNum() error {
	if !dao.NewUserInfoDao().IsUserExistByID(q.userId) {
		return ErrUserNotExist
	}
	return nil
}

func (q *QueryFollowerListFlow) prepareData() error {

	err := dao.NewUserInfoDao().InquiryFollowerListByUserId(q.userId, &q.userList)
	if err != nil {
		return err
	}
	//填充is_follow字段
	for _, v := range q.userList {
		v.IsFollow = cache.NewCacheSet().GainUserRelationState(q.userId, v.Id)
	}
	return nil
}

func (q *QueryFollowerListFlow) packData() error {
	q.FollowerList = &FollowerList{UserList: q.userList}

	return nil
}
