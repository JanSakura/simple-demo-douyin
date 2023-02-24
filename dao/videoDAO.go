package dao

import (
	"errors"
	"github.com/JanSakura/simple-demo-douyin/models"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type VideoDAO struct { //空结构体
}

var (
	videoDAO  *VideoDAO
	videoOnce sync.Once
)

func NewVideoDAO() *VideoDAO {
	videoOnce.Do(func() {
		videoDAO = new(VideoDAO)
	})
	return videoDAO
}

//全部写成Hook的形式

// AddVideoDao 视频和userinfo有多对一的关系，所以传入的Video参数一定要进行id的映射处理
func (v *VideoDAO) AddVideoDao(video *models.Video) error {
	if video == nil {
		return errors.New("AddVideo:nullptr")
	}
	return models.GlobalDB.Create(video).Error
}

// QueryVideoByVideoId 根据视频ID返回对应的视频
func (v *VideoDAO) QueryVideoByVideoId(videoId int64, video *models.Video) error {
	if video == nil {
		return errors.New("QueryVideoByVideoID video:nullptr")
	}
	return models.GlobalDB.Where("id=?", videoId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url",
			"favorite_count", "comment_count", "is_favorite", "title"}).First(video).Error
}

func (v *VideoDAO) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId() count:nullptr")
	}
	return models.GlobalDB.Model(&models.Video{}).Where("user_info_id=?", userId).Count(count).Error
}

func (v *VideoDAO) QueryVideoListByUserId(userId int64, videoList *[]*models.Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByUserId() VideoList:nullptr")
	}
	return models.GlobalDB.Where("user_info_id=?", userId).
		Select([]string{"id", "user_info_id", "play_url", "cover_url",
			"favorite_count", "comment_count", "is_favorite", "title"}).Find(videoList).Error
}

// QueryVideoListByLatestTimeAndLimitNum 按投稿时间倒序返回视频列表，并根据limit数目，返回相应的最大视频数
func (v *VideoDAO) QueryVideoListByLatestTimeAndLimitNum(limitNum int, latestTime time.Time, videoList *[]*models.Video) error {
	if videoList == nil {
		return errors.New("QueryVideoListByLatestTimeAndLimitNum() videoList:nullptr")
	}
	return models.GlobalDB.Model(&models.Video{}).Where("created_at<?", latestTime).
		Order("created_at").Limit(limitNum).
		Select([]string{"id", "user_info_id", "play_url", "cover_url",
			"favorite_count", "comment_count", "is_favorite", "title", "created_at", "updated_at"}).Find(videoList).Error
}

// AddFavorByUserIdAndVideoId 根据用户和视频ID，增加一个赞
func (v *VideoDAO) AddFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	return models.GlobalDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("update videos set favorite_count=favorite_count+1 where id=?", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("insert into `user_favor_videos`(`user_info_id`,`video_id`) values (?,?)", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// SubFavorByUserIdAndVideoId 减少一个赞
func (v *VideoDAO) SubFavorByUserIdAndVideoId(userId int64, videoId int64) error {
	return models.GlobalDB.Transaction(func(tx *gorm.DB) error {
		//-1前需要先判断是否合法,不能是负数
		if err := tx.Exec("UPDATE videos SET favorite_count=favorite_count-1 WHERE id = ? AND favorite_count>0", videoId).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM `user_favor_videos`  WHERE `user_info_id` = ? AND `video_id` = ?", userId, videoId).Error; err != nil {
			return err
		}
		return nil
	})
}

// QueryFavorVideoListByUserId 查询用户点赞的视频列表
func (v *VideoDAO) QueryFavorVideoListByUserId(userId int64, videoList *[]*models.Video) error {
	if videoList == nil {
		return errors.New("QueryFavorVideoListByUserId videoList:nullptr")
	}
	//左连接查询,自己写SQL语句
	if err := models.GlobalDB.Raw("select v.* from user_favor_video u,video v where u.user_info_id = ? and u.video_id = v.id", userId).
		Scan(videoList).Error; err != nil {
		return err
	}
	//如果列表长度，或第一个Video的id为0，则为空列表
	if len(*videoList) == 0 || (*videoList)[0].Id == 0 {
		return errors.New("点赞列表：空")
	}
	return nil
}

// IsVideoExistById 判断视频是否存在
func (v *VideoDAO) IsVideoExistById(id int64) bool {
	var video models.Video
	if err := models.GlobalDB.Where("id=?", id).Select("id").First(&video).Error; err != nil {
		log.Println("判断视频存在的错误：", err)
	}
	if video.Id == 0 {
		return false
	}
	return true
}
