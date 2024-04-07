package handlers

import (
  "context"
	"fmt"
	"net/http"
	"strings"
  "encoding/base64"
  "encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
)
  
type TeamsHandler struct {
  DB  *pgxpool.Pool
}

type TeamIn struct {
  Players []string
}

type TeamOut struct {
  Players []Player `json:"players"`
}

type Player struct {
  Name     string `json:"name"`
  Champion string `json:"champion"`
  Icon     []byte `json:"icon"`
}

func (p *Player) MarshalJSON() ([]byte, error) {
  type Alias Player
  return json.Marshal(&struct {
    *Alias
    Icon    string `json:"icon"`
  }{
    Alias: (*Alias)(p),
    Icon:   base64.StdEncoding.EncodeToString(p.Icon),
  })
}


func (h *TeamsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  playersQuery := r.URL.Query().Get("players")
  players := strings.Split(playersQuery, ",")

  var team TeamOut
  for i, player := range players {
    players[i] = strings.TrimSpace(player)
    team.Players = append(team.Players, Player{Name: players[i]})
  }

  rows, err := h.DB.Query(context.Background(), 
         "SELECT name, image FROM champions ORDER BY random() LIMIT $1", len(team.Players))
  if err != nil {
      http.Error(w, "Champion not found", http.StatusNotFound)
      return
  }
  defer rows.Close()

  var i int
  for rows.Next() {
    err := rows.Scan(&team.Players[i].Champion, &team.Players[i].Icon)
    if err != nil {
        http.Error(w, "Error assigning champions", http.StatusNotFound)
        fmt.Println(err)
        return
    }
    i++
  }

  json.NewEncoder(w).Encode(team)
  w.Header().Set("Content-Type", "application/json")
  
}
