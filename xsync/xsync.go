package xsync

type ContextMutex struct {
	ch chan struct{}
}

func NewContextMutex() *ContextMutex {
	return &ContextMutex{
		ch: make(chan struct{}, 1),
	}
}

func (cm *ContextMutex) TryLock() bool {
	select {
	case cm.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (cm *ContextMutex) Lock(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("failed to acquire lock: %w", ctx.Err())
	case cm.ch <- struct{}{}:
		return nil
	}
}

func (cm *ContextMutex) Unlock() {
	select {
	case <-cm.ch:
	default:
		panic("xsync.ContextMutex: Unlock before Lock")
	}
}
