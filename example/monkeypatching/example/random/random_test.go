package random

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var randomizer = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func TestRandomInt(t *testing.T) {
	defer func(random *rand.Rand) {
		randomizer = random
	}(randomizer)

	randomizer = rand.New(&stubSource{})
	result := randomizer.Int63()
	assert.Equal(t, int64(123), result)
}

type stubSource struct {
}

func (sr *stubSource) Int63() int64 { return 123 }

func (sr *stubSource) Seed(seed int64) {}
