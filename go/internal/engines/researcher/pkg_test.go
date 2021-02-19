package researcher_test

import (
	"context"
	"runtime"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jecolasurdo/tbtlarchivist/mocks/accessors/mock_messagebus"
	"github.com/jecolasurdo/tbtlarchivist/mocks/accessors/mock_messagebus/mock_acknowledger"
	"github.com/jecolasurdo/tbtlarchivist/mocks/researcher/mock_agent/mock_analystiface"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus/messagebustypes"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
	"github.com/jecolasurdo/tbtlarchivist/pkg/researcher/agent"
	"google.golang.org/protobuf/proto"
)

func Test_AgentHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Dependencies
	ctx := context.Background()
	pendingQueue := mock_messagebus.NewMockReceiver(ctrl)
	completedQueue := mock_messagebus.NewMockSender(ctrl)
	analyst := mock_analystiface.NewMockAnalystAPI(ctrl)

	// pendingQueue.Receive behavior/expectation
	acknack := mock_acknowledger.NewMockAckNack(ctrl)
	pri := &contracts.PendingResearchItem{
		// We need at least one field populated. Otherwise protobuf will notice
		// that the struct is a zero value and return an empty message.  That
		// would result in the agent thinking there is no work to do, which
		// isn't the behavior we're exercising.
		LeaseId: "FakeLeaseID",
	}
	priBytes, err := proto.Marshal(pri)
	if err != nil {
		panic(err)
	}
	inboundMsg := &messagebustypes.Message{
		Acknowledger: acknack,
		Body:         priBytes,
	}

	acknack.EXPECT().Ack().Times(1)
	acknack.EXPECT().Nack(gomock.Any()).Times(0)
	pendingQueue.EXPECT().Receive().Return(inboundMsg, nil).Times(1)

	// Analyst behavior/expectations
	completedWorkSrc := make(chan *contracts.CompletedResearchItem)
	analystErrSrc := make(chan error)
	analyst.EXPECT().Run(gomock.Any(), gomock.Any()).Return(completedWorkSrc, analystErrSrc).Times(1)
	close(completedWorkSrc)
	close(analystErrSrc)

	// completedQueue.Send behavior/expectations
	completedQueue.EXPECT().Send(gomock.Any()).Return(nil).Times(0)

	// Run SUT
	researchAgent := agent.StartResearchAgent(ctx, pendingQueue, completedQueue, analyst)
	for {
		select {
		case err, open := <-researchAgent.Errors:
			if !open {
				break
			}
			t.Fatal(err)
		case <-researchAgent.Done:
			return
		default:
			runtime.Gosched()
		}
	}
}
