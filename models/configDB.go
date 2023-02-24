package models

import "gorm.io/gorm"

var Info Config       //config
var GlobalDB *gorm.DB //全局MySQL DB

type Mysql struct {
	Host      string //IP
	Port      int    //端口
	Database  string //数据库名
	Username  string //登录用户名
	Password  string
	Charset   string //编码
	ParseTime bool   `toml:"parse_time"`
	Loc       string //地址方式
}

type Redis struct {
	IP       string
	Port     int
	Database int //Redis的数据库名是数字
}

type Server struct {
	IP   string //服务器IP
	Port int
}

type Path struct {
	FfmpegPath       string `toml:"ffmpeg_path"`
	StaticSourcePath string `toml:"static_source_path"`
}

type Config struct {
	DB     Mysql `toml:"mysql"`
	RDB    Redis `toml:"redis"`
	Server `toml:"server"`
	Path   `toml:"path"`
}
