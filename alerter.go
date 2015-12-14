package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Alerter struct {
	pool      *redis.Pool
	AlertChan chan *Alert
}

func NewAlerter(server, password string) (*Alerter, error) {
	return &Alerter{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		AlertChan: make(chan *Alert),
	}, nil
}

func (a *Alerter) Start() {
	go a.listenForChecks()
}

func (a *Alerter) GetChan() chan *Alert {
	return a.AlertChan
}

func (a *Alerter) listenForChecks() {
	for {
		select {
		case result := <-a.AlertChan:
			if err := a.processAlert(result); err != nil {
				log.Printf("ERROR: Unable to process alert, %s", err.Error())
			}
		case <-time.After(100 * time.Millisecond):
			// NOP, just breath
		}
	}
}

func (a *Alerter) processAlert(alert *Alert) error {
	if alert.Status != Resolved {
		value, _ := alert.Serialize()
		key, field, err := alert.prepAlertForHash()
		if err != nil {
			log.Printf("Unable to act on the following alert:\n%s\n", alert)
		} else {
			a.setHash(key, field, value)
		}
	} else {
		key, field, err := alert.prepAlertForHash()
		if err == nil {
			a.resolve(key, field)
		}
	}
	return nil
}

func (alert *Alert) prepAlertForHash() (string, string, error) {
	if alert.CapsuleID == "" {
		return "", "", errors.New("capsule_id required")
	} else if alert.DeploymentID == "" {
		return "", "", errors.New("deployment_id required")
	}
	latterPart := strings.TrimPrefix(alert.Name, alert.CapsuleName)
	checkName := strings.Replace(latterPart, "-", "", 1)
	key := fmt.Sprintf("%s:%s", alert.DeploymentID, alert.CapsuleID)
	return key, checkName, nil
}

func (a *Alerter) set(key string, value []byte, expire int) {
	conn := a.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", key, value); err != nil {
		fmt.Printf("ERROR: unable to set redis key %s\n", err.Error())
	}
}

func (a *Alerter) setWithExpires(key string, value []byte, expire int) {
	conn := a.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", key, value, "EX", expire); err != nil {
		fmt.Printf("ERROR: unable to set redis key %s\n", err.Error())
	}
}

func (a *Alerter) setHash(key, field string, value []byte) {
	conn := a.pool.Get()
	defer conn.Close()
	// expire the key after 30 minutes
	conn.Send("MULTI")
	conn.Send("HSET", key, field, value)
	conn.Send("EXPIRE", key, 60*30)
	_, err := redis.Ints(conn.Do("EXEC"))
	if err != nil {
		fmt.Printf("ERROR: unable to set redis key %s\n", err.Error())
	}
}

func (a *Alerter) resolve(key, field string) {
	conn := a.pool.Get()
	defer conn.Close()
	redis.Int(conn.Do("HDEL", key, field))
}

func (a *Alerter) get(key string) (interface{}, error) {
	conn := a.pool.Get()
	defer conn.Close()
	return conn.Do("GET", key)
}

func (a *Alerter) GetAll(key_search string) (map[string][]map[string]interface{}, error) {
	var alerts = make(map[string][]map[string]interface{})
	conn := a.pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", key_search))
	for _, v := range keys {
		capsule_id := strings.Split(v, ":")[1]
		current, _ := redis.StringMap(conn.Do("HGETALL", v))
		for _, alert := range current {
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(alert), &dat); err != nil {
				return alerts, err
			}
			alerts[capsule_id] = append(alerts[capsule_id], dat)
		}
	}
	return alerts, nil
}
