package product

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dclouisDan/ecommerce-api/types"
	"github.com/dclouisDan/ecommerce-api/utils"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	store types.ProductStore
}

func NewHandler(store types.ProductStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/products", h.handleGetProduct).Methods(http.MethodGet)
	router.HandleFunc("/products", h.handleCreateProduct).Methods(http.MethodPost)
}

func (h *Handler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	ps, err := h.store.GetProducts()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if utils.WriteJSON(w, http.StatusOK, ps) != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload

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

  err := h.store.CreateProduct(types.Product{
    Name: payload.Name,
    Description: payload.Description,
    Image: payload.Image,
    Price: payload.Price,
    Quantity: payload.Quantity,
  })
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if utils.WriteJSON(w, http.StatusCreated, nil) != nil {
    log.Fatal("Write JSON error")
  }
}
