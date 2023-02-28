package xcontext

type Mutex struct {
	ch chan struct{}
}

func NewMutex() *Mutex {
	return &Mutex{
		ch: make(chan struct{}, 1),
	}
}

func (m *Mutex) TryLock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (m *Mutex) Lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to acquire lock: %w", ctx.Err())
	case m.ch <- struct{}{}:
		return nil
	}
}

func (m *Mutex) Unlock() {
	select {
	case <-m.ch:
	default:
		panic("xcontext.Mutex: Unlock before Lock")
	}
}
