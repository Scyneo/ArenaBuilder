package main

import (
  "fmt"
  "log"
  "net/http"
  "learning/handlers"
  "learning/data"
)

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println("Received a request at my domain")
	w.Write([]byte("Hello, ME"))
}

func main() {
  dbPool, err := data.NewDb()
  if err != nil {
    log.Fatal("Cannot connect to database", err)
  }
	defer dbPool.Close()

  championHandler := &handlers.ChampionHandler{DB: dbPool}

  router := http.NewServeMux()
  router.HandleFunc("/", handle)
  router.Handle("/champions/{name}", championHandler)

  server := http.Server{
    Addr:    ":8080",
    Handler: router,
  } 

  fmt.Println("Starting GO API service...")
  fmt.Println(`
   ______     ______        ______     ______   __    
  /\  ___\   /\  __ \      /\  __ \   /\  == \ /\ \   
  \ \ \__ \  \ \ \/\ \     \ \  __ \  \ \  _-/ \ \ \  
   \ \_____\  \ \_____\     \ \_\ \_\  \ \_\    \ \_\ 
    \/_____/   \/_____/      \/_/\/_/   \/_/     \/_/ `)

  server.ListenAndServe()
}
