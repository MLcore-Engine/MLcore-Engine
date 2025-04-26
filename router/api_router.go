package router

import (
	"MLcore-Engine/controller"
	"MLcore-Engine/middleware"

	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.URLNormalizer())
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
				userManageRoute.GET("/", controller.ListUsers)                  // 获取所有用户
				userManageRoute.GET("/:id", controller.GetUser)                 // 获取单个用户
				userManageRoute.POST("/", controller.CreateUser)                // 创建用户
				userManageRoute.PUT("/:id", controller.UpdateUser)              // 更新用户
				userManageRoute.DELETE("/:id", controller.DeleteUser)           // 删除用户
				userManageRoute.PUT("/:id/status", controller.UpdateUserStatus) // 更新用户状态
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

		projectMembersRoute := apiRouter.Group("/project-members")
		projectMembersRoute.Use(middleware.UserAuth())
		{
			projectMembersRoute.GET("/user/:userId", controller.GetUserProjects)
			projectMembersRoute.POST("/", controller.AddUserToProject)
			projectMembersRoute.DELETE("/:projectId/:userId", controller.RemoveUserFromProject)
			projectMembersRoute.PUT("/", controller.UpdateUserProjectRole)
			projectMembersRoute.GET("/project/:projectId", controller.GetProjectMembers)
		}

		notebookRoute := apiRouter.Group("/notebook")
		notebookRoute.Use(middleware.UserAuth())
		{
			notebookRoute.POST("/", controller.CreateNotebook)
			// notebookRoute.PUT("/:id", controller.UpdateNotebook)
			notebookRoute.DELETE("/:id", controller.DeleteNotebook)
			// notebookRoute.GET("/:id", controller.GetNotebook)
			notebookRoute.GET("/get-all", controller.ListNotebooks)
			notebookRoute.GET("/reset/:id", controller.ResetNotebook)
		}

		pytorchJobRoute := apiRouter.Group("/pytorchtrain")
		pytorchJobRoute.Use(middleware.UserAuth())
		{
			pytorchJobRoute.POST("/", controller.CreateTrainingJob)
			pytorchJobRoute.DELETE("/:id", controller.DeleteTrainingJob)
			pytorchJobRoute.GET("/:id", controller.GetTrainingJob)
			pytorchJobRoute.GET("/get-all", controller.ListTrainingJobs)
		}

		tritonDeployRoute := apiRouter.Group("/triton")
		tritonDeployRoute.Use(middleware.UserAuth())
		{
			tritonDeployRoute.POST("/", controller.CreateTritonDeploy)
			tritonDeployRoute.DELETE("/:id", controller.DeleteTritonDeploy)
			tritonDeployRoute.PUT("/:id", controller.UpdateTritonDeploy)
			// tritonDeployRoute.GET("/:id", controller.GetTritonDeploy)
			tritonDeployRoute.GET("/get-all", controller.ListTritonDeploys)
			tritonDeployRoute.GET("/config", controller.GetTritonConfig)
		}

		// router/api_router.go 中添加
		datasetRoute := apiRouter.Group("/dataset")
		datasetRoute.Use(middleware.UserAuth())
		{
			datasetRoute.POST("/", controller.CreateDataset)
			datasetRoute.GET("/", controller.ListDatasets)
			datasetRoute.GET("/:id", controller.GetDataset)
			datasetRoute.PUT("/:id", controller.UpdateDataset)
			datasetRoute.DELETE("/:id", controller.DeleteDataset)
			datasetRoute.GET("/search", controller.SearchDatasets)
			datasetRoute.GET("/project/:projectId", controller.GetProjectDatasets)

			// 数据条目相关路由
			datasetRoute.GET("/:id/entries", controller.GetDatasetEntries)
			datasetRoute.GET("/:id/entry/:entryId", controller.GetDatasetEntry)
			datasetRoute.POST("/:id/entry", controller.CreateDatasetEntry)
			datasetRoute.PUT("/:id/entry/:entryId", controller.UpdateDatasetEntry)
			datasetRoute.DELETE("/:id/entry/:entryId", controller.DeleteDatasetEntry)

			// 导入导出相关路由
			datasetRoute.POST("/:id/import", middleware.UploadRateLimit(), controller.ImportDataset)
			datasetRoute.GET("/:id/export", controller.ExportDataset)
		}

	}
}
