package main

import (
	"context"
	"log"
	"os"
	"time"

	apicontrollers "github.com/Yefhem/mongo/dictionary/api/api_controllers"
	"github.com/Yefhem/mongo/dictionary/cms/controllers"
	"github.com/Yefhem/mongo/dictionary/configs"
	"github.com/Yefhem/mongo/dictionary/repository"
	"github.com/Yefhem/mongo/dictionary/services"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// --------> Crate Context and Get Config
	ctx           context.Context       = context.TODO()
	configuration configs.Configuration = configs.GetConfig()
	// --------> Database
	database *mongo.Database = configs.ConnectDB(ctx, configuration.MConf)
	// --------> Create Store
	store *sessions.CookieStore = sessions.NewCookieStore([]byte("secret"))
	// --------> Collection Layer
	wordCollection *mongo.Collection = database.Collection(configuration.MConf.WordCollection)
	userCollection *mongo.Collection = database.Collection(configuration.MConf.UserCollection)

	// --------> Repository Layer
	wordRepository repository.WordRepository = repository.NewWordRepository(ctx, wordCollection)
	userRepository repository.UserRepository = repository.NewUserRepository(ctx, userCollection)

	// --------> Service Layer
	wordService     services.WordService     = services.NewWordService(wordRepository)
	userService     services.UserService     = services.MewUserService(userRepository)
	sessionService  services.SessionService  = services.NewSessionService(userRepository, store)
	alertService    services.AlertService    = services.NewAlertService(store)
	seederService   services.SeederService   = services.NewSeederService(userRepository)
	paginateService services.PaginateService = services.NewPaginateService(wordRepository)

	// --------> controller Layer Cms
	wordController controllers.WordController = controllers.NewWordController(wordService, userService, sessionService, alertService)
	authController controllers.AuthController = controllers.NewAuthController(userService, sessionService, alertService)

	// --------> controller Layer Api
	searchController apicontrollers.SearchController = apicontrollers.NewSearchController(paginateService)
)

func init() {
	seederService.UploadUser()
}

func main() {
	log.Println("running...")

	e := echo.New()
	e.Debug = true

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[*] time=${time_rfc3339_nano}\n[*] status=${status}\n[*] method=${method}\n[*] uri=${uri}\n[*] ip=${remote_ip}\n[*] Latency=${latency_human}\n[*] Byte In=${bytes_in}\n[*] user Agent=${user_agent}\n-------->\n",
	}),
		middleware.Recover(),
		middleware.CORS(),
		middleware.BodyLimit("1024K"),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5)),
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout:      2 * time.Second,
			ErrorMessage: "Timeout Error",
		}))

	e.Static("/cms/assets", "cms/assets")
	e.Static("/uploads/*filepath", "uploads")

	a := e.Group("/api/v1/word")
	// --------> API
	a.GET("/search", searchController.SearchingWord)

	g := e.Group("/admin")

	// --------> Pages
	g.GET("/dashboard", wordController.DashboardIndex)
	g.GET("/new-word", wordController.DashboardNewWord)
	g.GET("/word-edit/:id", wordController.DashboardEditWord)
	// --------> Operations
	g.POST("/add-word", wordController.AddWord)
	g.POST("/update-word/:id", wordController.UpdateWord)
	g.GET("/delete-word/:id", wordController.DeleteWord)
	// --------> Login Page
	g.GET("/login", authController.LoginIndex)
	// --------> Login Operations
	g.POST("/user-login", authController.Login)
	g.GET("/user-logout", authController.Logout)

	port := os.Getenv("PORT")

	e.Start(":" + port)
}
