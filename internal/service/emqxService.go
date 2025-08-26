package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/Eatriceeveryday/data-stream-service/internal/config"
	"github.com/Eatriceeveryday/data-stream-service/internal/entities"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type EMQXService struct {
	client          mqtt.Client
	publishInterval atomic.Int32
	intervalCh      chan time.Duration
	ctx             context.Context
	cancel          context.CancelFunc
	sensorType      string
	id1             string
	id2             int
	key             string
}

func NewEmqxService(client mqtt.Client, cfg *config.Config) *EMQXService {
	ctx, cancel := context.WithCancel(context.Background())
	s := &EMQXService{
		client:     client,
		ctx:        ctx,
		cancel:     cancel,
		sensorType: cfg.SensorType,
		intervalCh: make(chan time.Duration, 1),
		id1:        cfg.ID1,
		id2:        cfg.ID2,
		key:        cfg.ApiKey,
	}
	s.publishInterval.Store(5)
	return s
}

func (s *EMQXService) publishMessage(topic string, payload []byte) error {
	token := s.client.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}

func (s *EMQXService) StartPublishing(ctx context.Context) {
	go func() {
		current := time.Duration(s.publishInterval.Load()) * time.Second
		if current <= 0 {
			current = 1 * time.Second
		}
		ticker := time.NewTicker(current)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Publisher Stoped")
				return

			case newDur := <-s.intervalCh:
				if newDur <= 0 {
					newDur = 1 * time.Second
				}
				ticker.Stop()
				ticker = time.NewTicker(newDur)
				fmt.Printf("Publisher interval changed to %v\n", newDur)

			case <-ticker.C:
				msg := entities.Message{
					Value:      math.Round(rand.Float64()*100*100) / 100,
					SensorType: s.sensorType,
					ID1:        s.id1,
					ID2:        s.id2,
					Key:        s.key,
					Timestamp:  time.Now().UTC().Format(time.RFC3339),
				}

				payload, _ := json.Marshal(msg)
				err := s.publishMessage("sensors", payload)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}

func (s *EMQXService) ChangeInterval(newInterval int32) {
	if newInterval <= 0 {
		return
	}
	s.publishInterval.Store(newInterval)
	d := time.Duration(newInterval) * time.Second
	select {
	case s.intervalCh <- d:
	default:
		select {
		case <-s.intervalCh:
		default:
		}
		select {
		case s.intervalCh <- d:
		default:
		}
	}
}
