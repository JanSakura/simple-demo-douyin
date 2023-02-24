package video

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
)

type listResponse struct {
	models.ResponseStatus
	*video.List
}

// ProxyQueryVideoList 代理类
type proxyQueryVideoList struct {
	c *gin.Context
}

func newProxyQueryVideoList(c *gin.Context) *proxyQueryVideoList {
	return &proxyQueryVideoList{c: c}
}

func GetVideoListHandler(c *gin.Context) {
	p := newProxyQueryVideoList(c)
	rawId, _ := c.Get("user_id")
	err := p.DoQueryVideoListByUserId(rawId)
	if err != nil {
		p.QueryVideoListError(err.Error())
	}
}

// DoQueryVideoListByUserId 根据userId字段进行查询
func (p *proxyQueryVideoList) DoQueryVideoListByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}

	videoList, err := video.QueryVideoListByUserId(userId)
	if err != nil {
		return err
	}

	p.QueryVideoListOk(videoList)
	return nil
}

func (p *proxyQueryVideoList) QueryVideoListError(msg string) {
	p.c.JSON(http.StatusOK, listResponse{
		ResponseStatus: models.ResponseStatus{
			StatusCode: 1,
			StatusMsg:  msg,
		}})
}

func (p *proxyQueryVideoList) QueryVideoListOk(videoList *video.List) {
	p.c.JSON(http.StatusOK, listResponse{
		ResponseStatus: models.ResponseStatus{
			StatusCode: 0,
		},
		List: videoList,
	})
}
