package datastore

import "github.com/jecolasurdo/tbtlarchivist/pkg/contracts"

type DataStorer interface {
	UpsertClipInfo(contracts.ClipInfo) error
	UpsertEpisodeInfo(contracts.EpisodeInfo) error
}
