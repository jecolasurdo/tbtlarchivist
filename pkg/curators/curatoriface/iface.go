package curatoriface

type Curator interface {
	Curate() (<-chan interface{}, <-chan error)
}
