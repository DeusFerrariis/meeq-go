package main

import (
  log "github.com/sirupsen/logrus"
  "net/http"
  "strconv"
  "encoding/json"
)


func RegisterRoutes(mex MessageExchange) {
  log.Info("Registering handlers")
  http.HandleFunc("/newTopic", CreateTopicHandler(mex))
  http.HandleFunc("/publish", PublishHandler(mex))
  http.HandleFunc("/subscribe", SubscribeHandler(mex))
  http.HandleFunc("/consume", ConsumeHandler(mex))
}

type HandleFunc = func(w http.ResponseWriter, r *http.Request)

func CreateTopicHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    mex.NewTopic(topic)
  }
}

func PublishHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    topic := r.URL.Query().Get("topic")
    b := r.FormValue("message_body")
    msg := Message{Body: b, Headers: make(map[string]string)}


    err := mex.Publish(topic, msg)
    if err != nil {
      w.Write([]byte("Error publishing message, " + err.Error()))
    }
  }
}

func SubscribeHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    subscriberID := r.URL.Query().Get("subscriberID")
    mex.Subscribe(topic, subscriberID)
  }
}

func ConsumeHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    subscriberID := r.URL.Query().Get("subscriberID")
    amount := r.URL.Query().Get("amount")

    a, err := strconv.Atoi(amount)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Error converting amount to integer, "))
    }

    msgs, err := mex.Consume(topic, subscriberID, a)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Error consuming messages, " + err.Error()))
      return
    }

    msgsJson, err := json.Marshal(msgs)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Error converting messages to JSON, " + err.Error()))
    }

    w.Write(msgsJson)

    return
  }
}
