package internal

import (
	"fmt"
	"log"
	"net"
	"time"

	"toko_buku_online/internal/config"
	"toko_buku_online/internal/entity"
	"toko_buku_online/internal/handler"
	"toko_buku_online/internal/logger"
	"toko_buku_online/internal/middleware"
	"toko_buku_online/internal/repository"
	"toko_buku_online/internal/service"
	"toko_buku_online/internal/usecase"
	"toko_buku_online/toko_buku_online/api/gen/go/toko/v1"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	log        logger.Logger
	authUc     usecase.AuhtUc
	categoryUc usecase.CategoryUc
	jwtService service.JwtService
	port       string
}

func (s *Server) Run() {
	s.log.Info("starting gRPC server...", nil)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	cfg := config.NewConfig()
	middleware := middleware.NewGRPCAuthMiddleware(s.jwtService, s.log, cfg.TokenConfig)

	// Register gRPC handler
	authHandler := handler.NewAuthHandler(s.authUc, s.log)
	categoryHandler := handler.NewCategoryHandler(s.log, s.categoryUc, middleware)
	toko.RegisterAuthServiceServer(grpcServer, authHandler)
	toko.RegisterCategoryServiceServer(grpcServer, categoryHandler)

	s.log.Info(fmt.Sprintf("gRPC server running on port %s", s.port), nil)
	// ✅ Aktifkan server reflection
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}

func NewServer() *Server {
	// Load env
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	// Logger
	log := logger.NewLogger()

	// Config
	cfg := config.NewConfig()

	// Connect DB
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database instance: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLife) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Minute)

	db = db.Debug()

	// Auto Migrate
	//func enum manual dipostgre
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

	// Dependency Injection
	// repo
	authRepo := repository.NewAuthRepo(*log, db)
	categoryRepo := repository.NewCategoryRepo(*log, db)

	// service
	jwtService := service.NewJwtService(cfg.TokenConfig, *log)

	// uc
	authUc := usecase.NewAuthUc(authRepo, *log, jwtService)
	categoryUc := usecase.NewCategoryUc(*log, categoryRepo)

	return &Server{
		log:        *log,
		authUc:     authUc,
		categoryUc: categoryUc,
		jwtService: jwtService,
		port:       cfg.ServerPort,
	}
}
