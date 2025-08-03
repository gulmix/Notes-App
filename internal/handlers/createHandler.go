package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Notes-App/internal/repo"
)

type CreateReq struct {
	Data      string `json:"data"`
	ExpiresIn int    `json:"expires_in"`
	MaxViews  int    `json:"max_views"`
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateReq
	db, err := repo.New("postgres://postgres:psql@localhost:5432/notes")
	if err != nil {
		http.Error(w, "Failed to connect", http.StatusInternalServerError)
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	expiresAt := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
	_, err = db.Pool.Exec(context.Background(), `INSERT INTO notes (id, content, expires_at, max_views) VALUES ($1, $2, $3, $4)`, &id, &req.Data, expiresAt, &req.MaxViews)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to save", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"id": id.String()})
}
