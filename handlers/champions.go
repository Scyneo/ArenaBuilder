package handlers

import (
	"context"
	"log"
	"net/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ChampionHandler struct {
  DB  *pgxpool.Pool
}

func (h *ChampionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var name string = cases.Title(language.English).String(r.PathValue("name"))
    var icon []byte

    err := h.DB.QueryRow(context.Background(), 
           "SELECT image FROM champions WHERE name = $1", name).Scan(&icon)
    if err != nil {
        http.Error(w, "Champion not found", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "image/png")

    _, err = w.Write(icon)
    if err != nil {
        // Handle error
        log.Println("Error writing icon to response:", err)
    }
}
