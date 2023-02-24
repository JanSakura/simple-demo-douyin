package video

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/handlers"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

//type feed_request struct {
//	last_time int64  //可选参数，限制返回视频的最新投稿时间戳，精确到秒
//	token     string //可选参数,登录用户设置
//}

type feed_response struct {
	models.ResponseStatus
	*models.FeedVideoList
}

func FeedVideoListHandler(c *gin.Context) {
	p := newProxyFeedVideoList(c)
	token, ok := c.GetQuery("token")
	//无登录状态
	if !ok {
		err := p.DoNoToken()
		if err != nil {
			p.FeedVideoListError(err.Error())
		}
		return
	}
	err := p.DoToken(token)
	if err != nil {
		p.FeedVideoListError(err.Error())
	}
}

// proxyFeedVideoList feed流代理列表,用结构体,可新增自定义方法
type proxyFeedVideoList struct {
	*gin.Context
}

// NewProxyFeedVideoList 新的feed流
func newProxyFeedVideoList(c *gin.Context) *proxyFeedVideoList {
	return &proxyFeedVideoList{Context: c}
}

// DoNoToken 未登录状态下的视频流推送
func (p *proxyFeedVideoList) DoNoToken() error {
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, intTime*1e6) //前端返回的时间戳是ms
	}
	videoList, err := video.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListOk(videoList)
	return nil
}

// DoToken 登录状态的视频流推送
func (p *proxyFeedVideoList) DoToken(token string) error {
	//jwt解析token
	claim, err := handlers.ParseToken(token)
	if err == nil {
		//判断是否过期
		if claim.VerifyExpiresAt(time.Now(), false) == false {
			return errors.New("token 超时")
		}
	}
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err != nil {
		latestTime = time.Unix(0, intTime*1e6)
	}
	//service interface
	videoList, err := video.QueryFeedVideoList(claim.UserId, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListOk(videoList)
	return nil
}

// FeedVideoListError feed流错误下的的JSON
func (p *proxyFeedVideoList) FeedVideoListError(err string) {
	p.JSON(http.StatusOK, feed_response{
		ResponseStatus: models.ResponseStatus{
			StatusCode: 1,
			StatusMsg:  err,
		},
	})
}

// FeedVideoListOk 返回FeedVideoList的JSON
func (p *proxyFeedVideoList) FeedVideoListOk(videoList *models.FeedVideoList) {
	p.JSON(http.StatusOK, feed_response{
		ResponseStatus: models.ResponseStatus{
			StatusCode: 0,
		},
		FeedVideoList: videoList,
	})
}
