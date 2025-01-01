package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Saparta/wishlist/wishlist/services/user-service/endpoints"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func setUpDb(channel chan *pgxpool.Pool) {
	godotenv.Load()
	var dsn string = os.Getenv("DSN")

	// Create a connection pool
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	channel <- dbPool
}

func main() {
	var dbChannel chan *pgxpool.Pool = make(chan *pgxpool.Pool)
	go setUpDb(dbChannel)

	var r *gin.Engine = gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	var dbPool *pgxpool.Pool = <-dbChannel
	defer dbPool.Close()

	r.GET("/users", func(ctx *gin.Context) { endpoints.GetUsers(ctx, dbPool) })

	r.Run()
}
