package archivist

import (
	"context"
)

// The PendingReseachWorker is a bit different from the other three workers
// in that this worker is the archivists sole message producer. The
// PendingResearchWorker determines what work should be done by the researchers
// as well as how much work should be queued up for the researcher pool. I'm
// not sure if it is InitializeDataStream's responsibility to create the individual
// pending-work-items, or if that is the responsibility of the ProcessDatum
// method. Essentially, the tasks are as follows:
// 1) determine how much work to generate
// 2) determine which work to generate
// 3) deliver the messages for each pending-work-item to the message bus
// Here's what we'll do. We'll defer to leaving everything one-to-one on the channels.
// IntializeDataStream will produce individual pending-work-items. Each of which
// will be placed on the DataStream channel,
// The poller will pick them off one at a time and send them to the ProcessDatum
// method one by one. This will allow the poller to operate on quanta such that
// so it can more fairly schedule time across all of the workers.

// A PendingResearchWorker can initialize a stream of pending research, and
// supplies a method for processing each pending work item.
type PendingResearchWorker struct{}

// InitializeDataStream determines what research work needs to be done, as well
// as how much capacity there is to do that work.
func (c *PendingResearchWorker) InitializeDataStream(ctx context.Context) (<-chan interface{}, error) {
	panic("not implemented")
}

// ProcessDatum takes an order for how much work to create, and processes it.
func (c *PendingResearchWorker) ProcessDatum(ctx context.Context, datum interface{}) error {
	// check to see how many consumers there are for a queue
	// compare the consumer count to the message count
	// Then determine how much work to create, ie consumerCount - messageCount
	// Create that much work (including leases) and send it to the queue
	panic("not implemented")
}
