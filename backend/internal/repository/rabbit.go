package repository

import (
	"campusclaw-plus/internal/domain"
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"net"
	"time"
)

type Rabbit struct {
	conn     *amqp.Connection
	AckDelay time.Duration
}

func OpenRabbit(address string) (*Rabbit, error) {
	c, e := amqp.DialConfig(address, amqp.Config{Heartbeat: 10 * time.Second, Dial: amqp.DefaultDial(5 * time.Second)})
	if e != nil {
		return nil, errors.New("broker unavailable")
	}
	r := &Rabbit{conn: c}
	ch, e := c.Channel()
	if e != nil {
		c.Close()
		return nil, e
	}
	defer ch.Close()
	for _, q := range []string{"campusclaw.work", "campusclaw.dead"} {
		if _, e = ch.QueueDeclare(q, true, false, false, false, amqp.Table{"x-queue-type": "quorum"}); e != nil {
			c.Close()
			return nil, e
		}
	}
	return r, nil
}
func (r *Rabbit) Close() error { return r.conn.Close() }
func (r *Rabbit) Publish(ctx context.Context, v domain.Outbox) error {
	ch, e := r.conn.Channel()
	if e != nil {
		return e
	}
	defer ch.Close()
	if e = ch.Confirm(false); e != nil {
		return e
	}
	confirm := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	returned := ch.NotifyReturn(make(chan amqp.Return, 1))
	raw, _ := json.Marshal(map[string]any{"schema_version": 1, "event_id": v.ID, "task_id": v.TaskID, "request_id": v.TaskID})
	queue := "campusclaw.work"
	if v.Kind == "dead" {
		queue = "campusclaw.dead"
	}
	if e = ch.PublishWithContext(ctx, "", queue, true, false, amqp.Publishing{DeliveryMode: amqp.Persistent, ContentType: "application/json", MessageId: v.ID, Body: raw}); e != nil {
		return e
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-returned:
		return errors.New("unroutable message")
	case v, ok := <-confirm:
		if !ok || !v.Ack {
			return errors.New("broker publish rejected")
		}
		select {
		case <-returned:
			return errors.New("unroutable message")
		default:
			return nil
		}
	}
}
func (r *Rabbit) Consume(ctx context.Context, handle func(context.Context, string) error) error {
	ch, e := r.conn.Channel()
	if e != nil {
		return e
	}
	defer ch.Close()
	if e = ch.Qos(1, 0, false); e != nil {
		return e
	}
	messages, e := ch.Consume("campusclaw.work", "", false, false, false, false, nil)
	if e != nil {
		return e
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-messages:
			if !ok {
				return errors.New("broker consumer closed")
			}
			var v struct {
				Version   int    `json:"schema_version"`
				EventID   string `json:"event_id"`
				TaskID    string `json:"task_id"`
				RequestID string `json:"request_id"`
			}
			if len(d.Body) > 4096 || json.Unmarshal(d.Body, &v) != nil || v.Version != 1 || len(v.TaskID) != 32 || len(v.EventID) != 32 {
				c, cancel := context.WithTimeout(ctx, 5*time.Second)
				e = r.Publish(c, domain.Outbox{ID: NewID(), Kind: "dead"})
				cancel()
				if e != nil {
					return e
				}
				if r.AckDelay > 0 {
					select {
					case <-ctx.Done():
						return nil
					case <-time.After(r.AckDelay):
					}
				}
				if e = d.Ack(false); e != nil {
					return e
				}
				continue
			}
			e = handle(ctx, v.TaskID)
			if ctx.Err() != nil {
				return nil
			}
			if e != nil {
				time.Sleep(time.Second)
				if e = d.Nack(false, true); e != nil {
					return e
				}
			} else {
				if r.AckDelay > 0 {
					select {
					case <-ctx.Done():
						return nil
					case <-time.After(r.AckDelay):
					}
				}
				if e = d.Ack(false); e != nil {
					return e
				}
			}
		}
	}
}

type BrokerProbe struct{ Address string }

func (p BrokerProbe) Ping(ctx context.Context) error {
	c, e := amqp.DialConfig(p.Address, amqp.Config{Dial: func(network, address string) (net.Conn, error) {
		c, e := (&net.Dialer{}).DialContext(ctx, network, address)
		if e == nil {
			if deadline, ok := ctx.Deadline(); ok {
				c.SetDeadline(deadline)
			}
		}
		return c, e
	}})
	if e != nil {
		return errors.New("broker unavailable")
	}
	return c.Close()
}
