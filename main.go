package main

import (
  "log"
  "net/http"
  "os"
  "fmt"
)

func main() {
  // get root directory and port from command line arguments
  // https://gobyexample.com/command-line-arguments
  // ./main /var/data :3000
  root := os.Args[1]
  // TCP port
  // eg.: ':3000'
  port := os.Args[2]

  fs := http.FileServer(http.Dir(root))
  http.Handle("/", fs)

  http.ListenAndServe(port, nil)
  log.Println("Listening...")
}
