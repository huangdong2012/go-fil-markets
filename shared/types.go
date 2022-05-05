package shared

import (
	"io"
	"strings"
	"sync"
)

// TipSetToken is the implementation-nonspecific identity for a tipset.
type TipSetToken []byte

// Unsubscribe is a function that gets called to unsubscribe from (storage|retrieval)market events
type Unsubscribe func()

// ReadSeekStarter implements io.Reader and allows the caller to seek to
// the start of the reader
type ReadSeekStarter interface {
	io.Reader
	SeekStart() error
}

func NewStringReadSeekStarter(str string) ReadSeekStarter {
	srss := &StringReadSeekStarter{
		Str:   str,
		rdrRW: &sync.RWMutex{},
	}
	_ = srss.SeekStart()
	return srss
}

type StringReadSeekStarter struct {
	Str string

	rdrRW *sync.RWMutex
	rdr   io.Reader
}

func (s *StringReadSeekStarter) Read(p []byte) (n int, err error) {
	s.rdrRW.RLock()
	defer s.rdrRW.RUnlock()

	return s.rdr.Read(p)
}

func (s *StringReadSeekStarter) SeekStart() error {
	s.rdrRW.Lock()
	defer s.rdrRW.Unlock()

	s.rdr = strings.NewReader(s.Str)
	return nil
}
