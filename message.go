package main

import (

)
// Message contains the value of the previous event in the log, and the index of
// the previous index.
type Message struct {
	SourceID int
	Source      string
	Destination string
	Index   int // Size of log array
	Epoch   int
	Vote	bool
}
