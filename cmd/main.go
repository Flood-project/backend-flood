package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Flood-project/backend-flood/config"
	accountHandler "github.com/Flood-project/backend-flood/internal/account_user/handler"
	accountRepository "github.com/Flood-project/backend-flood/internal/account_user/repository"
	accountUseCase "github.com/Flood-project/backend-flood/internal/account_user/usecase"
	buchaHandler "github.com/Flood-project/backend-flood/internal/bucha/handler"
	buchaRepository "github.com/Flood-project/backend-flood/internal/bucha/repository"
	buchaUseCase "github.com/Flood-project/backend-flood/internal/bucha/usecase"
	loginHandler "github.com/Flood-project/backend-flood/internal/login/handler"
	loginUseCase "github.com/Flood-project/backend-flood/internal/login/usecase"
	productHandler "github.com/Flood-project/backend-flood/internal/product/handler"
	productRepository "github.com/Flood-project/backend-flood/internal/product/repository"
	productUseCase "github.com/Flood-project/backend-flood/internal/product/usecase"
	"github.com/Flood-project/backend-flood/internal/token"
	tokenRepository "github.com/Flood-project/backend-flood/internal/token/repository"
	"github.com/Flood-project/backend-flood/pkg/router"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	
	db := config.ConnectDB()

	secretKey := os.Getenv("SECRET_KEY")
	tokenManager := token.NewJWT(secretKey)

	accountRepository := accountRepository.NewAccountRepository(db)
	accountUsecase := accountUseCase.AccountUseCase(accountRepository)	
	accountHandler := accountHandler.NewAccountHandler(accountUsecase, tokenManager)

	tokenRepository := tokenRepository.NewTokenRepository(db)

	loginUseCase := loginUseCase.NewLogin(accountRepository, tokenManager, tokenRepository)
	loginHandler := loginHandler.NewLoginHandler(loginUseCase)

	productRepository := productRepository.NewProductManager(db)
	productUseCase := productUseCase.NewProductUseCase(&productRepository)
	productHandler := productHandler.NewProductHandler(productUseCase)

	buchaRepository := buchaRepository.NewBuchaManager(db)
	buchaUseCase := buchaUseCase.NewBuchaUseCase(buchaRepository)
	buchaHandler := buchaHandler.NewBuchaHandler(buchaUseCase)

	server := router.CreateNewServer(tokenManager)
	server.MountAccounts(accountHandler)
	server.MountLogin(loginHandler)
	server.MountProducts(productHandler)
	server.MountBuchas(buchaHandler)

	log.Println("running on :8080")
	http.ListenAndServe(":8080", server.Router)
}