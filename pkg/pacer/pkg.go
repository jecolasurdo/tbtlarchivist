package pacer

import (
	"log"
	"math/rand"
	"time"

	"github.com/gonum/stat/distuv"
)

// Pace provides a means of pausing program execution for a period of time.
type Pace struct {
	basis        time.Duration
	lastCallTime time.Time
	jitter       *distuv.Normal
}

// SetPace returns a Pace where mu is the average wait time, sigma is the
// standard deviation of the wait time, and basis is the wait-time unit
// duration.
func SetPace(mu, sigma float64, basis time.Duration) Pace {
	return Pace{
		basis:        basis,
		lastCallTime: time.Now(),
		jitter: &distuv.Normal{
			Mu:     mu,
			Sigma:  sigma,
			Source: rand.New(rand.NewSource(time.Now().UnixNano())),
		},
	}
}

// Wait blocks until the time since the last call has exceeded a minimum pacing
// duration. The pace duration is centered around a mean wait time plus or
// minus a normally distributed jitter period.
func (p *Pace) Wait() {
	paceDuration := time.Duration(p.jitter.Rand()) * p.basis
	paceTime := p.lastCallTime.Add(paceDuration)
	if time.Now().Before(paceTime) {
		waitDuration := paceTime.Sub(time.Now())
		log.Printf("Pacing (%v)...", waitDuration)
		time.Sleep(waitDuration)
	}
	p.lastCallTime = time.Now()
}
