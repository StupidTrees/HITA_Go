package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"hita/config"
	"hita/controller"
	"hita/middleware"
	"hita/repository"
	"hita/utils/logger"
	"hita/utils/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//.初始化
func init() {

	err := config.LoadConfig("config.json")
	if err != nil {
		logger.Fatalln("LoadConfig failed:", err)
	}

	logger.Println("load config success")

	//连接到mysql数据库
	err = mysql.InitDB()
	if err != nil {
		logger.Fatalln("InitDB failed:", err)
	} else {
		_ = mysql.DB.AutoMigrate(&repository.User{})
		_ = mysql.DB.AutoMigrate(&repository.Timetable{})
		_ = mysql.DB.AutoMigrate(&repository.TermSubject{})
		_ = mysql.DB.AutoMigrate(&repository.Event{})
		_ = mysql.DB.AutoMigrate(&repository.History{})
		_ = mysql.DB.AutoMigrate(&repository.Article{})
		_ = mysql.DB.AutoMigrate(&repository.UserLikeArticle{})
		_ = mysql.DB.AutoMigrate(&repository.Comment{})
		_ = mysql.DB.AutoMigrate(&repository.UserLikeComment{})
		_ = mysql.DB.AutoMigrate(&repository.Follow{})
		_ = mysql.DB.AutoMigrate(&repository.Image{})
		_ = mysql.DB.AutoMigrate(&repository.Star{})
		_ = mysql.DB.AutoMigrate(&repository.Topic{})
		err = mysql.DB.AutoMigrate(&repository.Vote{})
		_ = mysql.DB.AutoMigrate(&repository.Message{})
		_ = mysql.DB.AutoMigrate(&repository.Info{})
		_ = mysql.DB.AutoMigrate(&repository.Inbox{})
		print(err)
	}

	logger.Println("init db success")

	logger.Println("process init success")
}

func main() {
	//创建一个gin框架的指针
	router := gin.New()

	//设置中间件
	router.Use(gin.Recovery())
	router.Use(middleware.TlsHandler()) //https

	frontRoutes := router.Group("/front")
	{
		frontRoutes.GET("user_num", controller.CountUserNum)
		frontRoutes.GET("latest_version_name", controller.GetLatestVersionName)
		frontRoutes.POST("suggestion",controller.MakeSuggestion)
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/sign_up", controller.SignUp)
		userRoutes.POST("/log_in", controller.LogIn)
	}
	router.GET("/profile/avatar", controller.GetAvatar)
	router.GET("/article/image", controller.GetArticleImage)
	syncRoutes := router.Group("/sync")
	{
		syncRoutes.POST("/sync", controller.Sync)
		syncRoutes.POST("/push", controller.Push)
	}
	router.Use(middleware.JWTAuthMiddleware)

	managerRoutes := router.Group("/manager")
	{
		managerRoutes.GET("check_update",controller.CheckUpdate)
	}
	profileRoutes := router.Group("/profile")
	{
		profileRoutes.POST("upload_avatar", controller.UploadAvatar)
		profileRoutes.GET("get", controller.GetBasicProfile)
		profileRoutes.POST("change_signature", controller.ChangeSignature)
		profileRoutes.POST("change_gender", controller.ChangeGender)
		profileRoutes.POST("change_nickname", controller.ChangeNickname)
		profileRoutes.POST("follow", controller.FollowOrUnFollow)
		profileRoutes.GET("gets", controller.GetUsers)
	}
	articleRoutes := router.Group("/article")
	{
		articleRoutes.POST("create", controller.CreateArticle)
		articleRoutes.POST("create_images", controller.CreateArticleWithImages)
		articleRoutes.GET("gets", controller.GetArticles)
		articleRoutes.GET("get", controller.GetArticle)
		articleRoutes.POST("like", controller.LikeOrUnlike)
		articleRoutes.POST("vote", controller.Vote)
		articleRoutes.POST("star", controller.StarOrUnStar)
		articleRoutes.POST("delete", controller.DeleteArticle)
	}
	topicRoutes := router.Group("/topic")
	{
		topicRoutes.GET("gets", controller.GetTopics)
		topicRoutes.GET("get", controller.GetTopic)
	}
	commentRoutes := router.Group("/comment")
	{
		commentRoutes.POST("create", controller.CreateComment)
		commentRoutes.GET("article", controller.GetCommentsOfArticle)
		commentRoutes.GET("reply", controller.GetCommentsOfComment)
		commentRoutes.GET("get", controller.GetCommentInfo)
		commentRoutes.POST("like", controller.LikeOrUnlikeComment)
		commentRoutes.POST("delete", controller.DeleteComment)
	}

	messageRoutes := router.Group("/message")
	{
		messageRoutes.GET("unread", controller.CountUnread)
		messageRoutes.GET("gets", controller.GetMessages)
	}

	//路由在指定端口Run起来（异步）
	go func() {
		//err := router.Run(":" + config.PORT)
		//println(err)
		err := router.RunTLS(":39999", "middleware/hita.pem", "middleware/hita.key")
		if err == nil {
			fmt.Println("TLS Running...")
		} else {
			fmt.Printf("Run TLS error!:%v\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	logger.Println("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	srv := &http.Server{
		Addr:    config.PORT,
		Handler: router,
	}
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalln("Server Shutdown: ", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logger.Println("timeout of 5 seconds.")
	}

	logger.Println("Server exiting")

}
