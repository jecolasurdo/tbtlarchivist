package curatoriface

import "google.golang.org/protobuf/reflect/protoreflect"

type Curator interface {
	Curate() (<-chan protoreflect.ProtoMessage, <-chan error)
}
