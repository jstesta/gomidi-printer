package midiprinter

import "github.com/jstesta/gomidi/midi"

type ByEvent []midi.Event

func (t ByEvent) Len() int      { return len(t) }
func (t ByEvent) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t ByEvent) Less(i, j int) bool {

	iIsMeta := false
	jIsMeta := false
	var iStatus byte = 0
	var jStatus byte = 0

	switch iType := t[i].(type) {
	case *midi.MetaEvent:
		// end of track meta event is always last ("greater" than
		// all other events)
		if iType.MetaType() == midi.META_END_OF_TRACK {
			return false
		}
		iIsMeta = true
		iStatus = iType.MetaType()
	}

	switch jType := t[j].(type) {
	case *midi.MetaEvent:
		// end of track meta event is always last ("greater" than
		// all other events)
		if jType.MetaType() == midi.META_END_OF_TRACK {
			return true
		}
		jIsMeta = true
		jStatus = jType.MetaType()
	}

	// for events with the same delta time
	if t[i].DeltaTime() == t[j].DeltaTime() {
		// meta events come first
		if iIsMeta && !jIsMeta {
			return true
		}

		// if both are meta events, order by status
		if iIsMeta && jIsMeta {
			return iStatus < jStatus
		}
	}

	// everything else is ordered by delta time
	return t[i].DeltaTime() < t[j].DeltaTime()
}
