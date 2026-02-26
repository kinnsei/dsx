package chanpubsub_test

import (
	"testing"

	"github.com/plaenen/webx/pubsub/chanpubsub"
	"github.com/plaenen/webx/pubsub/pubsubtest"
)

func TestSuite(t *testing.T) {
	ps := chanpubsub.New()
	defer ps.Close()
	pubsubtest.RunSuite(t, ps)
}
