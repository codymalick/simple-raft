package main

import ()

// Log records a request
type Log struct {
	Value     int
	Epoch     int
	Committed bool
}

// Commit sets a log to become committed
func (l *Log) Commit() {
	l.Committed = true
}
