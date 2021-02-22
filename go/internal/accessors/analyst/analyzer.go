package analyst

import (
	"context"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
)

// An Analyzer is anything that can take a PendingResearchItem, conduct an
// analysis, and return a channel of CompletedResearchItem in response.
type Analyzer interface {
	Run(context.Context, *contracts.PendingResearchItem)
	Errors() <-chan error
	CompletedWorkItems() <-chan *contracts.CompletedResearchItem
	Done() <-chan struct{}
}
