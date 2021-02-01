package archivist

import (
	"context"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

func (a *API) getCompletedResearchSource(ctx context.Context) (<-chan contracts.ResearchComplete, error) {
	// upsert research and update leases if applicable
	panic("not implemented")
}

func (a *API) processCompletedResearch(ctx context.Context, completedResearch contracts.ResearchComplete) error {
	panic("not implemented")
}
