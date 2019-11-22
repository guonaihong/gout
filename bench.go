package gout

import (
	"time"
)

type Bench struct {
	durations  time.Duration
	number     int
	concurrent int
	rate       int

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

func (b *Bench) Rate(rate int) *Bench {
	b.rate = rate
	return b
}

func (b *Bench) Durations(d time.Duration) *Bench {
	b.durations = d
	return b
}

func (b *Bench) Do() {
}
