package mock

import "sync"

type Logger struct {
	infos  []string
	errors []string
	sync.RWMutex
}

func (l *Logger) Info(i string, _ ...any) {
	l.Lock()
	defer l.Unlock()

	l.infos = append(l.infos, i)
}

func (l *Logger) Error(e string, _ ...any) {
	l.Lock()
	defer l.Unlock()

	l.errors = append(l.errors, e)
}

func (l *Logger) GotInfos() []string {
	l.RLock()
	defer l.RUnlock()

	return l.infos
}

func (l *Logger) GotErrors() []string {
	l.RLock()
	defer l.RUnlock()

	return l.errors
}
