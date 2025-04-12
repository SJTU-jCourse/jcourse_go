package embedding

import (
	"context"
	"log"
	"time"
)

func InitVectorBase() {
	initCtx, cancelInit := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelInit()
	if err := InitVectorStore(initCtx); err != nil {
		log.Fatalf("FATAL: Failed to initialize vector store client: %v", err)
	}
	log.Println("Vector store client initialized.")
}
