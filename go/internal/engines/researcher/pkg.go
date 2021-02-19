package researcher

import (
	"context"
	"log"
	"runtime"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/analyst"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
	"github.com/jecolasurdo/tbtlarchivist/pkg/utils"
	"google.golang.org/protobuf/proto"
)

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
func StartResearchAgent(ctx context.Context, pendingResearchQueue messagebus.Receiver, completedWorkQueue messagebus.Sender, analyzer analyst.Analyzer) *ResearchAgent {
	utils.PanicIfNil(pendingResearchQueue, completedWorkQueue, analyzer)

	errorSource := make(chan error)
	done := make(chan struct{})
	go func() {
		defer close(errorSource)
		defer close(done)

		msg, err := pendingResearchQueue.Receive()
		if err != nil {
			errorSource <- err
			return
		}

		if msg == nil || len(msg.Body) == 0 {
			log.Println("No pending-research to do.")
			return
		}

		pendingResearchItem := new(contracts.PendingResearchItem)
		err = proto.Unmarshal(msg.Body, pendingResearchItem)
		if err != nil {
			errorSource <- err
			err := msg.Acknowledger.Nack(true)
			if err != nil {
				errorSource <- err
			}
			return
		}

		err = msg.Acknowledger.Ack()
		if err != nil {
			errorSource <- err
			return
		}

		completedWorkSource, analystErrorSource := analyzer.Run(ctx, pendingResearchItem)
		utils.PanicIfNil(completedWorkSource, analystErrorSource)

		completedWorkSrcOpen, analystErrorSrcOpen := true, true
		for completedWorkSrcOpen || analystErrorSrcOpen {
			select {
			case completedWorkItem, open := <-completedWorkSource:
				if !open {
					completedWorkSrcOpen = false
					break
				}
				cwiBytes, err := proto.Marshal(completedWorkItem)
				if err != nil {
					errorSource <- err
					break
				}
				err = completedWorkQueue.Send(cwiBytes)
				if err != nil {
					errorSource <- err
				}
			case analystErr, open := <-analystErrorSource:
				if !open {
					analystErrorSrcOpen = false
					break
				}
				errorSource <- analystErr
			default:
				runtime.Gosched()
			}
		}

	}()

	return &ResearchAgent{
		Errors: errorSource,
		Done:   done,
	}
}
