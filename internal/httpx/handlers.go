package httpx

import (
	"JWTproject/internal/auth"
	"JWTproject/internal/httpx/response"
	"JWTproject/internal/models"
	"JWTproject/internal/service/user"
	"encoding/json"
	"errors"
	"net/http"
)

type HTTPHandlers struct {
	userService *user.UserService
}

func NewHTTPHandlers(userService *user.UserService) *HTTPHandlers {
	return &HTTPHandlers{userService: userService}
}

func (h *HTTPHandlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	var userHttp models.UserRequestDto

	if err := json.NewDecoder(r.Body).Decode(&userHttp); err != nil {
		response.WriteHttpError(w, err, http.StatusBadRequest)
		return
	}
	if err := userHttp.Validate(); err != nil {
		response.WriteHttpError(w, err, http.StatusBadRequest)
		return
	}

	userID, err := h.userService.CreateUser(userHttp)
	if err != nil {
		response.WriteHttpError(w, err, http.StatusInternalServerError)
		return
	}
	userResponse := models.NewUserRegRespDTO(userHttp.Name, userID)

	response.WriteJSON(w, userResponse, http.StatusCreated)
}
func (h *HTTPHandlers) LoginUserHandler(w http.ResponseWriter, r *http.Request) {

	var userHttp models.UserRequestDto

	if err := json.NewDecoder(r.Body).Decode(&userHttp); err != nil {
		response.WriteHttpError(w, err, http.StatusBadRequest)
		return
	}
	token, err := h.userService.Login(userHttp)
	if err != nil {
		response.WriteHttpError(w, err, http.StatusUnauthorized)
		return
	}

	resp := models.LoginTokenDTO{Token: token}

	response.WriteJSON(w, resp, http.StatusOK)
}
func (h *HTTPHandlers) GetUserHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		response.WriteHttpError(w, errors.New("user_id not found"), http.StatusUnauthorized)
		return
	}

	userDTO, err := h.userService.GetUserByID(id)
	if err != nil {
		response.WriteHttpError(w, err, http.StatusInternalServerError)
		return
	}

	response.WriteJSON(w, userDTO, http.StatusOK)
}
func (h *HTTPHandlers) ChangeUsernameHandler(w http.ResponseWriter, r *http.Request) {

	var updateUser models.UserUpdateNameDTO

	if err := json.NewDecoder(r.Body).Decode(&updateUser); err != nil {
		response.WriteHttpError(w, err, http.StatusBadRequest)
		return
	}

	if err := updateUser.Validate(); err != nil {
		response.WriteHttpError(w, err, http.StatusBadRequest)
		return
	}

	id, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		response.WriteHttpError(w, errors.New("user_id not found"), http.StatusUnauthorized)
		return
	}

	userDTO, err := h.userService.ChangeUserName(id, updateUser.Name)
	if err != nil {
		response.WriteHttpError(w, err, http.StatusInternalServerError)
		return
	}

	response.WriteJSON(w, userDTO, http.StatusOK)
}

// 204

func (h *HTTPHandlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {

	id, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		response.WriteHttpError(w, errors.New("user_id not found"), http.StatusUnauthorized)
		return
	}

	err := h.userService.DeleteUserByID(id)
	if err != nil {
		response.WriteHttpError(w, err, http.StatusInternalServerError)
		return
	}

	response.WriteJSON(w, nil, http.StatusNoContent)
}
