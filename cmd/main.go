package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flood-project/backend-flood/config"
	"github.com/Flood-project/backend-flood/internal/account_user/handler"
	"github.com/Flood-project/backend-flood/internal/account_user/repository"
	"github.com/Flood-project/backend-flood/internal/account_user/usecase"
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

	accountRepository := repository.NewAccountRepository(db)
	accountUsecase := usecase.AccountUseCase(accountRepository)	
	accountHandler := handler.NewAccountHandler(accountUsecase)

	server := router.CreateNewServer()
	server.MountAccounts(accountHandler)

	log.Println("running on :8080")
	http.ListenAndServe(":8080", server.Router)
}