package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cicd/models"
	"cicd/handlers"
	"cicd/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/stretchr/testify/assert"
)

var client *mongo.Client
var testDB *mongo.Database

func setupTestDB(t *testing.T) {
	var err error
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	testDB = client.Database("taskdb")
	err = testDB.Drop(context.Background())
	if err != nil {
		t.Fatalf("Failed to drop database: %v", err)
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	repo := repository.NewTaskRepository(testDB)
	handler := handlers.NewTaskHandler(repo)

	r.POST("/tasks", handler.CreateTask)
	r.GET("/tasks/:id", handler.GetTaskByID)
	r.GET("/tasks", handler.GetTasks)
	r.PUT("/tasks/:id", handler.UpdateTask)
	r.DELETE("/tasks/:id", handler.DeleteTask)

	return r
}

func TestCreateTask(t *testing.T) {
	setupTestDB(t)
	defer client.Disconnect(context.Background())

	r := setupRouter()

	newTask := `{"title":"Test Task","description":"This is a test task"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte(newTask)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var task models.Task
	err := json.NewDecoder(w.Body).Decode(&task)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, task.ID.Hex())
	assert.Equal(t, "Test Task", task.Title)
	assert.Equal(t, "This is a test task", task.Description)

	// Verify the task is inserted in the database
	var dbTask models.Task
	err = testDB.Collection("tasks").FindOne(context.Background(), bson.M{"_id": task.ID}).Decode(&dbTask)
	if err != nil {
		t.Fatal("Task not found in database")
	}
}

func TestGetTasks(t *testing.T) {
	setupTestDB(t)
	defer client.Disconnect(context.Background())

	collection := testDB.Collection("tasks")
	collection.InsertOne(context.Background(), models.Task{Title: "Task 1", Description: "First task", CreatedAt: time.Now()})
	collection.InsertOne(context.Background(), models.Task{Title: "Task 2", Description: "Second task", CreatedAt: time.Now()})

	r := setupRouter()
	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var tasks []models.Task
	err := json.NewDecoder(w.Body).Decode(&tasks)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, tasks, 2)
}

func TestGetTaskByID(t *testing.T) {
	setupTestDB(t)
	defer client.Disconnect(context.Background())

	collection := testDB.Collection("tasks")
	res, err := collection.InsertOne(context.Background(), models.Task{Title: "Single Task", Description: "Task description", CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	taskID := res.InsertedID.(primitive.ObjectID).Hex()

	r := setupRouter()
	req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var task models.Task
	err = json.NewDecoder(w.Body).Decode(&task)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Single Task", task.Title)
}

func TestUpdateTask(t *testing.T) {
	setupTestDB(t)
	defer client.Disconnect(context.Background())

	collection := testDB.Collection("tasks")
	res, err := collection.InsertOne(context.Background(), models.Task{Title: "Task to Update", Description: "Update this task", CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	taskID := res.InsertedID.(primitive.ObjectID).Hex()
	updateData := `{"title":"Updated Task","description":"Task has been updated"}`
	req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer([]byte(updateData)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var task models.Task
	err = collection.FindOne(context.Background(), bson.M{"_id": res.InsertedID}).Decode(&task)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Updated Task", task.Title)
}

func TestDeleteTask(t *testing.T) {
	setupTestDB(t)
	defer client.Disconnect(context.Background())

	collection := testDB.Collection("tasks")
	res, err := collection.InsertOne(context.Background(), models.Task{Title: "Task to Delete", Description: "Delete this task", CreatedAt: time.Now()})
	if err != nil {
		t.Fatal(err)
	}

	taskID := res.InsertedID.(primitive.ObjectID).Hex()

	r := setupRouter()
	req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	err = collection.FindOne(context.Background(), bson.M{"_id": res.InsertedID}).Err()
	assert.True(t, mongo.ErrNoDocuments == err, "Expected no documents error, got %v", err)
}
