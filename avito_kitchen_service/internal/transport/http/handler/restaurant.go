package handler

import (
	"avito-kitchen-service/internal/service"
	"avito-kitchen-service/internal/transport/http/dto"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RestaurantHandler struct {
	restaurantService *service.RestaurantService
}

func NewRestaurantHandler(service *service.RestaurantService) *RestaurantHandler {
	return &RestaurantHandler{restaurantService: service}
}

func (h *RestaurantHandler) GetAllRestaurants(w http.ResponseWriter, r *http.Request) {

	restaurants, err := h.restaurantService.GetAll(r.Context())

	if err != nil {
		WriteError(w, "INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	res := dto.RestaurantsDto{Restaurants: dto.ToDTORestaurants(restaurants)}

	WriteJSON(w, http.StatusOK, res)

}

func (h *RestaurantHandler) GetProducts(w http.ResponseWriter, r *http.Request) {

	restaurantIDStr := chi.URLParam(r, "id")

	restaurantID, err := strconv.ParseInt(restaurantIDStr, 10, 64)

	if err != nil {
		WriteError(w, "BAD_REQUEST", "invalid restaurant id", http.StatusBadRequest)
		return
	}

	products, err := h.restaurantService.GetProducts(r.Context(), restaurantID)

	if err != nil {
		WriteError(w, "NOT_FOUND", "Restaurant not found", http.StatusNotFound)
		return
	}

	res := dto.ProductsResponse{Products: products}

	WriteJSON(w, http.StatusOK, res)

}
