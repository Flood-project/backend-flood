package router

import (
	accountHandler "github.com/Flood-project/backend-flood/internal/account_user/handler"
	acionamentoHandler "github.com/Flood-project/backend-flood/internal/acionameto/handler"
	baseHandler "github.com/Flood-project/backend-flood/internal/base/handler"
	buchaHandler "github.com/Flood-project/backend-flood/internal/bucha/handler"
	loginHandler "github.com/Flood-project/backend-flood/internal/login/handler"
	"github.com/Flood-project/backend-flood/internal/middleware"
	ObjectStoreHandler "github.com/Flood-project/backend-flood/internal/object_store/handler"
	productHandler "github.com/Flood-project/backend-flood/internal/product/handler"
	"github.com/Flood-project/backend-flood/internal/token"
	tokenHandler "github.com/Flood-project/backend-flood/internal/token/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	Router       *chi.Mux
	TokenManager token.TokenManager
}

func CreateNewServer(tokenKey token.TokenManager) *Server {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	return &Server{
		Router:       r,
		TokenManager: tokenKey,
	}
}

func (s *Server) MountAccounts(handler *accountHandler.AccountHandler) {
	s.Router.Route("/accounts", func(r chi.Router) {
		//middleware para todas as rotas de accounts
		r.Use(middleware.CheckAuthentication(s.TokenManager))
		r.Post("/", handler.Create)
		r.Get("/", handler.Fetch)
		r.Get("/groupid", handler.FetchWithUserGroup)
		r.Get("/{id}", handler.GetByID)
		// r.Group(func(r chi.Router) {
		// 	r.Use(middleware.CheckAuthentication(s.TokenManager))
		// 	r.Post("/exampleRoutesWithMiddleware", handler.Create)
		// 	r.Get("/", handler.Fetch)
	})
}

func (s *Server) MountLogin(handler *loginHandler.LoginHandler, tokenHandler *tokenHandler.TokenHandler) {
	s.Router.Route("/login", func(r chi.Router) {
		r.Post("/", handler.Login)
		r.Post("/refresh", tokenHandler.RefreshToken)
	})
}

func (s *Server) MountProducts(handler *productHandler.ProductHandler) {
	s.Router.Route("/products", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.Fetch)
		r.Get("/buchas/acionamentos/bases", handler.FetchWithComponents)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
		r.Get("/params", handler.WithParams)
	})
}

func (s *Server) MountBase(handler *baseHandler.BaseHandler) {
	s.Router.Route("/bases", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.Fetch)
		r.Delete("/{id}", handler.Delete)
	})
}
func (s *Server) MountAcionamentos(handler *acionamentoHandler.AcionamentoHandler) {
	s.Router.Route("/acionamentos", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.Fetch)
		r.Delete("/{id}", handler.Delete)
	})
}

func (s *Server) MountBuchas(handler *buchaHandler.BuchaHandler) {
	s.Router.Route("/buchas", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.Fetch)
		r.Delete("/{id}", handler.Delete)
		r.Get("/params", handler.GetWithParams)
	})
}

func (s *Server) MountObjectStore(handler *ObjectStoreHandler.ObjectStoreHandler) {
	s.Router.Route("/files", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.CheckAuthentication(s.TokenManager))
			r.Get("/", handler.Fetch)
			r.Post("/{product_id}", handler.Create)
			r.Get("/url/{storageKey}", handler.GetFileUrl)
		})
		r.Get("/images/{storageKey}", handler.ServeImage)
	})
}
