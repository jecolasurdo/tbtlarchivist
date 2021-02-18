package analystiface

import (
	"context"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// AnalystAPI is anything that can take a PendingResearchItem, conduct an
// analysis, and return a channel of CompletedResearchItem in response.
type AnalystAPI interface {
	Run(context.Context, *contracts.PendingResearchItem) (<-chan *contracts.CompletedResearchItem, <-chan error)
}
