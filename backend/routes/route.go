package routes

import (
	"os"
	"permit-app/controller/divisionController"
	"permit-app/controller/domainController"
	"permit-app/controller/menuController"
	"permit-app/controller/moduleController"
	"permit-app/controller/notificationController"
	"permit-app/controller/permitController"
	"permit-app/controller/permitTypeController"
	"permit-app/controller/projectController"
	"permit-app/controller/referenceCategoryController"
	"permit-app/controller/referenceController"
	"permit-app/controller/roleController"
	"permit-app/controller/taskController"
	"permit-app/controller/taskRequestController"
	"permit-app/controller/userController"
	"permit-app/middleware"
	"permit-app/repo/divisionRepository"
	"permit-app/repo/domainRepository"
	"permit-app/repo/menuRepository"
	"permit-app/repo/moduleRepository"
	"permit-app/repo/notificationRepository"
	"permit-app/repo/permitRepository"
	"permit-app/repo/permitTypeRepository"
	"permit-app/repo/projectRepository"
	"permit-app/repo/referenceCategoryRepository"
	"permit-app/repo/referenceRepository"
	"permit-app/repo/roleRepository"
	"permit-app/repo/taskRepository"
	"permit-app/repo/userRepository"
	"permit-app/service/divisionService"
	"permit-app/service/domainService"
	"permit-app/service/menuService"
	"permit-app/service/moduleService"
	"permit-app/service/notificationService"
	"permit-app/service/permitService"
	"permit-app/service/permitTypeService"
	"permit-app/service/projectService"
	"permit-app/service/referenceCategoryService"
	"permit-app/service/referenceService"
	"permit-app/service/roleService"
	"permit-app/service/taskService"
	"permit-app/service/userService"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func corsConfig() cors.Config {
	allow := os.Getenv("CORS_ALLOW_ORIGINS")
	origins := []string{"http://192.168.75.232:5006", "http://192.168.75.233:5006", "https://192.168.75.233:5006", "http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:3131", "http://127.0.0.1:3131"}
	if allow != "" {
		parts := strings.Split(allow, ",")
		origins = make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				origins = append(origins, p)
			}
		}
	}

	return cors.Config{
		AllowOrigins: origins,
		// AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

func NewRoute(db *gorm.DB) *gin.Engine {
	validate := validator.New()

	// Repositories
	domainRepo := domainRepository.NewDomainRepository(db)
	divisionRepo := divisionRepository.NewDivisionRepository(db)
	permitTypeRepo := permitTypeRepository.NewPermitTypeRepository(db)
	permitRepo := permitRepository.NewPermitRepository(db)
	roleRepo := roleRepository.NewRoleRepository(db)
	userRepo := userRepository.NewUserRepository(db)
	menuRepo := menuRepository.NewMenuRepository(db)
	notificationRepo := notificationRepository.NewNotificationRepository(db)
	moduleRepo := moduleRepository.NewModuleRepository(db)
	referenceCategoryRepo := referenceCategoryRepository.NewReferenceCategoryRepository(db)
	referenceRepo := referenceRepository.NewReferenceRepository(db)
	projectRepo := projectRepository.NewProjectRepository(db)
	taskRepo := taskRepository.NewTaskRepository(db)

	// Services
	domainSvc := domainService.NewDomainService(domainRepo)
	divisionSvc := divisionService.NewDivisionService(divisionRepo)
	permitTypeSvc := permitTypeService.NewPermitTypeService(permitTypeRepo)
	permitSvc := permitService.NewPermitService(permitRepo)
	roleSvc := roleService.NewRoleService(roleRepo)
	userSvc := userService.NewUserService(userRepo)
	menuSvc := menuService.NewMenuService(menuRepo)
	notificationSvc := notificationService.NewNotificationService(notificationRepo, permitRepo, userRepo)
	moduleSvc := moduleService.NewModuleService(moduleRepo)
	referenceCategorySvc := referenceCategoryService.NewReferenceCategoryService(referenceCategoryRepo, moduleRepo)
	referenceSvc := referenceService.NewReferenceService(referenceRepo, referenceCategoryRepo)
	taskSvc := taskService.NewTaskService(taskRepo)
	projectSvc := projectService.NewProjectService(projectRepo, referenceRepo)

	// Controllers
	domainCtrl := domainController.NewDomainController(domainSvc)
	divisionCtrl := divisionController.NewDivisionController(divisionSvc)
	permitTypeCtrl := permitTypeController.NewPermitTypeController(permitTypeSvc)
	permitCtrl := permitController.NewPermitController(permitSvc)
	roleCtrl := roleController.NewRoleController(roleSvc)
	userCtrl := userController.NewUserController(userSvc)
	menuCtrl := menuController.NewMenuController(menuSvc)
	notificationCtrl := notificationController.NewNotificationController(notificationSvc, validate)
	moduleCtrl := moduleController.NewModuleController(moduleSvc)
	referenceCategoryCtrl := referenceCategoryController.NewReferenceCategoryController(referenceCategorySvc)
	referenceCtrl := referenceController.NewReferenceController(referenceSvc)
	taskCtrl := taskController.NewTaskController(taskSvc)
	taskRequestCtrl := taskRequestController.NewTaskRequestController(taskSvc)
	projectCtrl := projectController.NewProjectController(projectSvc)

	app := gin.Default()

	// CORS
	app.Use(cors.New(corsConfig()))
	app.OPTIONS("/*any", func(c *gin.Context) { c.Status(204) })

	// Serve static files
	app.Static("/file", "./file")

	/* API Routes */

	// Public routes (no authentication required)
	auth := app.Group("/auth")
	{
		auth.POST("/login", userCtrl.Login)
	}

	// Protected routes (authentication required)
	protected := app.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth endpoints (protected)
		authProtected := protected.Group("/auth")
		{
			authProtected.GET("/profile", userCtrl.GetProfile)
			authProtected.PUT("/profile", userCtrl.UpdateProfile)
			authProtected.POST("/switch-domain", userCtrl.SwitchDomain)
		}

		// Menu endpoints
		menu := protected.Group("/menus")
		{
			menu.POST("", menuCtrl.CreateMenu)
			menu.GET("", menuCtrl.GetAllMenus)
			menu.GET("/user", menuCtrl.GetUserMenus)
			menu.GET("/:id", menuCtrl.GetMenuByID)
			menu.PUT("/:id", menuCtrl.UpdateMenu)
			menu.DELETE("/:id", menuCtrl.DeleteMenu)
			menu.POST("/:id/roles", menuCtrl.AssignRolesToMenu)
		}

		// Domain endpoints
		domain := protected.Group("/domains")
		{
			domain.POST("", domainCtrl.Create)
			domain.GET("", domainCtrl.GetAll)
			domain.GET("/:id", domainCtrl.GetByID)
			domain.PUT("/:id", domainCtrl.Update)
			domain.DELETE("/:id", domainCtrl.Delete)
		}

		// Division endpoints
		division := protected.Group("/divisions")
		{
			division.POST("", divisionCtrl.Create)
			division.GET("", divisionCtrl.GetAll)
			division.GET("/:id", divisionCtrl.GetByID)
			division.PUT("/:id", divisionCtrl.Update)
			division.DELETE("/:id", divisionCtrl.Delete)
		}

		// Permit Type endpoints
		permitType := protected.Group("/permit-types")
		{
			permitType.POST("", permitTypeCtrl.Create)
			permitType.GET("", permitTypeCtrl.GetAll)
			permitType.GET("/:id", permitTypeCtrl.GetByID)
			permitType.PUT("/:id", permitTypeCtrl.Update)
			permitType.DELETE("/:id", permitTypeCtrl.Delete)
		}

		// Permit endpoints
		permit := protected.Group("/permits")
		{
			permit.POST("", permitCtrl.Create)
			permit.GET("", permitCtrl.GetAll)
			permit.GET("/search", permitCtrl.Search)
			permit.GET("/:id", permitCtrl.GetByID)
			permit.PUT("/:id", permitCtrl.Update)
			permit.DELETE("/:id", permitCtrl.Delete)
			permit.POST("/:id/upload", permitCtrl.UploadDocument)
			permit.GET("/:id/download", permitCtrl.DownloadDocument)
			permit.GET("/:id/preview", permitCtrl.PreviewDocument)
		}

		// Role endpoints
		role := protected.Group("/roles")
		{
			role.POST("", roleCtrl.Create)
			role.GET("", roleCtrl.GetAll)
			role.GET("/:id", roleCtrl.GetByID)
			role.PUT("/:id", roleCtrl.Update)
			role.DELETE("/:id", roleCtrl.Delete)
		}

		// User endpoints
		user := protected.Group("/users")
		{
			user.POST("", userCtrl.Register)
			user.GET("", userCtrl.GetAll)
			user.GET("/:id", userCtrl.GetByID)
			user.PUT("/:id", userCtrl.Update)
			user.DELETE("/:id", userCtrl.Delete)
			user.POST("/:id/change-password", userCtrl.ChangePassword)
			user.POST("/:id/domain-roles", userCtrl.AddDomainRole)
			user.DELETE("/:id/domain-roles/:domain_id/:role_id", userCtrl.RemoveDomainRole)
			user.PUT("/:id/domain-roles/:domain_id/:role_id/set-default", userCtrl.SetDefaultDomainRole)
		}

		// Notification endpoints

		// Module endpoints
		module := protected.Group("/modules")
		{
			module.POST("", moduleCtrl.CreateModule)
			module.GET("", moduleCtrl.GetModules)
			module.GET("/:id", moduleCtrl.GetModuleByID)
			module.PUT("/:id", moduleCtrl.UpdateModule)
			module.DELETE("/:id", moduleCtrl.DeleteModule)
			module.GET("/:id/categories", referenceCategoryCtrl.GetCategoriesByModuleID)
			module.GET("/:id/references", referenceCtrl.GetReferencesByModuleID)
		}

		// Reference Category endpoints
		referenceCategory := protected.Group("/reference-categories")
		{
			referenceCategory.POST("", referenceCategoryCtrl.CreateReferenceCategory)
			referenceCategory.GET("", referenceCategoryCtrl.GetReferenceCategories)
			referenceCategory.GET("/:id", referenceCategoryCtrl.GetReferenceCategoryByID)
			referenceCategory.PUT("/:id", referenceCategoryCtrl.UpdateReferenceCategory)
			referenceCategory.DELETE("/:id", referenceCategoryCtrl.DeleteReferenceCategory)
			referenceCategory.GET("/:id/references", referenceCtrl.GetReferencesByCategoryID)
		}

		// Reference endpoints
		reference := protected.Group("/references")
		{
			reference.POST("", referenceCtrl.CreateReference)
			reference.GET("", referenceCtrl.GetReferences)
			reference.GET("/:id", referenceCtrl.GetReferenceByID)
			reference.PUT("/:id", referenceCtrl.UpdateReference)
			reference.DELETE("/:id", referenceCtrl.DeleteReference)
		}

		// Project endpoints
		project := protected.Group("/projects")
		{
			project.POST("", projectCtrl.CreateProject)
			project.GET("", projectCtrl.GetProjects)
			project.GET("/:id", projectCtrl.GetProjectByID)
			project.PUT("/:id", projectCtrl.UpdateProject)
			project.DELETE("/:id", projectCtrl.DeleteProject)
			project.POST("/:id/change-status", projectCtrl.ChangeProjectStatus)

			project.GET("/:id/users", projectCtrl.GetUsersByProjectID)
		}

		// Domain-specific project endpoints
		protected.GET("/domains/:id/projects", projectCtrl.GetProjectsByDomainID)

		// User-specific project endpoints
		protected.GET("/users/:id/projects", projectCtrl.GetProjectsByUserID)

		notification := protected.Group("/notifications")
		{
			notification.GET("", notificationCtrl.GetNotifications)
			notification.GET("/unread", notificationCtrl.GetUnreadNotifications)
			notification.GET("/unread/count", notificationCtrl.GetUnreadCount)
			notification.POST("/read", notificationCtrl.MarkAsRead)
			notification.POST("/read/all", notificationCtrl.MarkAllAsRead)

			// Task endpoints
			notification.DELETE("/:id", notificationCtrl.DeleteNotification)
		}

		tasks := protected.Group("/tasks")
		{
			tasks.POST("", taskCtrl.Create)
			tasks.GET("", taskCtrl.GetAll)
			tasks.GET("/:id", taskCtrl.GetByID)
			tasks.GET("/code/:code", taskCtrl.GetByCode)
			tasks.PUT("/:id", taskCtrl.Update)
			tasks.DELETE("/:id", taskCtrl.Delete)
			tasks.POST("/:id/change-status", taskCtrl.ChangeStatus)
			tasks.POST("/:id/change-type", taskCtrl.ChangeType)
			tasks.POST("/:id/in-review", taskCtrl.InReview)
			tasks.POST("/:id/set-reason", taskCtrl.SetReason)
			tasks.POST("/:id/set-revision", taskCtrl.SetRevision)

			// Task approval endpoints
			tasks.POST("/:id/approvals/:approval_id/approve", taskRequestCtrl.ApproveTask)
			tasks.POST("/:id/approvals/:approval_id/reject", taskRequestCtrl.RejectTask)
		}

		// Task request endpoints (approval workflow)
		taskRequests := protected.Group("/task-requests")
		{
			taskRequests.GET("", taskRequestCtrl.GetAll)
		}
	}

	return app
}
