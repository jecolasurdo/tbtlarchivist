package utils

import (
	"log"
	"math/rand"
	"time"

	"github.com/gonum/stat/distuv"
)

// A Pacer provides a means of pausing program execution for a period of time.
type Pacer struct {
	basis        time.Duration
	lastCallTime time.Time
	jitter       distribution
}

type distribution interface {
	Rand() float64
}

// SetNormalPace returns a Pace where mu is the average wait time, sigma is the
// standard deviation of the wait time, and basis is the wait-time unit
// duration.
func SetNormalPace(mu, sigma float64, basis time.Duration) *Pacer {
	return &Pacer{
		basis:        basis,
		lastCallTime: time.Now(),
		jitter: distuv.Normal{
			Mu:     mu,
			Sigma:  sigma,
			Source: rand.New(rand.NewSource(time.Now().UnixNano())),
		},
	}
}

// SetUniformPace returns a pace where the wait time is uniformly distributed
// between min and max, and basis is the wait-time unit duration.
func SetUniformPace(min, max float64, basis time.Duration) *Pacer {
	return &Pacer{
		basis:        basis,
		lastCallTime: time.Now(),
		jitter: distuv.Uniform{
			Min:    min,
			Max:    max,
			Source: rand.New(rand.NewSource(time.Now().UnixNano())),
		},
	}
}

// Wait blocks until the time since the last call has exceeded a minimum pacing
// duration. The pace duration is centered around a mean wait time plus or
// minus a normally distributed jitter period.
func (p *Pacer) Wait() {
	paceDuration := time.Duration(p.jitter.Rand()) * p.basis
	paceTime := p.lastCallTime.Add(paceDuration)
	if time.Now().Before(paceTime) {
		waitDuration := paceTime.Sub(time.Now())
		log.Printf("Pacing (%v)...", waitDuration)
		time.Sleep(waitDuration)
	}
	p.lastCallTime = time.Now()
}
