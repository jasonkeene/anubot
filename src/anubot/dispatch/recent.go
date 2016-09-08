package dispatch

type recent struct {
	head     int
	capacity int
	fifo     []*string
	set      map[string]struct{}
}

func newRecent(capacity int) *recent {
	return &recent{
		capacity: capacity,
		fifo:     make([]*string, capacity),
		set:      make(map[string]struct{}),
	}
}

func (r *recent) insert(value string) {
	old := r.fifo[r.head]
	if old != nil {
		delete(r.set, *old)
	}
	r.fifo[r.head] = &value
	r.set[value] = struct{}{}
	r.head = (r.head + 1) % r.capacity
}

func (r *recent) lookup(v string) bool {
	_, ok := r.set[v]
	return ok
}
