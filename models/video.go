package models

import "time"

// MaxVideoNum 视频列表，最大数
const MaxVideoNum = 30

type Video struct {
	Id            int64       `json:"id,omitempty"`              //视频唯一标识
	UserInfoId    int64       `json:"-"`                         //查询状态的一个中间值
	Author        UserInfo    `json:"author,omitempty" gorm:"-"` //这里应该是作者对视频的一对多的关系，而不是视频对作者，故gorm不能存他，但json需要返回它
	PlayUrl       string      `json:"play_url,omitempty"`        //视频播放地址
	CoverUrl      string      `json:"cover_url,omitempty"`       //封面地址
	FavoriteCount int64       `json:"favorite_count,omitempty"`  //点赞数
	CommentCount  int64       `json:"comment_count,omitempty"`   //评论数
	IsFavorite    bool        `json:"is_favorite,omitempty"`     //是否点赞
	Title         string      `json:"title,omitempty"`           //视频标题
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}
type FeedVideoList struct {
	Videos   []*Video `json:"video_list,omitempty"` //视频list
	NextTime int64    `json:"next_time,omitempty"`  //本次返回的视频中最早的时间，作为下次请求的latestTime
}
