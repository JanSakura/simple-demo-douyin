package video

import (
	"github.com/JanSakura/simple-demo-douyin/dao"
	"github.com/JanSakura/simple-demo-douyin/models"
	"time"
)

type queryFeedVideoListFlow struct {
	userId     int64
	latestTime time.Time
	videos     []*models.Video
	nextTime   int64
	feedVideo  *models.FeedVideoList
}

func QueryFeedVideoList(userId int64, latestTime time.Time) (*models.FeedVideoList, error) {
	return NewQueryFeedVideoListFlow(userId, latestTime).Do()
}

func NewQueryFeedVideoListFlow(userId int64, latestTime time.Time) *queryFeedVideoListFlow {
	return &queryFeedVideoListFlow{userId: userId, latestTime: latestTime}
}

// Do 生成feed流视频列表和nextTime
func (q *queryFeedVideoListFlow) Do() (*models.FeedVideoList, error) {
	q.checkNum()
	if err := q.prepareTimeDate(); err != nil {
		return nil, err
	}
	if err := q.packFeedVideoList(); err != nil {
		return nil, err
	}
	return q.feedVideo, nil
}

// 传入的参数进行正常化处理
func (q *queryFeedVideoListFlow) checkNum() {
	if q.userId > 0 {
		//userId有效，可以定制登录用户的专属视频推荐
	}
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
}

// 准备下一次feed流的时间
func (q *queryFeedVideoListFlow) prepareTimeDate() error {
	//
	err := dao.NewVideoDAO().QueryVideoListByLatestTimeAndLimitNum(models.MaxVideoNum, q.latestTime, &q.videos)
	if err != nil {
		return err
	}
	//登录状态，则更新该视频的用户点赞状态，不是致命错误不返回，即err=_
	latestTime, _ := FillVideoListField(q.userId, &q.videos)

	//生成时间戳
	if latestTime != nil {
		q.nextTime = (*latestTime).UnixNano() / 1e6
		return nil
	}
	q.nextTime = time.Now().Unix() / 1e6
	return nil
}

// 打包feed流视频列表和nextTime
func (q *queryFeedVideoListFlow) packFeedVideoList() error {
	q.feedVideo = &models.FeedVideoList{
		Videos:   q.videos,
		NextTime: q.nextTime,
	}
	return nil
}
