package video

import (
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

// 声明视频和图片格式
var (
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
	}
	pictureIndexMap = map[string]struct{}{
		".jpg": {},
		".bmp": {},
		".png": {},
		".svg": {},
	}
)

// PublishVideoHandler 发布视频，并截取一帧画面作为封面
func PublishVideoHandler(c *gin.Context) {
	//准备参数
	rawId, _ := c.Get("user_id")

	userId, ok := rawId.(int64)
	if !ok {
		publishVideoError(c, "解析UserId出错")
		return
	}

	title := c.PostForm("title")

	form, err := c.MultipartForm()
	if err != nil {
		publishVideoError(c, err.Error())
		return
	}

	//支持多文件上传
	files := form.File["data"]
	for _, file := range files {
		suffix := filepath.Ext(file.Filename)    //得到后缀
		if _, ok := videoIndexMap[suffix]; !ok { //判断是否为视频格式
			publishVideoError(c, "不支持的视频格式")
			continue
		}
		//根据userId得到唯一的文件名
		name := video.NewVideoFileName(userId)
		filename := name + suffix
		savePath := filepath.Join("./static", filename)
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			publishVideoError(c, err.Error())
			continue
		}
		//截取一帧画面作为封面
		err = video.SaveVideoCover(name, true)
		if err != nil {
			publishVideoError(c, err.Error())
			continue
		}
		//数据库持久化
		err := video.PostVideo(userId, filename, name+video.GainCoverSuffix(), title)
		if err != nil {
			publishVideoError(c, err.Error())
			continue
		}
		publishVideoOk(c, file.Filename+"上传成功")
	}
}

func publishVideoError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.ResponseStatus{
		StatusCode: 1,
		StatusMsg:  msg})
}

func publishVideoOk(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, models.ResponseStatus{StatusCode: 0, StatusMsg: msg})
}
