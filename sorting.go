package midiprinter

import "github.com/jstesta/gomidi/midi"

type ByEvent []midi.Event

func (t ByEvent) Len() int      { return len(t) }
func (t ByEvent) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ByEvent) Less(i, j int) bool {

	iIsMeta := false
	jIsMeta := false

	switch iType := t[i].(type) {
	case *midi.MetaEvent:
		if iType.MetaType() == midi.META_END_OF_TRACK {
			return false
		}
		iIsMeta = true
	}

	switch jType := t[j].(type) {
	case *midi.MetaEvent:
		if jType.MetaType() == midi.META_END_OF_TRACK {
			return true
		}
		jIsMeta = true
	}

	if (t[i].DeltaTime() == t[j].DeltaTime()) {
		if (iIsMeta && !jIsMeta) {
			return true
		}
	}

	return t[i].DeltaTime() < t[j].DeltaTime()
}
