package webclient

type Streamer interface {
	Do(func(data []byte))
}

type streamer struct {
	C chan []byte
}

func (s *streamer) Do(fn func(data []byte)) {
	for line := range s.C {
		fn(line)
	}
}
