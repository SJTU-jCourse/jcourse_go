package embedding

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

func InitVectorBase(db *gorm.DB) {
	err := ensurePgVectorExtension(db)
	if err != nil {
		log.Fatalf("FATAL: Failed to enable pgvector: %v", err)
	}

	initCtx, cancelInit := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelInit()
	if err := InitVectorStore(initCtx); err != nil {
		log.Fatalf("FATAL: Failed to initialize vector store client: %v", err)
	}
	log.Println("Vector store client initialized.")
}

func ensurePgVectorExtension(db *gorm.DB) error {
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS vector;")
	if result.Error != nil {
		log.Printf("Error creating pgvector extension: %v", result.Error)
		return result.Error
	}
	log.Println("Checked/created pgvector extension.")
	return nil
}
