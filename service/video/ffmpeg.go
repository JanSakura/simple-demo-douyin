package video

///*
//#include<stdlib.h>
//int Cmd(const char* cmd){
//	return system(cmd);
//}
//*/
//import "C"
import (
	"errors"
	"fmt"
	"github.com/JanSakura/simple-demo-douyin/models"
	"log"
	"os/exec"
)

type VideoToCover struct {
	InputPath  string
	OutPutPath string
	StartTime  string //开始时间
	KeepTime   string //持续时间
	Filter     string //过滤器
	FrameCount int64
	check      bool
}

var videoToCover VideoToCover

func NewVideoToCover() *VideoToCover {
	return &videoToCover
}

// ffmpeg 参数,ffmpeg [a] [b] -i [c] [d] [e]
// 全局参数a、输入文件参数b,c输入文件名(包括扩展名),输出文件参数d,e输出文件名
// ffmpeg -y -ss 1 -t 5 -i 1.mp4 -c:v copy -c:a copy cut.mp4	截取5s的视频
// ffmpeg -y -i input.mp4 -ss 00:00:00 -t 00:00:01 output_%3d.jpg 截图
const (
	inputVideoPathPara = "-i"        //输入文件
	startTimePara      = "-ss"       //从后面数字的第0秒开始截取
	keepTimePara       = "-t"        //持续时间
	videoFilterPara    = "-vf"       //-filter:v的简写,视频过滤器
	formatToCoverPara  = "-f"        //指定使用avfoundation采集数据
	autoReWritePara    = "-y"        //覆盖输出文件,不询问
	framesPara         = "-frames:v" //获得image的格式
)

// 默认的视频和封面格式,作拼接
var (
	defaultVideoSuffix = ".mp4"
	defaultCoverSuffix = ".jpg"
)

// ChangeVideoSuffix 如果传的不是mp4格式的视频，那么默认视频格式转变为用户传的格式
func ChangeVideoSuffix(suffix string) {
	defaultVideoSuffix = suffix
}

func ChangeCoverSuffix(suffix string) {
	defaultCoverSuffix = suffix
}

func GainCoverSuffix() string {
	return defaultCoverSuffix
}

// 粘合string的函数
func paramJoin(s1, s2 string) string {
	return fmt.Sprintf(" %s %s ", s1, s2) //注意空格
}

func (v *VideoToCover) Check() {
	v.check = true
}

// GainCmdString 生成Ffmpeg的命令语句
func (v *VideoToCover) GainCmdString() (retStr string, err error) {
	if v.InputPath == "" || v.OutPutPath == "" {
		err = errors.New("未指定输入或输出路径")
		return
	}
	retStr = models.Info.FfmpegPath
	retStr += paramJoin(formatToCoverPara, v.InputPath)
	retStr += paramJoin(inputVideoPathPara, "cover")
	if v.Filter != "" {
		retStr += paramJoin(videoFilterPara, v.Filter)
	}
	if v.StartTime != "" {
		retStr += paramJoin(startTimePara, v.StartTime)
	}
	if v.KeepTime != "" {
		retStr += paramJoin(keepTimePara, v.KeepTime)
	}
	if v.FrameCount != 0 {
		retStr += paramJoin(framesPara, fmt.Sprintf("%d", v.FrameCount))
	}
	retStr += paramJoin(autoReWritePara, v.OutPutPath)
	return
}

func (v *VideoToCover) ExecCmd(cmd string) error {
	if v.check {
		log.Println(cmd)
	}
	//创建*exec.Cmd
	cmdState := exec.Command(cmd)
	if err := cmdState.Run(); err != nil {
		return errors.New("ffmpeg生成封面失败")
	}
	//
	////下面是用C语言写的,C的逻辑和语法，Go的形式
	//execStr := C.CString(cmd)
	////C 不会自动释放内存
	//defer C.free(unsafe.Pointer(execStr))
	//status := C.Cmd(execStr)
	//if status != 0 {
	//	return errors.New("ffmpeg生成封面失败")
	//}
	//
	return nil
}
