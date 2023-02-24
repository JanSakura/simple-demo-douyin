package cache

import "fmt"

// CacheSet 用Redis的Set数据类型存关系
// 参考了网上的一些建议，从Hash类型改为Set
type CacheSet struct {
}

var cacheSet CacheSet

const favor = "favor"
const relation = "relation"

func NewCacheSet() *CacheSet {
	return &cacheSet
}

// UpdateVideoFavorStateByUserIdAndVideoId 更新点赞状态
func (s *CacheSet) UpdateVideoFavorStateByUserIdAndVideoId(userId int64, videoId int64, state bool) {
	key := fmt.Sprintf("%s:%d", favor, userId)
	if state {
		rdb.SAdd(key, videoId)
		return
	}
	rdb.SRem(key, videoId)
}

// GainVideoFavorState 获得点赞
func (s *CacheSet) GainVideoFavorState(userId int64, videoId int64) bool {
	key := fmt.Sprintf("%s:%d", favor, userId)
	return rdb.SIsMember(key, videoId).Val()
}

// UpdateUserRelationState 更新用户关系状态，state=true为关注
func (s *CacheSet) UpdateUserRelationState(userId int64, followId int64, state bool) {
	key := fmt.Sprintf("%s:%d", relation, userId)
	if state {
		rdb.SAdd(key, followId)
	}
	rdb.SRem(key, followId)
}

// GainUserRelationState 增加关注状态
func (s *CacheSet) GainUserRelationState(userId int64, followId int64) bool {
	key := fmt.Sprintf("%s:%d", relation, userId)
	return rdb.SIsMember(key, followId).Val()
}
