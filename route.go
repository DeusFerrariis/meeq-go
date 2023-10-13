package main

import (
  log "github.com/sirupsen/logrus"
  "net/http"
)


func RegisterRoutes(mex MessageExchange) {
  log.Info("Registering endpoints")
  http.HandleFunc("/new-topic", CreateTopicHandler(mex))
  http.HandleFunc("/topic/publish", PublishHandler(mex))
  http.HandleFunc("/topic/subscribe", SubscribeHandler(mex))
  http.HandleFunc("/topic/consume", ConsumeHandler(mex))
  http.HandleFunc("/health", HealthCheckHandler())
  log.Info("Endpoints registered")
}

