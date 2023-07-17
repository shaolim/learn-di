package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	tm := time.Date(2019, 6, 17, 0, 0, 0, 0, time.UTC)

	// mock Now
	Now = func() time.Time {
		return tm
	}

	// restore mock Now
	defer func() {
		Now = func() time.Time {
			return time.Now()
		}
	}()

	actual := Now()
	assert.Equal(t, tm, actual)
}
