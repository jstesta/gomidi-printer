package midiprinter

import "github.com/jstesta/gomidi/midi"

var META_STATUS_STRING map[byte]string = map[byte]string{
	midi.META_SEQUENCE_NUMBER:     "Sequence Number",
	midi.META_TEXT_EVENT:          "Text Event",
	midi.META_COPYRIGHT_NOTICE:    "Copyright Notice",
	midi.META_SEQUENCE_NAME:       "Sequence Name",
	midi.META_INSTRUMENT_NAME:     "Instrument Name",
	midi.META_LYRIC:               "Lyric",
	midi.META_MARKER:              "Marker",
	midi.META_CUE_POINT:           "Cue Point",
	midi.META_PROGRAM_NAME:        "Program Name",
	midi.META_DEVICE_NAME:         "Device Name",
	midi.META_MIDI_CHANNEL_PREFIX: "MIDI Channel Prefix",
	midi.META_MIDI_PORT:           "MIDI Port",
	midi.META_END_OF_TRACK:        "End Of Track",
	midi.META_SET_TEMPO:           "Set Tempo",
	midi.META_SMPTE_OFFSET:        "SMPTE Offset",
	midi.META_TIME_SIGNATURE:      "Time Signature",
	midi.META_KEY_SIGNATURE:       "Key Signature",
	midi.META_SEQUENCER_SPECIFIC:  "Sequencer Specific",
}

var MIDI_STATUS_STRING map[byte]string = map[byte]string{
	midi.MIDI_NOTE_OFF:                "Note OFF",
	midi.MIDI_NOTE_ON:                 "Note ON",
	midi.MIDI_POLYPHONIC_KEY_PRESSURE: "Polyphonic Key Pressure",
	midi.MIDI_CONTROL_CHANGE:          "Control Change",
	midi.MIDI_PROGRAM_CHANGE:          "Program Change",
	midi.MIDI_CHANNEL_PRESSURE:        "Channel Pressure",
	midi.MIDI_PITCH_WHEEL_CHANGE:      "Pitch Wheel Change",
}

var SYSEX_STATUS_STRING map[byte]string = map[byte]string{
	midi.SYSEX_SYSTEM_EXCLUSIVE:      "System Exclusive",
	midi.SYSEX_UNDEFINED_1:           "Undefined (1)",
	midi.SYSEX_SONG_POSITION_POINTER: "Song Position Pointer",
	midi.SYSEX_SONG_SELECT:           "Song Select",
	midi.SYSEX_UNDEFINED_4:           "Undefined (4)",
	midi.SYSEX_UNDEFINED_5:           "Undefined (5)",
	midi.SYSEX_TUNE_REQUEST:          "Tune Request",
	midi.SYSEX_END_OF_EXCLUSIVE:      "End Of Exclusive",
	midi.SYSEX_TIMING_CLOCK:          "Timing Clock",
	midi.SYSEX_UNDEFINED_9:           "Undefined (9)",
	midi.SYSEX_START:                 "Start",
	midi.SYSEX_CONTINUE:              "Continue",
	midi.SYSEX_STOP:                  "Stop",
	midi.SYSEX_UNDEFINED_13:          "Undefined (13)",
	midi.SYSEX_ACTIVE_SENSING:        "Active Sensing",
	midi.SYSEX_RESET:                 "Reset",
}
