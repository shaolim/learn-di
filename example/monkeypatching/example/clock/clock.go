package clock

import "time"

type TimeFn = func() time.Time

var Now TimeFn = time.Now
