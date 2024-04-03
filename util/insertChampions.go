package main

import (
  "context"
  "fmt"
  // "net/http"
  "github.com/jackc/pgx/v5"
  "os"
  "io"
  "strings"
)

func main() {
  folderPath := "/mnt/c/Users/Mati/Desktop/champion/"
  result, err := os.ReadDir(folderPath)
  if err != nil {
    fmt.Println(err)
  }

  m := make(map[string][]byte)

  for _, fileName := range result {
    file, err := os.Open(folderPath + fileName.Name())
    if err != nil {
      fmt.Println(err)
    }

    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
      fmt.Println(err)
    }
    var championName string = strings.TrimSuffix(fileName.Name(), ".png")
    m[championName] = data

  }
  connectToDB(m)
}

func connectToDB(images map[string][]byte) {
  conn, err := pgx.Connect(context.Background(), "postgresql://postgres:admin@localhost:5432/lol")
  if err != nil {
    fmt.Println(err)
  }

  defer conn.Close(context.Background())

  for key, image := range images {
    _, err = conn.Exec(context.Background(), "INSERT INTO champions(name, image) VALUES ($1, $2)", key, image)
    if err != nil {
      fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
    }
  }
}
