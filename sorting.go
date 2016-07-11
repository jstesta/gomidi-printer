package midiprinter

import "github.com/jstesta/gomidi/midi"

// ByDeltaTime implements sort.Interface for []midi.Event based
// on the DeltaTime field
type ByDeltaTime []midi.Event

func (t ByDeltaTime) Len() int           { return len(t) }
func (t ByDeltaTime) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByDeltaTime) Less(i, j int) bool { return t[i].DeltaTime() < t[j].DeltaTime() }
