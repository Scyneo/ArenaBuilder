package main

import (
  "fmt"
  "log"
  "net/http"
  "learning/handlers"
  "learning/data"
)

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
  dbPool, err := data.NewDb()
  if err != nil {
    log.Fatal("Cannot connect to database", err)
  }
	defer dbPool.Close()

  championHandler := &handlers.ChampionHandler{DB: dbPool}
  teamsHandler := &handlers.TeamsHandler{DB: dbPool}

  router := http.NewServeMux()
  router.Handle("GET /champions/{name}", championHandler)
  router.Handle("GET /teams", teamsHandler)

  wrappedRouter := enableCORS(router)

  server := http.Server{
    Addr:    ":8080",
    Handler: wrappedRouter,
  } 

  fmt.Println("Starting GO API service...")
  fmt.Println(`
   ______     ______        ______     ______   __    
  /\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
  \ \ \__ \  \ \ \_\ \     \ \  __ \  \ \  _-/ \ \ \  
   \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
    \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)

  server.ListenAndServe()
}
