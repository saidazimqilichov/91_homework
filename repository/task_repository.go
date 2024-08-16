package repository

import (
	"cicd/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task) error {
	task.ID = primitive.NewObjectID()
	task.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, task)
	return err
}

func (r *TaskRepository) GetTaskByID(ctx context.Context, id string) (*models.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var task models.Task
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	return &task, err
}

func (r *TaskRepository) GetTasks(ctx context.Context) ([]*models.Task, error) {
	var tasks []*models.Task
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, id string, update map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": update})
	return err
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
