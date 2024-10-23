package router

import (
	"MLcore-Engine/controller"
	"MLcore-Engine/middleware"

	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiRouter := router.Group("/api")
	apiRouter.Use(middleware.GlobalAPIRateLimit())
	{
		apiRouter.GET("/status", controller.GetStatus)
		apiRouter.GET("/notice", controller.GetNotice)
		apiRouter.GET("/about", controller.GetAbout)
		apiRouter.GET("/verification", middleware.CriticalRateLimit(), controller.SendEmailVerification)
		apiRouter.GET("/reset_password", middleware.CriticalRateLimit(), middleware.TurnstileCheck(), controller.SendPasswordResetEmail)
		apiRouter.POST("/user/reset", middleware.CriticalRateLimit(), controller.ResetPassword)
		apiRouter.GET("/oauth/github", middleware.CriticalRateLimit(), controller.GitHubOAuth)
		apiRouter.GET("/oauth/wechat", middleware.CriticalRateLimit(), controller.WeChatAuth)
		apiRouter.GET("/oauth/wechat/bind", middleware.CriticalRateLimit(), middleware.UserAuth(), controller.WeChatBind)
		apiRouter.GET("/oauth/email/bind", middleware.CriticalRateLimit(), middleware.UserAuth(), controller.EmailBind)
		apiRouter.POST("/oauth/token", middleware.CriticalRateLimit(), controller.GetToken)

		userRoute := apiRouter.Group("/user")
		{
			userRoute.POST("/register", middleware.CriticalRateLimit(), controller.Register)
			userRoute.POST("/login", middleware.CriticalRateLimit(), controller.Login)
			userRoute.GET("/logout", controller.Logout)

			selfRoute := userRoute.Group("/")
			{
				selfRoute.GET("/self", controller.GetSelf)
				selfRoute.PUT("/self", controller.UpdateSelf)
				selfRoute.DELETE("/self", controller.DeleteSelf)
			}

			userManageRoute := userRoute.Group("/manage")
			userManageRoute.Use(middleware.AdminAuth())
			{
				userManageRoute.GET("/", controller.GetAllUsers)
				userManageRoute.GET("/search", controller.SearchUsers)
				userManageRoute.GET("/:id", controller.GetUser)
				userManageRoute.POST("/", controller.CreateUser)
				userManageRoute.POST("/manage", controller.ManageUser)
				userManageRoute.PUT("/", controller.UpdateUser)
				userManageRoute.DELETE("/:id", controller.DeleteUser)
			}

			adminRoute := userRoute.Group("/")
			adminRoute.Use(middleware.RootAuth())
			{
				adminRoute.GET("/", controller.GetAllUsers)
				adminRoute.GET("/search", controller.SearchUsers)
				adminRoute.GET("/:id", controller.GetUser)
				adminRoute.POST("/", controller.CreateUser)
				adminRoute.POST("/manage", controller.ManageUser)
				adminRoute.PUT("/", controller.UpdateUser)
				adminRoute.DELETE("/:id", controller.DeleteUser)
			}
		}
		optionRoute := apiRouter.Group("/option")
		optionRoute.Use(middleware.RootAuth())
		{
			optionRoute.GET("/", controller.GetOptions)
			optionRoute.PUT("/", controller.UpdateOption)
		}
		fileRoute := apiRouter.Group("/file")
		fileRoute.Use(middleware.AdminAuth())
		{
			fileRoute.GET("/", controller.GetAllFiles)
			fileRoute.GET("/search", controller.SearchFiles)
			fileRoute.POST("/", middleware.UploadRateLimit(), controller.UploadFile)
			fileRoute.DELETE("/:id", controller.DeleteFile)
		}

		projectRoute := apiRouter.Group("/project")
		projectRoute.Use(middleware.UserAuth())
		{
			projectRoute.POST("/", controller.CreateProject)
			projectRoute.PUT("/:id", controller.UpdateProject)
			projectRoute.DELETE("/:id", controller.DeleteProject)
			projectRoute.GET("/:id", controller.GetProject)
			projectRoute.GET("/get-all", controller.ListProjects)
		}

		projectMembershipsRoute := apiRouter.Group("/project-memberships")
		projectMembershipsRoute.Use(middleware.UserAuth())
		{
			projectMembershipsRoute.GET("/user/:userId", controller.GetUserProjects)
			projectMembershipsRoute.POST("/", controller.AddUserToProject)
			projectMembershipsRoute.DELETE("/:projectId/:userId", controller.RemoveUserFromProject)
			projectMembershipsRoute.PUT("/", controller.UpdateUserProjectRole)
			projectMembershipsRoute.GET("/project/:projectId", controller.GetProjectMembers)
		}

		notebookRoute := apiRouter.Group("/notebook")
		notebookRoute.Use(middleware.UserAuth())
		{
			notebookRoute.POST("/", controller.CreateNotebook)
			notebookRoute.PUT("/:id", controller.UpdateNotebook)
			notebookRoute.DELETE("/:id", controller.DeleteNotebook)
			notebookRoute.GET("/:id", controller.GetNotebook)
			notebookRoute.GET("/get-all", controller.ListNotebooks)
			notebookRoute.POST("/reset/:id", controller.ResetNotebook)
		}
	}
}
