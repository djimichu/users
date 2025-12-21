package httpx

import (
	"JWTproject/internal/auth"
	"JWTproject/internal/repository"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
)

type HTTPServer struct {
	Handlers *HTTPHandlers
}

func NewHTTPServer(handlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{Handlers: handlers}
}

func (s *HTTPServer) Start(port string, jwtManager *auth.JWTManager, userRepo *repository.UserRepo) error {

	router := chi.NewRouter()
	//register
	router.Post("/register", s.Handlers.CreateUserHandler)

	//login
	router.Post("/login", s.Handlers.LoginUserHandler)

	router.Group(func(r chi.Router) {

		r.Use(JWTAuthMiddleware(jwtManager, userRepo))

		r.Get("/user", s.Handlers.GetUserHandler)

		r.Patch("/user", s.Handlers.ChangeUsernameHandler)

		r.Delete("/user", s.Handlers.DeleteUserHandler)
	})

	if err := http.ListenAndServe(port, router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
