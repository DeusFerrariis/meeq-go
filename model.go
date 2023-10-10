package main

type Message struct {
  Body string `json:"body"`
  Headers map[string]string `json:"headers, omitempty"`
}

func NewMessage(body string) Message {
  return Message{Body: body, Headers: make(map[string]string)}
}
