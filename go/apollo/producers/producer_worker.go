package producers

import (
	"fmt"
	"sync"
)

// ProducerWorker publishes records to a channel
type ProducerWorker struct {
	id      int
	records []map[string]interface{}
	channel chan<- map[string]interface{}
	wg      *sync.WaitGroup
}

// NewProducerWorker creates a new producer worker
func NewProducerWorker(id int, records []map[string]interface{}, channel chan<- map[string]interface{}, wg *sync.WaitGroup) *ProducerWorker {
	return &ProducerWorker{
		id:      id,
		records: records,
		channel: channel,
		wg:      wg,
	}
}

// Start begins producing records to the channel
func (pw *ProducerWorker) Start() {
	go func() {
		defer pw.wg.Done()
		fmt.Printf("Producer Worker %d: Starting to produce %d records\n", pw.id, len(pw.records))

		for _, record := range pw.records {
			pw.channel <- record
			fmt.Printf("Producer Worker %d: Published record: %v\n", pw.id, record)
		}

		fmt.Printf("Producer Worker %d: Finished producing records\n", pw.id)
	}()
}
