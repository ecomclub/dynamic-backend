package main

import (
  "log"
  "net/http"
  "os"
  "fmt"
  "github.com/go-redis/redis"
)

func main() {
  // get root directory and port from command line arguments
  // https://gobyexample.com/command-line-arguments
  // ./main /var/data :3000
  root := os.Args[1]
  // TCP port
  // eg.: ':3000'
  port := os.Args[2]

  if len(os.Args) >= 4 {
    // ./main /var/data :3000 /var/log/go.log
    file := os.Args[3]
    // log to file
    f, err := os.OpenFile(file, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
    if err != nil {
      log.Fatalf("Error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)
  }

  log.Println("------")
  log.Println("Starting dynamic backend")

  // root must end with bar
  lastRootChar := root[len(root) - 1:]
  if lastRootChar != "/" {
    root += "/"
  }
  log.Println("Server root")
  log.Println(root)

  // init Redis client
  // https://github.com/go-redis/redis
  client := redis.NewClient(&redis.Options{
    Addr: "127.0.0.1:6379",
    // no password set
		Password: "",
    // use default DB
		DB: 0,
	})

	pong, err := client.Ping().Result()
  if err != nil {
    log.Println(err)
  } else {
    log.Println("Redis ping")
    log.Println(pong)
  }

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    // get channel ID from Redis
    val, err := client.Get(r.Host).Result()
    if err == nil {
      // fmt.Fprintf(w, "Key value: %q\n", val)
      file := fmt.Sprintf("%s%s%s", root, val, r.URL.Path)
      // log.Println("GET file")
      // log.Println(file)
      http.ServeFile(w, r, file)
    } else {
      w.WriteHeader(http.StatusNotFound)
      w.Write([]byte("Not Found!\n"))
    }
  })

  log.Println("Listening...")
  log.Println(port)
  log.Fatal(http.ListenAndServe(port, nil))
}
