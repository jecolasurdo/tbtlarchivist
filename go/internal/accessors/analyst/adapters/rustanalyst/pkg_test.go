package rustanalyst_test

import (
	"context"
	"log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst/adapters/rustanalyst"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/mocks/accessors/mock_analyst"
)

func Test_AdapterRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	readCloser := mock_analyst.NewMockReadCloser(ctrl)
	writeCloser := mock_analyst.NewMockWriteCloser(ctrl)

	cmd := mock_analyst.NewMockCommand(ctrl)
	cmd.EXPECT().Start().Return(nil).Times(1)
	cmd.EXPECT().Wait().Return(nil).Times(1)
	cmd.EXPECT().StdinPipe().Return(writeCloser, nil).Times(1)
	cmd.EXPECT().StdoutPipe().Return(readCloser, nil).Times(1)
	cmdBuilder := mock_analyst.NewMockCommandBuilder(ctrl)
	cmdBuilder.EXPECT().CommandContext(gomock.Any(), gomock.Any()).
		Return(cmd).
		Times(1)

	adapter := &rustanalyst.Adapter{
		PathResolver: func() (string, error) { return "", nil },
		CmdBuilder:   cmdBuilder,
	}

	pendingResearchItem := &contracts.PendingResearchItem{}
	adapter.Run(context.Background(), pendingResearchItem)

	for {
		select {
		case item, open := <-adapter.CompletedWorkItems():
			if !open {
				break
			}
			log.Println(item)
		case err, open := <-adapter.Errors():
			if !open {
				break
			}
			log.Println(err)
		case <-adapter.Done():
			log.Println("Done")
			return
		}
	}
}