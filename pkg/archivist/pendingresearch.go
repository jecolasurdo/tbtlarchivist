package archivist

import (
	"context"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

func (a *API) getPendingResearchSource(ctx context.Context) (<-chan contracts.ResearchPending, error) {
	// check to see how many consumers there are for a queue
	// compare the consumer count to the message count
	// Then determine how much work to create, ie consumerCount - messageCount
	// Create that much work (including leases) and send it to the queue
	panic("not implemented")
}

func (a *API) processPendingResearch(ctx context.Context, pendingResearch contracts.ResearchPending) error {
	panic("not implemented")
}
