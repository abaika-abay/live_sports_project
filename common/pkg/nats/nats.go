package nats

import (
	"log"

	"github.com/nats-io/nats.go"
)

type NATS struct {
	Conn *nats.Conn
}

func NewNATS(url string) (*NATS, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to NATS")
	return &NATS{Conn: conn}, nil
}

func (n *NATS) Close() {
	if n.Conn != nil {
		n.Conn.Close()
	}
}
