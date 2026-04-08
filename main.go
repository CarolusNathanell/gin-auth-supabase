package main

import (
	"context"
	"log"
	"os"
	"time"

	"gin-auth-supabase/src/auth"
	"gin-auth-supabase/src/db"
	headCountLog "gin-auth-supabase/src/head_count_log"
	"gin-auth-supabase/src/snapshots"
	"gin-auth-supabase/src/sources"
	websock "gin-auth-supabase/src/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// 1. Connect to Supabase Postgres via pgxpool
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer pool.Close()

	// 2. Initialize sqlc Queries and DI
	queries := db.New(pool)
	authService := auth.NewService(queries)
	SourcesService := sources.NewService(queries)
	// attendanceLogService := attendanceLog.NewService(queries)
	headCountLogService := headCountLog.NewService(queries)
	snapshotsService := snapshots.NewService(queries)

	authHandler := auth.NewHandler(authService)
	SourcesHandler := sources.NewHandler(SourcesService)
	// attendanceLogHandler := attendanceLog.NewHandler(attendanceLogService)
	headCountLogHandler := headCountLog.NewHandler(headCountLogService)
	snapshotsHandler := snapshots.NewHandler(snapshotsService)

	wsHub := websock.NewWSHub()
	go wsHub.Run()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/logs", websock.ReceiveLogs(wsHub))
	r.GET("/ws", websock.HandleWS(wsHub))

	r.Static("public/snapshots", "./public/snapshots/")

	Api := r.Group("/api")
	{
		authApi := Api.Group("/auth")
		{
			authApi.POST("/register", authHandler.HandleRegister)
			authApi.POST("/login", authHandler.HandleLogin)
		}

		profileApi := Api.Group("/profile")
		profileApi.Use(auth.AuthMiddleware())
		{
			profileApi.GET("", authHandler.HandleRequest)
			profileApi.PUT("", authHandler.HandleUpdate)
		}

		SourcesApi := Api.Group("/sources")
		SourcesApi.Use(auth.AuthMiddleware())
		{
			SourcesApi.POST("", SourcesHandler.HandleAdd)
			SourcesApi.GET("", SourcesHandler.HandleRequest)
			SourcesApi.GET("/:id", SourcesHandler.HandleRequestById)
			SourcesApi.PUT("/:id", SourcesHandler.HandleUpdateById)
			SourcesApi.DELETE("/:id", SourcesHandler.HandleDeleteById)
		}

		// attendanceLogApi := Api.Group("/al")
		// attendanceLogApi.Use(auth.AuthMiddleware())
		// {
		// 	attendanceLogApi.POST("/add", attendanceLogHandler.HandleAdd)
		// 	attendanceLogApi.GET("/request", attendanceLogHandler.HandleRequest)
		// }

		headCountLogApi := Api.Group("/logs")
		{
			headCountLogApi.POST("", headCountLogHandler.HandleAdd)
			headCountLogApi.GET("/:sourceId", headCountLogHandler.HandleRequestBySource)
		}

		SnapshotsApi := Api.Group("/snapshots")
		{
			SnapshotsApi.POST("", snapshotsHandler.HandleAdd)
			SnapshotsApi.GET("/:sourceId", snapshotsHandler.HandleRequestsBySource)
			SnapshotsApi.GET("/:sourceId/:snapshotId", snapshotsHandler.HandleRequestById)
			SnapshotsApi.DELETE("/:sourceId/:snapshotId", snapshotsHandler.HandleDeleteById)
		}
	}

	r.Run(":8080")
}
