package webclient

type Streamer interface {
	Do(func(data interface{}))
}

type streamer struct {
	C chan interface{}
}

func (s *streamer) Do(fn func(data interface{})) {
	for line := range s.C {
		fn(line)
	}
}
