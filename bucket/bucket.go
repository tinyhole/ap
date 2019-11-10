package bucket

import (
	"errors"
	"github.com/tinyhole/ap/transport"
	"sync"
)

var (
	ErrNotFoundSocket = errors.New("error not found socket")
)

type Bucket interface {
	Add(int64, transport.Socket)
	Remove(int64)
	Get(int64) (transport.Socket, error)
}

type socketBucket struct {
	bucket map[int64]transport.Socket
	mu     sync.Mutex
}

var (
	DefaultSocketBucket = &socketBucket{
		bucket: make(map[int64]transport.Socket),
	}
)

func (s *socketBucket) Add(fid int64, socket transport.Socket) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bucket[fid] = socket
}

func (s *socketBucket) Remove(fid int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bucket[fid]; ok {
		delete(s.bucket, fid)
	}
}

func (s *socketBucket) Get(fid int64) (transport.Socket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s, ok := s.bucket[fid]; ok {
		return s, nil
	}
	return nil, ErrNotFoundSocket
}
