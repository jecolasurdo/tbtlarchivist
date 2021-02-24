package utils

type Backoff struct{}

func (b *Backoff) Wait() error {
	panic("not implemented")
}

func (b *Backoff) Reset() {
	panic("not implmeneted")
}
