package rustanalyst

import "github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/analyst"

// The Adapter spawns a child analyst-rust process, and marshals messages between
// the caller and the child processes.
type Adapter struct{}

var _ analyst.Analyzer = (*Adapter)(nil)
