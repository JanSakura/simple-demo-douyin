package video

import (
	"errors"
	"fmt"
	"github.com/JanSakura/simple-demo-douyin/cache"
	"github.com/JanSakura/simple-demo-douyin/dao"
	"github.com/JanSakura/simple-demo-douyin/models"
	"log"
	"path/filepath"
	"time"
)

// GetVideoFileUrl 返回视频文件的路径
func GetVideoFileUrl(fileName string) string {
	base := fmt.Sprintf("http://%s:%d/static/%s", models.Info.IP, models.Info.Port, fileName)
	return base
}

// NewVideoFileName 返回userId+视频数得到的文件名
func NewVideoFileName(userId int64) string {
	var count int64
	err := dao.NewVideoDAO().QueryVideoCountByUserId(userId, &count)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d", userId, count)
}

// FillVideoListField 添加视频的作者信息，通过userId判断作者与视频的关系
func FillVideoListField(userId int64, videos *[]*models.Video) (*time.Time, error) {
	listLen := len(*videos)
	if videos == nil || listLen == 0 {
		return nil, errors.New("FillVideoListField videos:null")
	}
	daos := dao.NewUserInfoDao()
	p := cache.NewCacheSet()                     //redis 缓存
	latestTime := (*videos)[listLen-1].CreatedAt //最近的投稿时间
	//MySQL添加作者信息,向Redis缓存更新is_follow状态
	for i := 0; i < listLen; i++ {
		var userInfo models.UserInfo
		err := daos.InquiryUserInfoByID((*videos)[i].UserInfoId, &userInfo)
		if err != nil {
			continue
		}
		userInfo.IsFollow = p.GainUserRelationState(userId, userInfo.Id) //根据Redis缓存更新是否被点赞
		(*videos)[i].Author = userInfo
		//userId>0 认为是登录状态,其余都是未登录
		if userId > 0 {
			(*videos)[i].IsFavorite = p.GainVideoFavorState(userId, (*videos)[i].Id)
		}
	}
	return &latestTime, nil
}

// SaveVideoCover 将视频的一帧当做封面保存到static,check控制是否打印FFmpeg命令
func SaveVideoCover(name string, check bool) error {
	cov := NewVideoToCover()
	if check {
		cov.Check()
	}
	//生成图片的传入和传出路径
	cov.InputPath = filepath.Join(models.Info.StaticSourcePath, name+defaultVideoSuffix)
	cov.OutPutPath = filepath.Join(models.Info.StaticSourcePath, name+defaultCoverSuffix)
	cov.FrameCount = 1
	cmdStr, err := cov.GainCmdString()
	if err != nil {
		return err
	}
	return cov.ExecCmd(cmdStr)
}
