// sequencer is a package returning unique seq no to be used by all routines
// tbd. as it is not supporting multiproces mode as the seq no must be
// persistent and shared between instances

package sequencer

import (
	"sync"
)

var (
	seqNo int
	mutex sync.Mutex
)

func NextSeqNo() (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	seqNo := seqNo + 1

	return seqNo, nil
}
