package rustanalyst_test

import (
	"context"
	"encoding/binary"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst/adapters/rustanalyst"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/mocks/accessors/mock_analyst"
	"google.golang.org/protobuf/runtime/protoiface"
)

func Test_AdapterRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	readCloser := mock_analyst.NewMockReadCloser(ctrl)
	// completedResearchItem := &contracts.CompletedResearchItem{}
	// framedResearchItem := mustFrame(completedResearchItem)
	readCloser.EXPECT().Read(gomock.Any()).Return(0, nil).AnyTimes()

	writeCloser := mock_analyst.NewMockWriteCloser(ctrl)
	pendingResearchItem := &contracts.PendingResearchItem{}
	marshaledResearchItem := mustMarshal(pendingResearchItem)
	writeCloser.EXPECT().Write(marshaledResearchItem).Return(len(marshaledResearchItem), nil).Times(1)
	writeCloser.EXPECT().Close().Return(nil).Times(1)

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	adapter.Run(ctx, pendingResearchItem)
	defer cancel()

	for {
		select {
		case item, open := <-adapter.CompletedWorkItems():
			if !open {
				break
			}
			log.Println("item", item)
		case err, open := <-adapter.Errors():
			if !open {
				break
			}
			log.Println("error", err)
		case <-adapter.Done():
			log.Println("Done")
			return
		}
	}
}

func mustMarshal(msg protoiface.MessageV1) []byte {
	protoBytes, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return protoBytes
}

func mustFrame(msg protoiface.MessageV1) []byte {
	protoBytes := mustMarshal(msg)
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(len(protoBytes)))
	return append(bs, protoBytes...)
}
