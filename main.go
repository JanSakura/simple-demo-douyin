package main

import (
	"fmt"
	"github.com/JanSakura/simple-demo-douyin/models"
	"github.com/JanSakura/simple-demo-douyin/routers"
)

func main() {
	r := routers.InitRouter()
	err := r.Run(fmt.Sprintf(":%d", models.Info.Port)) //监听端口
	if err != nil {
		return
	}
}
