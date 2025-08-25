package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Flood-project/backend-flood/config"
	accountHandler "github.com/Flood-project/backend-flood/internal/account_user/handler"
	"github.com/Flood-project/backend-flood/internal/account_user/repository"
	accountUseCase "github.com/Flood-project/backend-flood/internal/account_user/usecase"
	loginHandler "github.com/Flood-project/backend-flood/internal/login/handler"
	loginUseCase "github.com/Flood-project/backend-flood/internal/login/usecase"
	"github.com/Flood-project/backend-flood/internal/token"
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

	accountRepository := repository.NewAccountRepository(db)
	accountUsecase := accountUseCase.AccountUseCase(accountRepository)	
	accountHandler := accountHandler.NewAccountHandler(accountUsecase, tokenManager)

	loginUseCase := loginUseCase.NewLogin(accountRepository, tokenManager)
	loginHandler := loginHandler.NewLoginHandler(loginUseCase)

	server := router.CreateNewServer()
	server.MountAccounts(accountHandler)
	server.MountLogin(loginHandler)

	log.Println("running on :8080")
	http.ListenAndServe(":8080", server.Router)
}