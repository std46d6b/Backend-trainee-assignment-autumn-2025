package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/config"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/delivery/server"
	pullrequestservice "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/service/pull_request"
	teamservice "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/service/team"
	userservice "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/service/user"

	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/store/postgres"
	postgresrepo "github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/store/postgres/repo"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal()
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	txManager := postgres.NewTxManager(pool)
	builder := postgres.NewStatementBuilder()
	repoFactory := postgresrepo.NewRepoFactory(builder)

	teamService := teamservice.NewTeamService(txManager, pool, repoFactory)
	userService := userservice.NewUserService(txManager, pool, repoFactory)
	pullRequestService := pullrequestservice.NewPullRequestService(txManager, pool, repoFactory)

	server := server.NewServer(teamService, userService, pullRequestService)

	go func() {
		err = server.Start(fmt.Sprintf("%s:%d", cfg.WebServerConfig.Address, cfg.WebServerConfig.Port))
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.WebServerConfig.ShutdownTimeout)

	defer cancel()

	if err = server.Stop(ctx); err != nil {
		log.Printf("error stopping server: %v\n", err)
	}

	log.Println("Server stopped")
}
