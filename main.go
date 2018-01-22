package main

import (
  "log"
  "net/http"
  "os"
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
    log.Println(pong)
  }

  fs := http.FileServer(http.Dir(root))
  http.Handle("/", fs)

  log.Println("Listening...")
  http.ListenAndServe(port, nil)
}
