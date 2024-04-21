package mock

import "sync"

type Logger struct {
	Infos  []string
	Errors []string
	sync.RWMutex
}

func (l *Logger) Info(i string, _ ...any) {
	l.Lock()
	defer l.Unlock()

	l.Infos = append(l.Infos, i)
}

func (l *Logger) Error(e string, _ ...any) {
	l.Lock()
	defer l.Unlock()

	l.Errors = append(l.Errors, e)
}

func (l *Logger) GotInfos() []string {
	l.RLock()
	defer l.RUnlock()

	return l.Infos
}

func (l *Logger) GotErrors() []string {
	l.RLock()
	defer l.RUnlock()

	return l.Errors
}
