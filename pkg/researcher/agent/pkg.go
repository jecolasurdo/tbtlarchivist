package agent

import "github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"

// A ResearchAgent is responsible for gathering a pending work item from the
// pending work queue, spawning an Analyst sub-process, communicating with that
// process, and reporting completed research results back to the completed work
// queue.
type ResearchAgent struct {
	Errors <-chan error
	Done   <-chan struct{}
}

// StartResearchAgent initializes a research agent. The agent will attempt to
// consume a pending work item from the pending work queue. If there is no work
// found on the queue, the agent will exit. Otherwise, the agent will spawn an
// Analyst process, and assign the work to that process.  As the Analyst
// completes its work, it is reported back to the Agent, who then forwards the
// results to the completed work queue.
func StartResearchAgent(ctx, queue messagebus.SenderReceiver) *ResearchAgent {
	errorSource := make(chan error)
	done := make(chan struct{})
	go func() {
		defer close(errorSource)
		defer close(done)

	}()

	return &ResearchAgent{
		Errors: errorSource,
		Done:   done,
	}
}
