package main

import (
  log "github.com/sirupsen/logrus"
  "github.com/redis/go-redis/v9"
  "context"
  "errors"
  "encoding/json"
)

type MessageExchange interface {
  NewTopic(name string) error
  Publish(topic string, msg Message) error
  Subscribe(topic string, subscriberID string) error
  Consume(topic string, subscriberID string, amount int) ([]Message, error)
  CheckSubscriptions(topic string, subscriberID string) error
}

type RedisClient struct {
  client *redis.Client
  ctx context.Context
}

func (rc *RedisClient) NewTopic(name string) error {
  log.Info("Creating new topic: " + name)

  if _, err := rc.client.RPush(rc.ctx, "topics", name).Result(); err != nil {
    log.Error("Error creating topic: ", err.Error())
    return errors.New("Error creating topic")
  }

  log.Info("Topic created successfully")

  return nil
}

func (rc *RedisClient) Publish(topic string, msg Message) error {
  log.Info("Publishing message to topic: " + topic)
  key := "topic:" + topic + ":messages"

  log.Info("Marshalling message to json")
  msgJSON, err := json.Marshal(&msg)
  if err != nil {
    log.Error("Error marshalling, " + err.Error())
    return errors.New("Error publishing message, " + err.Error())
  }

  log.Info("Publishing marshalled message to redis")
  if _, err := rc.client.RPush(rc.ctx, key, string(msgJSON)).Result(); err != nil {
    log.Error("Error publishing message, " + err.Error())
    return errors.New("Error publishing message")
  }

  log.Info("Message published successfully")
  return nil
}

func (rc *RedisClient) Subscribe(topic string, subscriberID string) error {
  log.Info("Subscribing to topic: " + topic)
  key := "topic:" + topic + ":subscribers"

  if _, err := rc.client.RPush(rc.ctx, key, subscriberID).Result(); err != nil {
    log.Error("Error subscribing to topic, " + err.Error())
    return errors.New("Error subscribing to topic")
  }

  log.Info("Subscribed to topic successfully")
  return nil
}

func (rc *RedisClient) Consume(topic string, subscriberID string, amount int) ([]Message, error) {
  log.Info("Consuming messages from topic: " + topic)
  if err := rc.CheckSubscriptions(topic, subscriberID); err != nil {
    return nil, errors.New("Error consuming messages")
  }

  key := "topic:" + topic + ":messages"

  var msgs []Message

  r := rc.client.LPopCount(rc.ctx, key, amount)
  if r.Err() == redis.Nil {
    log.Info("No messages to consume")
    return make([]Message, 0), nil
  }

  if r.Err() != nil {
    log.Error("Error consuming messages, " + r.Err().Error())
    return nil, errors.New("Error, retrieving messages failed")
  }

  for _, msg := range r.Val() {
    var m Message
    if err := json.Unmarshal([]byte(msg), &m); err != nil {
      log.Error("Error consuming messages, " + err.Error())  
    }
    msgs = append(msgs, m)
  }

  return msgs, nil
}

func (rc *RedisClient) CheckSubscriptions(topic string, subscriberID string) error {
  log.Info("Checking subscriptions for topic: " + topic)
  key := "topic:" + topic + ":subscribers"

  subscribers, err := rc.client.LRange(rc.ctx, key, 0, -1).Result()

  if err != nil {
    log.Error("Error checking subscriptions, " + err.Error())
    return errors.New("Error checking subscriptions")
  }

  for _, id := range subscribers {
    if id == subscriberID {
      return nil
    }
  }

  return errors.New("Subscriber is not subscribed to topic")
}
