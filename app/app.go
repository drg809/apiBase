package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	database "github.com/nikola43/fibergormapitemplate/database"
	middlewares "github.com/nikola43/fibergormapitemplate/middleware"
	websocketsManager "github.com/nikola43/fibergormapitemplate/models"
	"github.com/nikola43/fibergormapitemplate/routes"
	"github.com/nikola43/fibergormapitemplate/utils"
	"github.com/nikola43/fibergormapitemplate/wasabis3instance"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Import Fiber Swagger
	swagger "github.com/arsmn/fiber-swagger/v2"
)

var httpServer *fiber.App

type App struct {
}

func (a *App) Initialize(port string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	PROD := os.Getenv("PROD")

	MYSQL_USER := os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE")

	S3_ACCESS_KEY := os.Getenv("S3_ACCESS_KEY")
	S3_SECRET_KEY := os.Getenv("S3_SECRET_KEY")
	S3_ENDPOINT := os.Getenv("S3_ENDPOINT")
	S3_BUCKET_NAME := os.Getenv("S3_BUCKET_NAME")
	S3_BUCKET_REGION := os.Getenv("S3_BUCKET_REGION")

	X_API_KEY := os.Getenv("X_API_KEY")
	FROM_EMAIL := os.Getenv("FROM_EMAIL")
	FROM_EMAIL_PASSWORD := os.Getenv("FROM_EMAIL_PASSWORD")

	if PROD == "0" {
		MYSQL_USER = os.Getenv("MYSQL_USER_DEV")
		MYSQL_PASSWORD = os.Getenv("MYSQL_PASSWORD_DEV")
		MYSQL_DATABASE = os.Getenv("MYSQL_DATABASE_DEV")

		S3_ACCESS_KEY = os.Getenv("S3_ACCESS_KEY_DEV")
		S3_SECRET_KEY = os.Getenv("S3_SECRET_KEY_DEV")
		S3_ENDPOINT = os.Getenv("S3_ENDPOINT_DEV")
		S3_BUCKET_NAME = os.Getenv("S3_BUCKET_NAME_DEV")
		S3_BUCKET_REGION = os.Getenv("S3_BUCKET_REGION_DEV")

		X_API_KEY = os.Getenv("X_API_KEY_DEV")
		FROM_EMAIL = os.Getenv("FROM_EMAIL_DEV")
		FROM_EMAIL_PASSWORD = os.Getenv("FROM_EMAIL_PASSWORD_DEV")
	}

	fmt.Println(S3_ACCESS_KEY)
	fmt.Println(S3_SECRET_KEY)
	fmt.Println(S3_ENDPOINT)
	fmt.Println(S3_BUCKET_NAME)
	fmt.Println(S3_BUCKET_REGION)
	fmt.Println(MYSQL_USER)
	fmt.Println(MYSQL_PASSWORD)
	fmt.Println(MYSQL_DATABASE)
	fmt.Println(X_API_KEY)
	fmt.Println(FROM_EMAIL)
	fmt.Println(FROM_EMAIL_PASSWORD)

	wasabis3instance.WasabiS3Client = utils.New(
		S3_ACCESS_KEY,
		S3_SECRET_KEY,
		S3_ENDPOINT,
		S3_BUCKET_REGION)

	InitializeDatabase(
		MYSQL_USER,
		MYSQL_PASSWORD,
		MYSQL_DATABASE)

	//database.Migrate()
	//fakedatabase.CreateFakeData()

	InitializeHttpServer(port)
}

func HandleRoutes(api fiber.Router) {
	//app.Use(middleware.Logger())

	routes.UserRoutes(api)
	routes.SignUpRoutes(api)
	routes.PresaleRoutes(api)
}

func InitializeHttpServer(port string) {
	httpServer = fiber.New(fiber.Config{
		BodyLimit: 2000 * 1024 * 1024, // this is the default limit of 4MB
	})
	httpServer.Use(middlewares.XApiKeyMiddleware)
	httpServer.Use(cors.New(cors.Config{}))

	ws := httpServer.Group("/ws")

	// Setup the middleware to retrieve the data sent in first GET request
	ws.Use(middlewares.WebSocketUpgradeMiddleware)

	// Pull out in another function
	// all the ikisocket callbacks and listeners
	InitializeSocketListeners()

	ws.Get("/:id", ikisocket.New(func(kws *ikisocket.Websocket) {
		websocketsManager.SocketInstance = kws

		// Retrieve the user id from endpoint
		userID := kws.Params("id")

		// Add the connection to the list of the connected users
		// The UUID is generated randomly and is the key that allow
		// ikisocket to manage Emit/EmitTo/Broadcast
		websocketsManager.SocketUsers[userID] = kws.UUID

		// Every websocket connection has an optional session key => value storage
		kws.SetAttribute("user_id", userID)

		//Broadcast to all the connected users the newcomer
		// kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userID, kws.UUID)), true)
		//Write welcome message
		kws.Emit([]byte(fmt.Sprintf("Socket connected")))
	}))

	api := httpServer.Group("/api") // /api
	v1 := api.Group("/v1")          // /api/v1

	v1.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
			"message": "ON",
		})
	})

	v1.Get("/swagger/*", swagger.Handler) // default

	/*
		v1.Get("/swagger/*", swagger.New(swagger.Config{ // custom
			URL: "http://127.0.0.1:3001/doc.json",
			DeepLinking: false,
			// Expand ("list") or Collapse ("none") tag groups by default
			DocExpansion: "none",
			// Prefill OAuth ClientID on Authorize popup
			OAuth: &swagger.OAuthConfig{
				AppName:  "OAuth Provider",
				ClientID: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
			},
			// Ability to change OAuth2 redirect uri location
			OAuth2RedirectUrl: "http://127.0.0.1:3001/swagger/oauth2-redirect.html",
		}))
	*/

	HandleRoutes(v1)

	err := httpServer.Listen(port)
	if err != nil {
		log.Fatal(err)
	}
}

func InitializeDatabase(user, password, databaseName string) {
	connectionString := fmt.Sprintf(
		"%s:%s@/%s?parseTime=true",
		user,
		password,
		databaseName,
	)

	DB, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	database.GormDB, err = gorm.Open(mysql.New(mysql.Config{Conn: DB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}
}

func InitializeSocketListeners() {

	// Multiple event handling supported
	ikisocket.On(ikisocket.EventConnect, func(ep *ikisocket.EventPayload) {
		fmt.Println(fmt.Sprintf("Connection socket event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})

	// On message event
	ikisocket.On(ikisocket.EventMessage, func(ep *ikisocket.EventPayload) {
		fmt.Println(fmt.Sprintf("Message socket event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})

	// On disconnect event
	ikisocket.On(ikisocket.EventDisconnect, func(ep *ikisocket.EventPayload) {
		// Remove the user from the local users
		delete(websocketsManager.SocketUsers, ep.Kws.GetStringAttribute("user_id"))
		fmt.Println(fmt.Sprintf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	ikisocket.On(ikisocket.EventClose, func(ep *ikisocket.EventPayload) {
		// Remove the user from the local users
		delete(websocketsManager.SocketUsers, ep.Kws.GetStringAttribute("user_id"))
		fmt.Println(fmt.Sprintf("Close event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})

	// On error event
	ikisocket.On(ikisocket.EventError, func(ep *ikisocket.EventPayload) {
		fmt.Println(fmt.Sprintf("Error event - User: %s", ep.Kws.GetStringAttribute("user_id")))
	})
}
