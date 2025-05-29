package main

import (
	"context"
	"log"

	"github.com/abaika-abay/live_sports_project/common/pkg/config"
	"github.com/abaika-abay/live_sports_project/common/pkg/db"
	"github.com/abaika-abay/live_sports_project/common/pkg/logger"
	"github.com/abaika-abay/live_sports_project/common/pkg/nats"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	logger.InitLogger(cfg.Log.Level)
	mongoDB, err := db.NewMongoDB(cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		log.Fatalf("MongoDB error: %v", err)
	}
	defer mongoDB.Disconnect(context.Background())

	natsConn, err := nats.NewNATS(cfg.NATS.URL)
	if err != nil {
		log.Fatalf("NATS error: %v", err)
	}
	defer natsConn.Close()

	logger.InfoLogger.Println("All services initialized successfully")
}
