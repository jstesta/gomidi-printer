package midiprinter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jstesta/gomidi/midi"
	"math"
)

func ParseMetaEventData(e *midi.MetaEvent) string {

	data := e.Data()

	switch e.MetaType() {
	case midi.META_TEXT_EVENT,
		midi.META_COPYRIGHT_NOTICE,
		midi.META_SEQUENCE_NAME,
		midi.META_INSTRUMENT_NAME,
		midi.META_LYRIC,
		midi.META_MARKER,
		midi.META_CUE_POINT:
		return fmt.Sprintf("%s", data)

	case midi.META_SEQUENCE_NUMBER:
	case midi.META_PROGRAM_NAME:
	case midi.META_DEVICE_NAME:
	case midi.META_MIDI_CHANNEL_PREFIX:
	case midi.META_MIDI_PORT:
		return fmt.Sprintf("Midi Port: %v", data[0])
	case midi.META_END_OF_TRACK:
	case midi.META_SET_TEMPO:
		// 3 bytes representing 24-bit tempo
		var tempo int32
		tempo = int32(data[0])<<16 | int32(data[1])<<8 | int32(data[2])
		return fmt.Sprintf("Tempo: %v (msec/qtr-note)",
			tempo)
	case midi.META_SMPTE_OFFSET:
	case midi.META_TIME_SIGNATURE:
		// 4 bytes representing time signature
		return fmt.Sprintf("Time Signature: %v/%v\nClocks per qtr-note: %v\n32nd-notes per qtr-note: %v",
			data[0],
			math.Pow(2, float64(data[1])),
			data[2],
			data[3])

	case midi.META_KEY_SIGNATURE:
		// 2 bytes representing key signature
		sf := data[0]
		var sfi int32
		if sf == 0 {
			sfi = 0
		} else if sf>>7 == 1 {
			sfi = int32(sf) - 256
		} else {
			sfi = 256 - int32(sf)
		}

		var signature string
		if sfi < 0 {
			signature = fmt.Sprintf("%v flat(s)", sfi)
		} else if sfi > 0 {
			signature = fmt.Sprintf("%v sharp(s)", sfi)
		} else {
			signature = "key of C"
		}

		var key string
		if data[0] == 0 {
			key = "major"
		} else {
			key = "minor"
		}

		return fmt.Sprintf("%s, %s key",
			signature,
			key)

	case midi.META_SEQUENCER_SPECIFIC:
	}
	return fmt.Sprintf("%X", data)
}

func ToBitString(t interface{}) (s string) {

	switch t := t.(type) {

	case byte:
		for i := 0; i < 8; i++ {
			s = fmt.Sprintf("%d", t&1) + s
			t >>= 1
		}
		return

	case int32, uint32:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.LittleEndian, t)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		for _, b := range buf.Bytes() {
			s = ToBitString(b) + s
		}
		return

	case int:
		return ToBitString(int32(t))

	case uint:
		return ToBitString(uint32(t))

	default:
		return ""
	}
}