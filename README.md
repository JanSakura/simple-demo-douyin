# 一、项目介绍

> 

项目基本信息：实现基础功能，喜欢列表，用户评论，关系列表

Github地址：https://github.com/JanSakura/simple-demo-douyin

# 二、项目分工

> 好的团队协作可以酌情加分哟～请组长和组员做好项目分工与监督。

| **团队成员** | **主要贡献**                 |
| ------------ | ---------------------------- |
| 队长本人     | 负责开发，文档编写，讨论结构 |
| 队员         | debug调试，演示，讨论结构    |
|              |                              |

# 三、项目实现

### 3.1 技术选型与相关开发文档

> 可以补充场景分析环节，明确要解决的问题和前提假设，比如按当前选型和架构总体预计需要xxx存储空间，xxx台服务器......。

#### 项目采用的技术、应用：

数据库：MySQL 8、Redis 5

语言，框架：Go语言，Gin、Gorm框架

外部应用：ffmpeg 5.12，gcc 8.1

#### 场景分析：

##### 1 如何存用户的核心信息：密码。

选定了固定盐值，直接写在源代码里，并使用SHA256。

使用动态盐值，就需要将其存在数据库里，一旦被拖库和不加salt没什么区别；SHA-1虽然短，但彩虹表已相当完善，所以采用了SHA-256。

##### 2 用户的关系如何存储

采用普通的MySQL存关系，在人数少时还没什么，一旦人数变大，数据就会膨胀的特别厉害；并且MySQL是要IO查询的，所以采用Redis存储，相应的可以类比到用户与视频的关系，视频和评论的关系。

##### 3 使用Redis的什么数据类型存储关系

原本预计使用HASH这个字典类型，后来参考了网上的说法，采用了Set。Set可以交叉并集合操作，可以更好地区分用户关系。比如粉丝集，关注集，两个的交集是互粉（可以设定为好友关系），差集即不在两个集合里的都是无关的人。

##### 4 如何确认用户状态

采用了JWT标准，通过ExpireTime来判断，旧的StandClaims被废弃，采用RegisteredClaims。（比较明显的是ExpiresAt从int64格式变成了 NumericDate即time Time格式）

##### 5 处理实时操作

采用Redis处理点赞等即时性操作，先给出反应，同时存储到MySQL。

#### 实现方法选用

##### 空结构体还是空接口

*空接口，在传数据的情况下，可能出现并不是真**nil**的情况，因为Go判断nil是看前8字节，而空接口前8字节会存类型，后8字节为nil，会出现误判*

##### Go的指针与引用

在返回结构体的地址时，受C++影响，认为本质都是一样的，不过就是const的区别；

但Go的官方文档说明传递上只有值传递，Go是通过类似于引用的方式完成的

所以在使用中要注意

### 3.2 架构设计

```Go
│  .gitattributes
│  .gitignore
│  go.mod
│  go.sum
│  main.go
│  README.md
│
├─cache
│      cacheSet.go
│      initRedis.go
│
├─dao
│      commentDAO.go
│      messageDAO.go
│      userInfoDAO.go
│      userLoginDAO.go
│      videoDAO.go
│
├─handlers
│  │  getUserId.go
│  │  jwt.go
│  │
│  ├─comment
│  │      getComment_handler.go
│  │      postComment_handler.go
│  │
│  ├─userInfo
│  │      getFollowerList_handler.go
│  │      getFollowList_handler.go
│  │      postRelation_handler.go
│  │      userInfo_handler.go
│  │
│  ├─userLoginRegister
│  │      genPassword_handler.go
│  │      postUserLogin_handler.go
│  │      postUserRegister_handler.go
│  │
│  └─video
│          feedVideoList_handler.go
│          getFavorVideoList_handler.go
│          getVideoList_handler.go
│          postFavorState_handler.go
│          postPublishVideo_handler.go
│
├─Init
│      config.toml
│      initDB.go
│
├─lib
│      .gitattributes
│      ffmpeg.exe
│      ffprobe.exe
│
├─models
│      comment.go
│      configDB.go
│      message.go
│      responseStatus.go
│      user.go
│      video.go
│
├─routers
│      routers.go
│
├─service
│  ├─comment
│  │      commentUtil.go
│  │      getComment.go
│  │      postComment.go
│  │
│  ├─userInfo
│  │      getFollowerList.go
│  │      getFollowList.go
│  │      postRelation.go
│  │
│  ├─userLoginRegister
│  │      postUserLogin.go
│  │      postUserRegister.go
│  │
│  └─video
│          feed_videoList.go
│          ffmpeg.go
│          getFavorVideoList.go
│          getVideoList.go
│          postFavorState.go
│          postPublishVideo.go
│          videoUtil.go
│
└─static
```

![img](https://qofh5bjs1x.feishu.cn/space/api/box/stream/download/asynccode/?code=Y2E4YWIzYmQxNzdlN2MzYmFmNTI1ZWVmMWVhNDc3YTVfOFFHeVJ0cWt6azlBT0NVWUd1ekEzWFZ4RERHY3pHR0NfVG9rZW46Ym94Y252RGFrTWxDclQwamZxa2Y5d2EyRUdiXzE2Nzk4MzIxNTU6MTY3OTgzNTc1NV9WNA)

### 3.3 项目代码介绍

- 系统入口main.go
- routers里InitRouter加载路由和SQL数据库
- Init初始化MySQL数据库，其中全局配置文件采用toml格式
- 各路由启动对应的Handler， 并做中间件检测，比如是否需要JWT检测，像feed接口要能直接播放视频，无需用户登录而产生的token
- Handler调用对应的service功能，解析上层的数据，对上下层的数据做判断，和简单的数据操作，比如SHA256加密
- service调用MySQL的dao ，有需要的使用Redis做缓存和操作，对得到的数据进行处理，还可能处理与系统的交互，比如使用ffmpeg命令，
- dao抽象操作的数据结构体，并返回对应models的CRUD结果，不做检测和处理
- cache初始化Redis，抽象数据类型，并实现Redis的数据库操作，也不做检测和处理

- ​     **代表性代码展示：**

- ####  建表结构体

-  用Gorm自动建表，需要定义好外键关系，不然容易一开始就被Automigrate被卡住。

```Go
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
```

####     数据库初始化日志配置

可以自定义一些日志输出格式，和数据库使用设置

```Go
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
```

####     JWT创建

新的JWT，弃用了StandardClaims，因为时间校对问题，不再使用int64格式，而是采用NumericDate即time Time。自己在转化中需要注意int和time Time

```Go
// 用于签名的字符串
var jwtKey = []byte("DouyinKey") //jwt的秘钥
// myClaims 创建Claim
type myClaims struct {
   UserId               int64
   jwt.RegisteredClaims //最新版的StandardClaims已废弃
}

// GenerateToken 生成JWT:封装生成Token
func GenerateToken(user models.UserLogin) (string, error) {
   //Token过期时间，如24小时
   const TokenExpireDuration = time.Hour * 24
   //创建自己的声明
   claims := &myClaims{
      UserId: user.UserInfoId, //自定义字段
      RegisteredClaims: jwt.RegisteredClaims{
         ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
         IssuedAt:  jwt.NewNumericDate(time.Now()), //DefaultClaims是int64格式,即time.now().Unix()
         Issuer:    "douyin_demo",                  //签发人
         Subject:   "GenToken",                     //主题
      }}
   //使用指定的签名方法创建签名对象,如SHA256
   token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
   //使用指定的secret签名，获得完整编码后的字符串token
   tokenStr, err := token.SignedString(jwtKey)
   if err != nil {
      return "", err
   }
   return tokenStr, nil
}
```

- ####     JWT超时比较

-  JWT的NumericDate虽然是time Time，但并不能直接使用Go标准的时间函数，要注意形参bool值的选择

- ```Go
  //token超时
  //通过ExpiresAt值判断是否过期，如果没有设置，则返回true，没有过期
  if !tokenMsg.VerifyExpiresAt(time.Now(), false) {
     c.JSON(http.StatusOK, models.ResponseStatus{
        StatusCode: http.StatusPaymentRequired,
        StatusMsg:  "token过期",
     })
     c.Abort()
     return
  }
  ```

#### 简单的密码加盐

```Go
func SHA256(s string) string {
   //SHA256 口令
   auth := sha256.New()
   auth.Write([]byte(s))
   return hex.EncodeToString(auth.Sum(nil))
}
func SHAPassword(c *gin.Context) {
   password := c.Query("password")
   if password == "" {
      password = c.PostForm("password")
   }
   salt := "salt" //简单盐值
   c.Set("password", SHA256(password+salt))
   c.Next()
}
```

#### 具体service结构体操作

操作使用结构体指针，后面的全是针对结构体指针的接口函数

```Go
func PostFavorState(userId, videoId, actionType int64) error {
   return NewPostFavorStateFlow(userId, videoId, actionType).Do()
}
type PostFavorStateFlow struct {
   userId     int64
   videoId    int64
   actionType int64
}
func NewPostFavorStateFlow(userId, videoId, action int64) *PostFavorStateFlow {
   return &PostFavorStateFlow{
      userId:     userId,
      videoId:    videoId,
      actionType: action,
   }
}
```

#### 使用Redis和MySQL操作点赞

```Go
// AddOperation 点赞操作
func (p *PostFavorStateFlow) AddOperation() error {
   //视频点赞数目+1
   err := dao.NewVideoDAO().AddFavorByUserIdAndVideoId(p.userId, p.videoId)
   if err != nil {
      return errors.New("不要重复点赞")
   }
   //对应的用户是否点赞的映射状态更新
   cache.NewCacheSet().UpdateVideoFavorStateByUserIdAndVideoId(p.userId, p.videoId, true)
   return nil
}
```

#### ffmpeg执行

可以采用C语言和Go的执行方式

```Go
///*
//#include<stdlib.h>
//int Cmd(const char* cmd){
// return system(cmd);
//}
//*/
//import "C"

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
   // return errors.New("ffmpeg生成封面失败")
   //}
   //
   return nil
}
```

#### DAO层和cache层空结构体设计

使用对应的model层表的空结构体，并使用指针，和考虑锁问题

```Go
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
```

对应的Redis没有表，只有数据类型，所以只要将数据类型定义为空结构体

```Go
type CacheSet struct {
}
var cacheSet CacheSet
```

# 四、测试结果

功能测试：

比如调用系统的ffmpeg命令，是粘合字符串到一起，不知道会是什么情况，只能是看结果情况，不利于修改

```Go
// 粘合string的函数
func Test_paramJoin(s1, s2 string) string {
   return fmt.Sprintf(" %s %s ", s1, s2) //注意空格
}
func (v *VideoToCover) Test_GainCmdString() (retStr string, err error) {
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
```

# 五、Demo 演示视频 （必填）

暂时无法在飞书文档外展示此内容

# 六、项目总结与反思

### 目前仍存在的问题

1 Gorm学的浅，会基本的增删改查

2 函数和变量命名，为了省事刚开始全是大写开头，导致多了后调用会有特别多东西提示，应该先预期好实现的功能，区分出能否继承，和外界调用。

### 已识别出的优化项

1 采用Redis，提升了响应，和大数据存储查询修改。

### 架构演进的可能性

​      router和所需的handler、service都是一一对应的，也可以单独部署。

### 项目过程中的反思与总结

​       1 项目开始框架不清楚，导致花了很多冤枉时间：比如JWT生成token和解析token，是伴随着其他handler前后的，因此刚开始混着来。

2 命名比较混乱。因为DAO层的函数是要能被上层调用的，而基础操作基本就增删改查，很容易前面的一样，调用时看着眼花，所以想着采用不同的近义词英文替代，部分的service也出现这种情况，导致名字挺混乱，一下子有的看不出是干什么的。后面在写handler时，就相对统一了规范，比如是get请求的，就统一以Get开头，再加上对应的router。

# 七、其他补充资料