package user

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

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers(r.Context())
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request)  {
	var input struct {
		Name string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "name and email are required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.CreateUser(r.Context(), input.Name, input.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user in the database
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int32 `json:"id"`
		Name string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "name, email, and id are required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.UpdateUser(r.Context(), input.ID, input.Name, input.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser deletes a user from the database
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int32 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := h.repo.DeleteUser(r.Context(), input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUser retrieves a user from the database
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID int32 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUser(r.Context(), input.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
