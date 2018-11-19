package main

import (
  "log"
  "net/http"
  "os"
  "strings"
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
    // get store ID from Redis
    val, err := client.Get(r.Host).Result()
    if err == nil {
      // fmt.Fprintf(w, "Key value: %q\n", val)
      // {storeId}@{storeObjectId}@{channelId}[@{defaultLang}]
      s := strings.Split(val, "@")
      storeId := s[0]
      storeObjectId := s[1]
      channelId := s[2]
      var defaultLang string
      if len(s) > 3 {
        defaultLang = s[3]
      }

      slug := strings.TrimPrefix(r.URL.Path, "/")
      // replace / with $ on slug
      key := storeId + "@" + strings.Replace(slug, "/", "$", -1)

      val, err = client.Get(key).Result()
      if err == nil {
        // slug found
        // [resource]@[id]
        s = strings.Split(val, "@")
        resource := s[0]
        id := s[1]

        // write cookies to use same info client side
        http.SetCookie(w, &http.Cookie{Name: "Ecom.store_id", Value: storeId, MaxAge: 60})
        http.SetCookie(w, &http.Cookie{Name: "Ecom.store_object_id", Value: storeObjectId, MaxAge: 60})
        http.SetCookie(w, &http.Cookie{Name: "Ecom.channel_id", Value: channelId, MaxAge: 60})
        if defaultLang != "" {
          http.SetCookie(w, &http.Cookie{Name: "Ecom.default_lang", Value: defaultLang, MaxAge: 60})
        }
        http.SetCookie(w, &http.Cookie{Name: "Ecom." + r.URL.Path + ":resource", Value: resource, MaxAge: 60})
        http.SetCookie(w, &http.Cookie{Name: "Ecom." + r.URL.Path + ":_id", Value: id, MaxAge: 60})
        http.SetCookie(w, &http.Cookie{Name: "Ecom.path:_id", Value: id, MaxAge: 60})
        http.SetCookie(w, &http.Cookie{Name: "Ecom.path._id", Value: id, MaxAge: 30})
        // debug
        log.Println("&Resource ID")
        log.Println(id)

        // files from channel directory
        dir := root + channelId
        // try compinled on dist folder
        file := dir + "/.dist/_" + resource + ".html"
        if _, err := os.Stat(file); os.IsNotExist(err) {
          // dist file does not exists
          // try on channel's root directory
          file = dir + "/_" + resource + ".html"
        }

        http.ServeFile(w, r, file)
        return
      }
    }

    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("Not Found!\n"))
  })

  log.Println("Listening...")
  log.Println(port)
  log.Fatal(http.ListenAndServe(port, nil))
}
