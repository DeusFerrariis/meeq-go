package main

import (
  "net/http"
  "os"
  "context"
  "strconv"
  "github.com/redis/go-redis/v9"
  toml "github.com/pelletier/go-toml/v2"
  //. "github.com/WAY29/icecream-go/icecream"
  log "github.com/sirupsen/logrus"
)

func main() {
  log.Info("Starting server")
  log.Info("Creating redis client connection to localhost:6379")

  cfg := FetchConfig()
  rc := NewRedisClient("localhost", 6379)

  RegisterRoutes(&rc)

  p := strconv.Itoa(cfg.Server.Port)
  log.Info("Listening on port [" + p + "]")
  http.ListenAndServe(":" + p, nil)
}

type MeeqConfig struct {
  Server struct {
    Port int
  }
}

func FetchConfig() MeeqConfig {
  p := FetchConfigPath()

  c, err := os.ReadFile(p)
  if err != nil {
    log.Fatal("Error reading meeq.toml, " + err.Error())
  }

  config := MeeqConfig{}
  if err := toml.Unmarshal(c, &config); err != nil {
    log.Fatal("Error parsing meeq.toml, " + err.Error())
  }

  return config
}

func FetchConfigPath() string {
  if c := os.Getenv("MEEQ_CONFIG_PATH"); c != "" {
    return c
  }

  if c := os.Getenv("XDG_CONFIG_DIR"); c != "" {
    return c + "/meeq.toml"
  }

  c := os.Getenv("HOME")
  return c + "/.config/meeq.toml"
}

func NewRedisClient(host string, port int) RedisClient {
  p := ":" + strconv.Itoa(port)
  addr := host + p

  return RedisClient{
    client: redis.NewClient(&redis.Options{
      Addr: addr,
      Password: "",
      DB: 0,
    }),
    ctx: context.Background(),
  }
}
