package amqpserver

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/oklog/run"
	"log/slog"
)

type HandlerFunc func(ctx context.Context, msg *message.Message) error

func (f HandlerFunc) ServeMsg(ctx context.Context, msg *message.Message) error {
	return f(ctx, msg)
}

type Handler interface {
	ServeMsg(ctx context.Context, msg *message.Message) error
}

type Server struct {
	handler Handler
}

func NewServer(handler Handler) *Server {
	return &Server{handler: handler}
}

func (l *Server) ListenAndServe(msgs <-chan *message.Message) {
	for msg := range msgs {
		err := l.handler.ServeMsg(msg.Context(), msg)
		if err != nil {
			msg.Nack()
			continue
		}
		msg.Ack()
	}
}

type Listener interface {
	ListenAndServe(context.Context) error
	Shutdown() error
}

func RunWithGroup(
	g *run.Group,
	amqpURI string,
	subscriberFactory func(subscriber *amqp.Subscriber) Listener,
) {
	var listener Listener
	g.Add(func() error {
		subscriber, err := amqp.NewSubscriber(
			amqp.NewDurableQueueConfig(amqpURI),
			watermill.NewSlogLogger(slog.Default()),
		)
		if err != nil {
			return err
		}
		listener = subscriberFactory(subscriber)
		return listener.ListenAndServe(context.Background())
	}, func(error) {
		if listener != nil {
			_ = listener.Shutdown()
		}
	})
}
