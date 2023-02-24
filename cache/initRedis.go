package cache

import (
	"context"
	"fmt"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/go-redis/redis"
)

// 传递上下文消息，在不同的goroutine之间同步，background作为初始的上下文向下传递
var ctx = context.Background()
var rdb *redis.Client

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", models.Info.RDB.IP, models.Info.RDB.Port),
		Password: "",                       //用的Windows本机安装的redis默认没密码
		DB:       models.Info.RDB.Database, //使用的数据库根据配置
	})
}
