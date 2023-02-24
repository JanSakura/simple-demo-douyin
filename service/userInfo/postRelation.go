package userInfo

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/cache"
	"github.com/JanSakura/simple-demo-douyin/dao"
)

const (
	FOLLOW = 1
	CANCEL = 2
)

var (
	ErrIvdAct    = errors.New("未定义操作")
	ErrIvdFolUsr = errors.New("关注用户不存在")
)

func PostFollowAction(userId, userToId int64, actionType int) error {
	return NewPostFollowActionFlow(userId, userToId, actionType).Do()
}

type PostFollowActionFlow struct {
	userId     int64
	userToId   int64
	actionType int
}

func NewPostFollowActionFlow(userId int64, userToId int64, actionType int) *PostFollowActionFlow {
	return &PostFollowActionFlow{userId: userId, userToId: userToId, actionType: actionType}
}

func (p *PostFollowActionFlow) Do() error {
	var err error
	if err = p.checkNum(); err != nil {
		return err
	}
	if err = p.publish(); err != nil {
		return err
	}
	return nil
}

func (p *PostFollowActionFlow) checkNum() error {
	//由于userId是经过乐token鉴权故不需要check，只需要检查userToId
	if !dao.NewUserInfoDao().IsUserExistByID(p.userToId) {
		return ErrIvdFolUsr
	}
	if p.actionType != FOLLOW && p.actionType != CANCEL {
		return ErrIvdAct
	}
	//自己不能关注自己
	if p.userId == p.userToId {
		return ErrIvdAct
	}
	return nil
}

func (p *PostFollowActionFlow) publish() error {
	userDAO := dao.NewUserInfoDao()
	var err error
	switch p.actionType {
	case FOLLOW:
		err = userDAO.AddUserFollow(p.userId, p.userToId)
		//更新redis的关注信息
		cache.NewCacheSet().UpdateUserRelationState(p.userId, p.userToId, true)
	case CANCEL:
		err = userDAO.CancelUserFollow(p.userId, p.userToId)
		cache.NewCacheSet().UpdateUserRelationState(p.userId, p.userToId, false)
	default:
		return ErrIvdAct
	}
	return err
}
