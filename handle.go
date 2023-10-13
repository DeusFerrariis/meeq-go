package main

import (
  "net/http"
  "strconv"
  "encoding/json"
  log "github.com/sirupsen/logrus"
)

type HandleFunc = func(w http.ResponseWriter, r *http.Request)

func HealthCheckHandler() func (w http.ResponseWriter, r *http.Request) {
  return func (w http.ResponseWriter, r *http.Request) {
    log.Info("Health check")
    // Check redis connection
    w.WriteHeader(http.StatusOK)
    return
  }
}

func CreateTopicHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    if err := mex.NewTopic(topic); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("Error creating topic"))
    }
  }
}

func PublishHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    topic := r.URL.Query().Get("topic")
    b := r.FormValue("message_body")
    h := r.FormValue("message_headers")

    var hJSON map[string]string

    if h != "" {
      if err := json.Unmarshal([]byte(h), &hJSON); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("Error parsing headers, " + err.Error()))
      }
    } else {
      hJSON = make(map[string]string)
    }

    msg := Message{Body: b, Headers: hJSON}

    err := mex.Publish(topic, msg)
    if err != nil {
      w.Write([]byte("Error publishing message, " + err.Error()))
    }
  }
}

func SubscribeHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    subscriberID := r.URL.Query().Get("sid")
    if err := mex.Subscribe(topic, subscriberID) ; err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("Error subscribing"))
    }
  }
}

func ConsumeHandler(mex MessageExchange) HandleFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    topic := r.URL.Query().Get("topic")
    subscriberID := r.URL.Query().Get("sid")

    amount := r.URL.Query().Get("amount")
    if amount == "" {
      amount = "1"
    }

    a, err := strconv.Atoi(amount)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Error converting amount to integer, "))
    }

    msgs, err := mex.Consume(topic, subscriberID, a)
    if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      w.Write([]byte("Error consuming messages"))
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
