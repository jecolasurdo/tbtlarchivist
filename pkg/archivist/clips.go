package archivist

import (
	"context"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

func (a *API) getCuratedClipSource(ctx context.Context) (<-chan contracts.ClipInfo, error) {
	// similar process to episode handling, except clips are unique by name only
	panic("not implemented")
}

func (a *API) processCuratedClip(ctx context.Context, clip contracts.ClipInfo) error {
	panic("not implemented")
}
