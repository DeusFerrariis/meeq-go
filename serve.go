package main

import (
  "net/http"
  //"log"
  "github.com/redis/go-redis/v9"
  "context"
  //. "github.com/WAY29/icecream-go/icecream"
  log "github.com/sirupsen/logrus"
)

func main() {
  log.Info("Starting server")
  log.Info("Creating redis client on localhost:6379")

  rc := RedisClient{
    client: redis.NewClient(&redis.Options{
      Addr: "localhost:6379",
      Password: "",
      DB: 0,
    }),
    ctx: context.Background(),
  }

  RegisterRoutes(&rc)

  log.Info("Listening on port 8080")
  http.ListenAndServe(":8080", nil)
}
