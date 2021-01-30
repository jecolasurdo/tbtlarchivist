package messagebusiface

type MessageBus interface {
	Send([]byte) error
	Receive() ([]byte, error)
}
