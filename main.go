package main

import (
	"MLcore-Engine/common"
	"MLcore-Engine/middleware"
	"MLcore-Engine/model"
	"MLcore-Engine/router"
	"embed"
	"log"
	"os"
	"strconv"

	_ "MLcore-Engine/docs"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed  web/build
var buildFS embed.FS

//go:embed web/build/index.html
var indexPage []byte

// @title MLcore-Engine API
// @version 1.0
// @description This is the API documentation for MLcore-Engine.
// @host localhost:3000
// @BasePath /api
func main() {

	if err := common.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// fmt.Println(viper.GetString("server.mode"))

	common.SetupGinLog()
	common.SysLog("Gin Template " + common.Version + " started")
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Initialize SQL Database
	err := model.InitDB()
	if err != nil {
		common.FatalLog(err)
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.FatalLog(err)
		}
	}()

	// Initialize Redis
	err = common.InitRedisClient()
	if err != nil {
		common.FatalLog(err)
	}

	// Initialize options
	model.InitOptionMap()

	// Initialize HTTP server
	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())
	//server.Use(gzip.Gzip(gzip.DefaultCompression))
	server.Use(middleware.CORS())

	// Initialize session store
	if common.RedisEnabled {
		opt := common.ParseRedisOption()
		store, _ := redis.NewStore(opt.MinIdleConns, opt.Network, opt.Addr, opt.Password, []byte(common.SessionSecret))
		server.Use(sessions.Sessions("session", store))
	} else {
		store := cookie.NewStore([]byte(common.SessionSecret))
		server.Use(sessions.Sessions("session", store))
	}
	server.RedirectTrailingSlash = false
	server.RedirectFixedPath = false
	router.SetRouter(server, buildFS, indexPage)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}
	err = server.Run(":" + port)
	if err != nil {
		log.Println(err)
	}
}
