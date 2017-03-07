package main

import (

)
// Message contains the value of the previous event in the log, and the index of
// the previous index.
type Message struct {
	Source      int
	Destination int
	PrevIndex   int
	PrevEpoch   int
}
