package mempubsub_test

import (
	"github.com/areknoster/attendgo/test/testsubscribers"
	"testing"
	"time"

	"github.com/areknoster/attendgo/domain"
	. "github.com/areknoster/attendgo/mempubsub"
	"github.com/stretchr/testify/assert"
)

func TestPubSub(t *testing.T) {
	ps := NewPubSub(Config{BufSize: 3})
	defer ps.Close()

	sb1 := &testsubscribers.StringsBuilder{}
	ps.Register(sb1)
	ps.Publish(domain.EventKeyClicked{Glyph: 'a'})
	time.Sleep(time.Millisecond) // this is a race condition, but safe enough
	assert.EqualValues(t, "a", sb1.String())

	sb2 := &testsubscribers.StringsBuilder{}
	ps.Register(sb2)

	ps.Publish(domain.EventKeyClicked{Glyph: 'b'})
	time.Sleep(time.Millisecond)
	assert.EqualValues(t, "ab", sb1.String())
	assert.EqualValues(t, "b", sb2.String())
}
