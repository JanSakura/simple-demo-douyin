package video

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/cache"
	"github.com/JanSakura/simple-demo-douyin/dao"
)

const (
	ADD = 1
	SUB = 2
)

func PostFavorState(userId, videoId, actionType int64) error {
	return NewPostFavorStateFlow(userId, videoId, actionType).Do()
}

type PostFavorStateFlow struct {
	userId     int64
	videoId    int64
	actionType int64
}

func NewPostFavorStateFlow(userId, videoId, action int64) *PostFavorStateFlow {
	return &PostFavorStateFlow{
		userId:     userId,
		videoId:    videoId,
		actionType: action,
	}
}

func (p *PostFavorStateFlow) Do() error {
	var err error
	if err = p.checkNum(); err != nil {
		return err
	}

	switch p.actionType {
	case ADD:
		err = p.AddOperation()
	case SUB:
		err = p.SubOperation()
	default:
		return errors.New("未定义的操作")
	}
	return err
}

// AddOperation 点赞操作
func (p *PostFavorStateFlow) AddOperation() error {
	//视频点赞数目+1
	err := dao.NewVideoDAO().AddFavorByUserIdAndVideoId(p.userId, p.videoId)
	if err != nil {
		return errors.New("不要重复点赞")
	}
	//对应的用户是否点赞的映射状态更新
	cache.NewCacheSet().UpdateVideoFavorStateByUserIdAndVideoId(p.userId, p.videoId, true)
	return nil
}

// SubOperation  取消点赞
func (p *PostFavorStateFlow) SubOperation() error {
	//视频点赞数目-1
	err := dao.NewVideoDAO().SubFavorByUserIdAndVideoId(p.userId, p.videoId)
	if err != nil {
		return errors.New("点赞数目已经为0")
	}
	//对应的用户是否点赞的映射状态更新
	cache.NewCacheSet().UpdateVideoFavorStateByUserIdAndVideoId(p.userId, p.videoId, false)
	return nil
}

func (p *PostFavorStateFlow) checkNum() error {
	if !dao.NewUserInfoDao().IsUserExistByID(p.userId) {
		return errors.New("用户不存在")
	}
	if p.actionType != ADD && p.actionType != SUB {
		return errors.New("未定义的行为")
	}
	return nil
}
