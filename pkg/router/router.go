package router

import (
	accountHandler "github.com/Flood-project/backend-flood/internal/account_user/handler"
	loginHandler "github.com/Flood-project/backend-flood/internal/login/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Server struct {
	Router *chi.Mux
}

func CreateNewServer() *Server {
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
		Router: r,
	}
}

func (s *Server) MountAccounts(handler *accountHandler.AccountHandler) {
	s.Router.Route("/accounts", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.Fetch)
		r.Get("/{id}", handler.GetByID)
	})
}

func (s *Server) MountLogin(handler *loginHandler.LoginHandler) {
	s.Router.Route("/login", func(r chi.Router) {
		r.Post("/", handler.Login)
	})
}
