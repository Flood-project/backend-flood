package main

import (
	//"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Flood-project/backend-flood/config"
	accountHandler "github.com/Flood-project/backend-flood/internal/account_user/handler"
	accountRepository "github.com/Flood-project/backend-flood/internal/account_user/repository"
	accountUseCase "github.com/Flood-project/backend-flood/internal/account_user/usecase"
	acionamentoHandler "github.com/Flood-project/backend-flood/internal/acionameto/handler"
	acionamentoRepository "github.com/Flood-project/backend-flood/internal/acionameto/repository"
	acionamentoUseCase "github.com/Flood-project/backend-flood/internal/acionameto/usecase"
	
	auditLogRepository "github.com/Flood-project/backend-flood/internal/audit_log/repository"
	auditLogHandler "github.com/Flood-project/backend-flood/internal/audit_log/handler"
	auditLogUseCase "github.com/Flood-project/backend-flood/internal/audit_log/usecase"
	"github.com/Flood-project/backend-flood/internal/base/handler"
	baseRepository "github.com/Flood-project/backend-flood/internal/base/repository"
	baseUseCase "github.com/Flood-project/backend-flood/internal/base/usecase"
	buchaHandler "github.com/Flood-project/backend-flood/internal/bucha/handler"
	buchaRepository "github.com/Flood-project/backend-flood/internal/bucha/repository"
	buchaUseCase "github.com/Flood-project/backend-flood/internal/bucha/usecase"
	loginHandler "github.com/Flood-project/backend-flood/internal/login/handler"
	loginUseCase "github.com/Flood-project/backend-flood/internal/login/usecase"
	"github.com/Flood-project/backend-flood/internal/middleware"
	objectStoreHandler "github.com/Flood-project/backend-flood/internal/object_store/handler"
	objectStoreRepository "github.com/Flood-project/backend-flood/internal/object_store/repository"
	objectStoreUseCase "github.com/Flood-project/backend-flood/internal/object_store/usecase"
	productHandler "github.com/Flood-project/backend-flood/internal/product/handler"
	productRepository "github.com/Flood-project/backend-flood/internal/product/repository"
	productUseCase "github.com/Flood-project/backend-flood/internal/product/usecase"
	"github.com/Flood-project/backend-flood/internal/token"
	tokenHandler "github.com/Flood-project/backend-flood/internal/token/handler"
	tokenRepository "github.com/Flood-project/backend-flood/internal/token/repository"
	"github.com/Flood-project/backend-flood/pkg/router"

	//"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	db := config.ConnectDB()

	minIOConn, err := config.NewMinIO()
	if err != nil {
		log.Println("erro na conexão", err)
	}

	secretKey := os.Getenv("SECRET_KEY")
	tokenManager := token.NewJWT(secretKey)

	accountRepository := accountRepository.NewAccountRepository(db)
	accountUsecase := accountUseCase.AccountUseCase(accountRepository)	
	accountHandler := accountHandler.NewAccountHandler(accountUsecase, tokenManager)

	tokenRepository := tokenRepository.NewTokenRepository(db)

	loginUseCase := loginUseCase.NewLogin(accountRepository, tokenManager, tokenRepository)
	loginHandler := loginHandler.NewLoginHandler(loginUseCase)



	tokenHandler := tokenHandler.NewTokenHandler(tokenManager, accountUsecase)
	
	buchaRepository := buchaRepository.NewBuchaManager(db)
	buchaUseCase := buchaUseCase.NewBuchaUseCase(buchaRepository)
	buchaHandler := buchaHandler.NewBuchaHandler(buchaUseCase)
	
	acionamentoRepository := acionamentoRepository.NewAcionamentoManagement(db)
	acionamentoUseCase := acionamentoUseCase.NewAcionamentoUseCase(acionamentoRepository)
	acionamentoHandler := acionamentoHandler.NewAcionamentoHandler(acionamentoUseCase)
	
	baseRepository := baseRepository.NewBaseManagement(db)
	baseUseCase := baseUseCase.NewBaseUseCase(baseRepository)
	baseHandler := handler.NewBaseHandler(baseUseCase)

	objectStoreRepository := objectStoreRepository.NewObjectStoreUseCase(db, minIOConn)
	objectStoreUseCase := objectStoreUseCase.NewObjectStoreUseCase(objectStoreRepository, *minIOConn)
	objectStoreHandler := objectStoreHandler.NewObjectStoreHandler(objectStoreUseCase)

	productRepository := productRepository.NewProductManager(db, objectStoreRepository)
	productUseCase := productUseCase.NewProductUseCase(&productRepository)
	productHandler := productHandler.NewProductHandler(productUseCase)

	auditLogRepo := auditLogRepository.NewAuditLogManagement(db)
	auditLogUseCase := auditLogUseCase.NewAuditLogUseCase(auditLogRepo)
	auditMiddleware := middleware.NewAuditMiddleware(auditLogUseCase, tokenManager)
	auditHandler := auditLogHandler.NewBaseHandler(auditLogUseCase)

	server := router.CreateNewServer(tokenManager)
	server.Router.Use(auditMiddleware.GlobalAuditLog)
	server.MountAccounts(accountHandler, auditMiddleware)
	server.MountLogin(loginHandler, &tokenHandler)
	server.MountProducts(productHandler)
	server.MountBuchas(buchaHandler)
	server.MountAcionamentos(acionamentoHandler)
	server.MountBase(baseHandler)
	server.MountObjectStore(objectStoreHandler)
	server.MountLogs(auditHandler)

	log.Println("running on :8080")
	log.Println("minIO running on :9000", minIOConn)
	http.ListenAndServe(":8080", server.Router)
}