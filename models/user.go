package models

// UserLogin 用户登录结构体,对应MySQL数据库的user表，和UserInfo属于一对一关系
type UserLogin struct {
	Id         int64    `gorm:"primary_key"`
	Username   string   `gorm:"primary_key"`
	Password   string   `gorm:"size:200;notnull"` //密码要用盐值加密
	UserInfo   UserInfo //必须写上，不然查不到UserInfo，会变成：user_infos
	UserInfoId int64    //gorm建表外键关联
}

// UserInfo 信息表
type UserInfo struct {
	Id            int64       `json:"id" gorm:"id,omitempty"`                         //用户ID
	Name          string      `json:"name" gorm:"name,omitempty"`                     //用户名
	FollowCount   int64       `json:"follow_count" gorm:"follow_count,omitempty"`     //关注总数
	FollowerCount int64       `json:"follower_count" gorm:"follower_count,omitempty"` //粉丝总数
	IsFollow      bool        `json:"is_follow" gorm:"is_follow,omitempty"`           //是否关注true已关注
	User          *UserLogin  `json:"-"`                                              //用户与账号密码之间的一对一
	Videos        []*Video    `json:"-"`                                              //用户与投稿视频的一对多
	Follows       []*UserInfo `json:"-" gorm:"many2many:user_relations;"`             //用户之间的多对多
	FavorVideos   []*Video    `json:"-" gorm:"many2many:user_favor_videos;"`          //用户与点赞视频之间的多对多
	Comments      []*Comment  `json:"-"`                                              //用户与评论的一对多
}
