package main

import (
	"cicd/config"
	"cicd/handlers"
	"cicd/repository"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := repository.NewMngoClient(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	taskRepo := repository.NewTaskRepository(client.Database(cfg.MongoDBName))
	handler := handlers.NewTaskHandler(taskRepo)

	router := gin.Default()

	router.POST("/tasks", handler.CreateTask)
	router.GET("/tasks", handler.GetTasks)
	router.GET("/tasks/:id", handler.GetTaskByID)
	router.PUT("/tasks/:id", handler.UpdateTask)
	router.DELETE("/tasks/:id", handler.DeleteTask)

	log.Printf("Server started on port %s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
