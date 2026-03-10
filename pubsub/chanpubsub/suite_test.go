package chanpubsub_test

import (
	"testing"

	"github.com/laenen-partners/dsx/pubsub/chanpubsub"
	"github.com/laenen-partners/dsx/pubsub/pubsubtest"
)

func TestSuite(t *testing.T) {
	ps := chanpubsub.New()
	defer ps.Close()
	pubsubtest.RunSuite(t, ps)
}
