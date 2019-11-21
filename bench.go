package gout

import (
	"time"
)

type Bench struct {
	durations  time.Duration
	number     int
	concurrent int

	g *routerGroup
}

func (b *Bench) Concurrent(c int) *Bench {
	b.concurrent = c
	return b
}

func (b *Bench) Number(n int) *Bench {
	b.number = n
	return b
}

func (b *Bench) Durations(d time.Duration) *Bench {
	b.durations = d
	return b
}

func (b *Bench) Do() {
}
