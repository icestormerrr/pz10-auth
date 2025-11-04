package main

import (
	"log"
	"net/http"

	router "github.com/icestormerrr/pz10-auth/internal/delivery/http"
	"github.com/icestormerrr/pz10-auth/internal/delivery/http/handlers"
	"github.com/icestormerrr/pz10-auth/internal/repos"
	"github.com/icestormerrr/pz10-auth/internal/services"
	"github.com/icestormerrr/pz10-auth/internal/utils/config"
	"github.com/icestormerrr/pz10-auth/internal/utils/jwt"
)

func main() {
	cfg := config.Load()

	sessionRepo := repos.NewSessionRedisRepo(cfg)
	userRepo := repos.NewUserInMemoryRepo()
	jwtValidator, err := jwt.NewRS256(cfg.PrivateRsaKey, cfg.PublicRsaKey, cfg.AccessTTL)
	if err != nil {
		log.Fatal("cannot parse RSA keys: ", err)
	}

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, sessionRepo, jwtValidator)

	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	mux := router.Build(authHandler, userHandler, jwtValidator)
	log.Println("listening on", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
