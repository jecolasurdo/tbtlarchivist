package pacer

import (
	"log"
	"time"

	"github.com/gonum/stat/distuv"
)

// Pace provides a means of pausing program execution for a period of time.
type Pace struct {
	basis     time.Duration
	startTime time.Time
	jitter    distuv.Normal
}

// SetPace returns a Pace where mu is the average wait time, sigma is the
// standard deviation of the wait time, and basis is the wait-time unit
// duration.
func SetPace(mu, sigma float64, basis time.Duration) Pace {
	return Pace{
		basis:     basis,
		startTime: time.Now(),
		jitter: distuv.Normal{
			Mu:    mu,
			Sigma: sigma,
		},
	}
}

// Wait blocks until the time since the last call has exceeded a minimum pacing
// duration. The pace duration is centered around a mean wait time plus or
// minus a normally distributed jitter period.
func (p Pace) Wait() {
	paceDuration := time.Duration(p.jitter.Rand()) * p.basis
	paceTime := p.startTime.Add(paceDuration)
	if time.Now().Before(paceTime) {
		waitDuration := paceTime.Sub(time.Now())
		log.Printf("Pacing (%v)...", waitDuration)
		time.Sleep(waitDuration)
	}
	p.startTime = time.Now()
}
