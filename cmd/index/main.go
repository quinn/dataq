package main

import (
	"context"
	"log"

	"go.quinn.io/dataq/boot"
)

func main() {
	b, err := boot.New()
	if err != nil {
		log.Fatalf("Failed to initialize boot: %v", err)
	}

	if err := b.Index.Rebuild(context.Background()); err != nil {
		log.Fatalf("Failed to rebuild index: %v", err)
	}
}
