package video

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
)

type favorVideoListResponse struct {
	models.ResponseStatus
	*video.FavorList
}

type proxyFavorVideoListHandler struct {
	*gin.Context
	userId int64
}

func newProxyFavorVideoListHandler(c *gin.Context) *proxyFavorVideoListHandler {
	return &proxyFavorVideoListHandler{Context: c}
}

func GetFavorVideoListHandler(c *gin.Context) {
	newProxyFavorVideoListHandler(c).Do()
}

func (p *proxyFavorVideoListHandler) Do() {
	//解析参数
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	//正式调用
	favorVideoList, err := video.QueryFavorVideoList(p.userId)
	if err != nil {
		p.SendError(err.Error())
		return
	}

	//成功返回
	p.SendOk(favorVideoList)
}

func (p *proxyFavorVideoListHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId
	return nil
}

func (p *proxyFavorVideoListHandler) SendError(msg string) {
	p.JSON(http.StatusOK, favorVideoListResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 1, StatusMsg: msg}})
}

func (p *proxyFavorVideoListHandler) SendOk(favorList *video.FavorList) {
	p.JSON(http.StatusOK, favorVideoListResponse{ResponseStatus: models.ResponseStatus{StatusCode: 0},
		FavorList: favorList,
	})
}
