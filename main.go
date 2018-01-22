package main

import (
  "log"
  "net/http"
  "os"
)

func main() {
  // get root directory from first command line argument
  // https://gobyexample.com/command-line-arguments
  root := os.Args[1]

  fs := http.FileServer(http.Dir(root))
  http.Handle("/", fs)

  log.Println("Listening...")
  http.ListenAndServe(":3000", nil)
}
