package internal

import (
	"fmt"
	"time"
	"toko_buku_online/internal/config"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

type Server struct {
	log logger.Logger
}

func (s *Server) initRoute() {
}

func (s *Server) Run() {
	s.log.Info("server running", nil)
	s.initRoute()
}

func NewServer() *Server {
	//logger
	log := logger.NewLogger()
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	cfg := config.NewConfig()

	//postgres
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance")
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLife) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Minute)

	// === AUTO MIGRATE ===
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Category{},
		&entity.Book{},
		&entity.Order{},
		&entity.OrderItem{},
	); err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	log.Info("success connect database", nil)

	return &Server{
		log: *log,
	}
}
