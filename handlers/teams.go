package handlers

import (
  "math/rand"
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
  Players [][2]Player `json:"players"`
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
  shuffle(players)

  var teams []Player
  for _, player := range players {
    player = strings.TrimSpace(player)
    teams = append(teams, Player{Name: player})
  }

  rows, err := h.DB.Query(context.Background(), 
         "SELECT name, image FROM champions ORDER BY random() LIMIT $1", len(teams))
  if err != nil {
      http.Error(w, "Champion not found", http.StatusNotFound)
      return
  }
  defer rows.Close()

  var i int
  for rows.Next() {
    err := rows.Scan(&teams[i].Champion, &teams[i].Icon)
    if err != nil {
        http.Error(w, "Error assigning champions", http.StatusNotFound)
        fmt.Println(err)
        return
    }
    i++
  }
  var finalTeam TeamOut
  for i := 0; i < len(teams); i += 2 {
    if (i+1) == len(teams) {
      finalTeam.Players = append(finalTeam.Players, [2]Player{teams[i], {Name: "xd", Champion: "xd", Icon: []byte{}}})
      break
    }
    finalTeam.Players = append(finalTeam.Players, [2]Player{teams[i], teams[i+1]})
  }

  json.NewEncoder(w).Encode(finalTeam)
  w.Header().Set("Content-Type", "application/json")
}

func shuffle(slice []string) {
  rand.Shuffle(len(slice), func(i, j int) {
      slice[i], slice[j] = slice[j], slice[i]
  })
}
