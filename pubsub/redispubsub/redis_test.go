package redispubsub_test

import (
	"context"
	"testing"

	"github.com/plaenen/webx/pubsub/pubsubtest"
	"github.com/plaenen/webx/pubsub/redispubsub"
	"github.com/redis/go-redis/v9"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	container, err := tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		t.Fatalf("starting Redis container: %v", err)
	}
	t.Cleanup(func() { container.Terminate(ctx) })

	endpoint, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("getting Redis connection string: %v", err)
	}

	opts, err := redis.ParseURL(endpoint)
	if err != nil {
		t.Fatalf("parsing Redis URL: %v", err)
	}

	client := redis.NewClient(opts)
	t.Cleanup(func() { client.Close() })

	ps := redispubsub.New(client)
	t.Cleanup(func() { ps.Close() })

	pubsubtest.RunSuite(t, ps)
}
