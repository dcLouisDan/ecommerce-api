package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dclouisDan/ecommerce-api/config"
	"github.com/dclouisDan/ecommerce-api/service/auth"
	"github.com/dclouisDan/ecommerce-api/types"
	"github.com/dclouisDan/ecommerce-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.LoginUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
    log.Printf("Parse Error: %s", err.Error())
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

  u, err := h.store.GetUserByEmail(payload.Email)
  if err != nil {
    utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
    return
  }

  if !auth.ComparePasswords(u.Password, []byte(payload.Password)) {
    utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
    return
  }

  secret := []byte(config.Envs.JWTSecret)
  token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

  if utils.WriteJSON(w, http.StatusOK, map[string]string{"toker": token}) != nil {
    log.Fatal("Write JSON error")
  }

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get JSON payload
	var payload types.RegisterUserPayload

	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := utils.Validator.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	// check if user exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exists", payload.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if it doesn't we create a new user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if utils.WriteJSON(w, http.StatusCreated, nil) != nil {
    log.Fatal("Write JSON error")
  }
}
