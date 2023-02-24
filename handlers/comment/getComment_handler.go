package comment

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/comment"
	"github.com/JanSakura/simple-demo-douyin/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type listResponse struct {
	models.ResponseStatus
	*comment.List
}
type favorVideoListResponse struct {
	models.ResponseStatus
	*video.FavorList
}

func GetCommentListHandler(c *gin.Context) {
	newProxyCommentListHandler(c).Do()
}

type proxyCommentListHandler struct {
	*gin.Context

	videoId int64
	userId  int64
}

func newProxyCommentListHandler(context *gin.Context) *proxyCommentListHandler {
	return &proxyCommentListHandler{Context: context}
}

func (p *proxyCommentListHandler) Do() {
	//解析参数
	if err := p.parseNum(); err != nil {
		p.SendError(err.Error())
		return
	}
	//正式调用
	commentList, err := comment.QueryCommentList(p.userId, p.videoId)
	if err != nil {
		p.SendError(err.Error())
		return
	}
	//成功返回
	p.SendOk(commentList)
}

func (p *proxyCommentListHandler) parseNum() error {
	rawUserId, _ := p.Get("user_id")
	userId, ok := rawUserId.(int64)
	if !ok {
		return errors.New("userId解析出错")
	}
	p.userId = userId

	rawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(rawVideoId, 10, 64)
	if err != nil {
		return err
	}
	p.videoId = videoId

	return nil
}

func (p *proxyCommentListHandler) SendError(msg string) {
	p.JSON(http.StatusOK, favorVideoListResponse{
		ResponseStatus: models.ResponseStatus{StatusCode: 1, StatusMsg: msg}})
}

func (p *proxyCommentListHandler) SendOk(commentList *comment.List) {
	p.JSON(http.StatusOK, listResponse{ResponseStatus: models.ResponseStatus{StatusCode: 0},
		List: commentList,
	})
}
