package main

import (
	"context"
	"log"
	"time"

	"github.com/elgntt/avito-internship-2023/internal/api"
	"github.com/elgntt/avito-internship-2023/internal/config"
	"github.com/elgntt/avito-internship-2023/internal/pkg/db"
	"github.com/elgntt/avito-internship-2023/internal/repository"
	"github.com/elgntt/avito-internship-2023/internal/service"
)

// @title Backend-trainee-assignment-2023
// @version 1.0
// @description API Dynamic User Segmentation service

// @host localhost:8080
// @BasePath /

func main() {
	dbCfg, err := config.GetDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := db.OpenDB(ctx, dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	historyRepo := repository.NewHistoryRepo(pool)
	segmentRepo := repository.NewSegmentRepo(pool)
	userRepo := repository.NewUserRepo(pool)

	historyService := service.NewHistoryService(
		historyRepo,
		segmentRepo,
	)
	r := api.New(
		service.NewUserService(
			userRepo,
			segmentRepo,
			historyRepo,
		),
		historyService,
		service.NewSegmentService(
			segmentRepo,
			historyRepo,
			userRepo,
		),
	)

	go ClearExpiredSegmentsWorker(ctx, historyService)

	serverCfg := config.GetServerConfig()

	log.Println("Server has been successfully started on the port:" + serverCfg.HTTPPort)
	log.Fatal(r.Run(serverCfg.HTTPPort))
}

func ClearExpiredSegmentsWorker(ctx context.Context, s *service.HistoryService) {
	workerInterval := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-workerInterval.C:
			err := s.DeleteExpiredUserSegments(ctx)
			if err != nil {
				log.Println("Worker err:", err)
			}
			log.Println("Success")
		}
	}

}
