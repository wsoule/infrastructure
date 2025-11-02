package product

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.ListProducts(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Stock       int32  `json:"stock"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "name, description, price, and stock are required", http.StatusBadRequest)
		return
	}

	product, err := h.repo.CreateProduct(r.Context(), input.Name, input.Description, input.Price, input.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct updates a product in the database
func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID          int32  `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Stock       int32  `json:"stock"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "id, name, description, price, and stock are required", http.StatusBadRequest)
		return
	}

	product, err := h.repo.UpdateProduct(r.Context(), input.ID, input.Name, input.Description, input.Price, input.Stock)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// DeleteProduct deletes a product from the database
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int32 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := h.repo.DeleteProduct(r.Context(), input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetProduct retrieves a product from the database
func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int32 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// Alternative: Parse from URL path parameter
	// id := r.PathValue("id")
	// idInt, err := strconv.Atoi(id)

	product, err := h.repo.GetProduct(r.Context(), input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
