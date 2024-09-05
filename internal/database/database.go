package database

import (
	"go-tasker/schemas"
	"go-tasker/types"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	GetLists() ([]schemas.List, error)
	CreateList(payload types.CreateListPayload) (*schemas.List, error)
	UpdateList(listID string, payload types.UpdateListPayload) (*schemas.List, error)
	DeleteList(listID string) error
	GetTasks(listID string) ([]schemas.Task, error)
	CreateTask(listID string, payload types.CreateTaskPayload) (*schemas.Task, error)
	UpdateTask(taskID string, payload types.UpdateTaskPayload) (*schemas.Task, error)
	UpdateTaskDone(taskID string, payload types.UpdateTaskDonePayload) (*schemas.Task, error)
	DeleteTask(taskID string) error
	DeleteAllTasks(listID string) error
}

type service struct {
	db *gorm.DB
}

var (
	dbUrl      = os.Getenv("DB_URL")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Create DB and connect
	db, err := gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the Schema
	err = db.AutoMigrate(&schemas.List{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&schemas.Task{})
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}
