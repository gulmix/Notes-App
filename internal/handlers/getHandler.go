package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gulmix/Notes-App/internal/repo"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	db, err := repo.New("postgres://postgres:psql@localhost:5432/notes")
	if err != nil {
		http.Error(w, "Failed to connect", http.StatusInternalServerError)
	}
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
	}
	var data string
	var expires_at time.Time
	var max_views int
	var views int
	err = db.Pool.QueryRow(context.Background(), `SELECT content, expires_at, max_views, views FROM notes WHERE id = $1`, id).Scan(&data, &expires_at, &max_views, &views)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Note not found", http.StatusNotFound)
	}
	_, err = db.Pool.Exec(context.Background(), `UPDATE notes SET views = views + 1 WHERE id = $1`, id)
	if err != nil {
		log.Fatal(err)
	}
	if (!expires_at.IsZero() && time.Now().After(expires_at)) || (max_views > 0 && views >= max_views) {
		db.Pool.Exec(context.Background(), `DELETE FROM notes WHERE id = &1`, id)
		http.Error(w, "Note expired", http.StatusGone)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(data))
}
