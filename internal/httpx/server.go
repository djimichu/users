package httpx

import (
	"JWTproject/internal/auth"
	"github.com/go-chi/chi"
)

type HTTPServer struct {
	Handlers *HTTPHandlers
}

func (s *HTTPServer) Start(jwtManager *auth.JWTManager) {

	router := chi.NewRouter()
	//register
	router.Post("/register", s.Handlers.CreateUserHandler)

	//login
	router.Post("/login", s.Handlers.LoginUserHandler)

	router.Group(func(r chi.Router) {

		r.Use(auth.JWTAuthMiddleware(jwtManager))

		r.Get("/user", s.Handlers.GetUserHandler)

		r.Patch("/user", s.Handlers.ChangeUsernameHandler)

		r.Delete("/user", s.Handlers.DeleteUserHandler)
	})

}
