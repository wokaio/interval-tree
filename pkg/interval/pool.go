package fdinterval

import (
	"errors"
	"fmt"

	cmap "github.com/orcaman/concurrent-map"
)

type IntervalPool struct {
	imap cmap.ConcurrentMap
}

func NewIntervalPool() *IntervalPool {
	intervalPool := &IntervalPool{
		imap: cmap.New(),
	}
	return intervalPool
}

func (sender *IntervalPool) GetIntervalPtr(key string) (*IntervalNode, error) {
	if tmp, ok := sender.imap.Get(key); ok {
		return tmp.(*IntervalNode), nil
	}
	return nil, errors.New(fmt.Sprintf("IntervalNode store is empty"))
}

func (sender *IntervalPool) SetIntervalPtr(key string, root *IntervalNode) {
	sender.imap.Set(key, root)
}

func (sender *IntervalPool) RemoveIntervalPtr(key string) bool {
	ok := sender.imap.Has(key)
	if ok {
		sender.imap.Remove(key)
		return true
	}
	return false
}
