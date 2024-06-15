package product

import (
	"net/http"

	"github.com/dclouisDan/ecommerce-api/types"
	"github.com/dclouisDan/ecommerce-api/utils"
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
