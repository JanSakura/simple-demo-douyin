package Init

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/JanSakura/simple-demo-douyin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strings"
	"time"
)

// 注意包名，init会与init函数冲突，顾为大写的Init
// 初始化数据库连接配置,配置文件格式toml
func init() {
	if _, err := toml.DecodeFile(
		"C:\\Users\\Administrator\\go\\src\\github.com\\JanSakura\\simple-demo-douyin\\Init\\config.toml",
		&models.Info); err != nil {
		panic(err)
	}
	//去除左右的空格
	strings.Trim(models.Info.Server.IP, " ")
	strings.Trim(models.Info.RDB.IP, " ")
	strings.Trim(models.Info.DB.Host, " ")
}

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		models.Info.DB.Username, models.Info.DB.Password, models.Info.DB.Host, models.Info.DB.Port, models.Info.DB.Database,
		models.Info.DB.Charset, models.Info.DB.ParseTime, models.Info.DB.Loc)
	log.Println(dsn)
	//写sql语句logger配置
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //ioWriter(日志输出的目标,前缀和日志包含的内容)
		logger.Config{
			SlowThreshold:             time.Second,   //慢SQL阈值,1秒
			LogLevel:                  logger.Silent, //日志级别:silent
			IgnoreRecordNotFoundError: true,          //忽略记录未找到(ErrRecordNotFound)错误
			Colorful:                  false,         //不启用彩色日志
		})
	var err error
	models.GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,      //缓存预编译命令
		SkipDefaultTransaction: true,      //禁用默认事务操作
		Logger:                 newLogger, //自定义日志
	})
	if err != nil {
		log.Panicf("初始化MySQL异常：%v", err)
	}
	err = models.GlobalDB.AutoMigrate(&models.UserLogin{},
		&models.UserInfo{},
		&models.Video{},
		&models.Comment{})
	if err != nil {
		panic(err)
	}
}
