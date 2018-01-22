package main

import (
  "log"
  "net/http"
  "os"
  "fmt"
)

func main() {
  // get root directory from first command line argument
  // https://gobyexample.com/command-line-arguments
  root := os.Args[1]
  // TCP port
  port := os.Args[2]
  // listen localhost only
  bind := fmt.Sprintf("%s%s", "127.0.0.1:", port)

  fs := http.FileServer(http.Dir(root))
  http.Handle("/", fs)

  http.ListenAndServe(bind, nil)
  log.Println("Listening...")
}
