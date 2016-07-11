package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"sam/midiprinter"

	"github.com/jstesta/gomidi"
	"github.com/jstesta/gomidi/cfg"
	"github.com/jstesta/gomidi/midi"
	"sort"
)

var column = "|"
var plane = "-"
var intersection = "+"
var leftPad = " "
var rightPad = " "

func main() {
	var (
		midiFile = flag.String("input", "", "Filesystem location of MIDI file to parse")
	)
	flag.Parse()

	stdOutLogger := log.New(os.Stdout, "gomidi ", log.LUTC|log.LstdFlags)
	stdOutLogger.Printf("midiFile: %v", *midiFile)

	m, err := gomidi.ReadMidiFromFile(*midiFile, cfg.GomidiConfig{
		Log: stdOutLogger,
	})
	if err != nil {
		stdOutLogger.Fatal(err)
	}

	f, err := os.Create("out.log")
	if err != nil {
		log.Fatal(err)
	}
	printerLogger := log.New(f, "", 0)
	//printerLogger := log.New(os.Stdout, "", 0)
	printMidi(m, printerLogger)
}

func printMidi(m *midi.Midi, logr *log.Logger) error {
	colWidths := [4]int{30, 30, 30, 30}
	cfg := midiprinter.NewPrinterConfig(colWidths[:], leftPad, rightPad, plane, intersection, column)

	var headerSpacerRow = midiprinter.BuildHeaderSpacerRow(cfg)
	logr.Println(headerSpacerRow)

	var header = midiprinter.BuildHeaderItemRow(cfg, "HEADER")
	logr.Println(header)

	var spacerRow = midiprinter.BuildSpacerRow(cfg)
	logr.Println(spacerRow)

	itemRowHeader := midiprinter.BuildItemRowJustified(cfg,
		"-",
		"Division",
		"Number of Tracks")
	logr.Println(itemRowHeader)
	logr.Println(spacerRow)

	itemRow := midiprinter.BuildItemRow(cfg,
		m.Division().Type(),
		strconv.Itoa(m.NumberOfTracks()),
		"")
	logr.Println(itemRow)
	logr.Println(spacerRow)

	for n, t := range m.Tracks() {
		header = midiprinter.BuildHeader(
			cfg,
			fmt.Sprintf("TRACK #"+strconv.Itoa(n+1)))
		logr.Println(header)
		logr.Println(spacerRow)

		itemRowHeader := midiprinter.BuildItemRowJustified(cfg,
			"-",
			"Delta Time",
			"Type",
			"Data",
			"Note")
		logr.Println(itemRowHeader)
		logr.Println(spacerRow)

		events := t.Events()
		sort.Sort(midiprinter.ByEvent(events))

		for _, e := range events {

			dataFormat := "%X"
			eventType := "default"
			note := ""
			switch t := e.(type) {
			case *midi.MidiEvent:
				eventType = "Midi Event"
				note = midiTypeNote(t.Status())
			case *midi.SysexEvent:
				eventType = "Sysex Event"
				note = sysexTypeNote(t.Status())
			case *midi.MetaEvent:
				eventType = "Meta Event"
				switch t.MetaType() {
				case midi.META_TEXT_EVENT,
					midi.META_COPYRIGHT_NOTICE,
					midi.META_SEQUENCE_NAME,
					midi.META_INSTRUMENT_NAME,
					midi.META_LYRIC,
					midi.META_MARKER,
					midi.META_CUE_POINT:
					dataFormat = "%s"
				}
				note = metaTypeNote(t.MetaType())
			}

			itemRow := midiprinter.BuildItemRow(cfg,
				e.DeltaTime(),
				eventType,
				fmt.Sprintf(dataFormat, e.Data()),
				note)
			logr.Println(itemRow)
		}
		logr.Println(spacerRow)
	}

	return nil
}

func metaTypeNote(b byte) string {
	switch b {
	default:
		return "default"
	case midi.META_SEQUENCE_NUMBER:
		return "Sequence Number"
	case midi.META_TEXT_EVENT:
		return "Text Event"
	case midi.META_COPYRIGHT_NOTICE:
		return "Copyright Notice"
	case midi.META_SEQUENCE_NAME:
		return "Sequence Name"
	case midi.META_INSTRUMENT_NAME:
		return "Instrument Name"
	case midi.META_LYRIC:
		return "Lyric"
	case midi.META_MARKER:
		return "Marker"
	case midi.META_CUE_POINT:
		return "Cue Point"
	case midi.META_MIDI_CHANNEL_PREFIX:
		return "MIDI Channel Prefix"
	case midi.META_END_OF_TRACK:
		return "End Of Track"
	case midi.META_SET_TEMPO:
		return "Set Tempo"
	case midi.META_SMPTE_OFFSET:
		return "SMPTE Offset"
	case midi.META_TIME_SIGNATURE:
		return "Time Signature"
	case midi.META_KEY_SIGNATURE:
		return "Key Signature"
	case midi.META_SEQUENCER_SPECIFIC:
		return "Sequencer Specific"
	}
}

func midiTypeNote(b byte) string {
	switch (b & 0xF0) >> 4 {
	default:
		return "default"
	case midi.MIDI_NOTE_OFF:
		return "Note OFF"
	case midi.MIDI_NOTE_ON:
		return "Note ON"
	case midi.MIDI_POLYPHONIC_KEY_PRESSURE:
		return "Polyphonic Key Pressure"
	case midi.MIDI_CONTROL_CHANGE:
		return "Control Change"
	case midi.MIDI_PROGRAM_CHANGE:
		return "Program Change"
	case midi.MIDI_CHANNEL_PRESSURE:
		return "Channel Pressure"
	case midi.MIDI_PITCH_WHEEL_CHANGE:
		return "Pitch Wheel Change"
	}
}

func sysexTypeNote(b byte) string {
	switch b {
	default:
		return "default"
	case midi.SYSEX_SYSTEM_EXCLUSIVE:
		return "System Exclusive"
	case midi.SYSEX_UNDEFINED_1:
		return "Undefined (1)"
	case midi.SYSEX_SONG_POSITION_POINTER:
		return "Song Position Pointer"
	case midi.SYSEX_SONG_SELECT:
		return "Song Select"
	case midi.SYSEX_UNDEFINED_4:
		return "Undefined (4)"
	case midi.SYSEX_UNDEFINED_5:
		return "Undefined (5)"
	case midi.SYSEX_TUNE_REQUEST:
		return "Tune Request"
	case midi.SYSEX_END_OF_EXCLUSIVE:
		return "End Of Exclusive"
	case midi.SYSEX_TIMING_CLOCK:
		return "Timing Clock"
	case midi.SYSEX_UNDEFINED_9:
		return "Undefined (9)"
	case midi.SYSEX_START:
		return "Start"
	case midi.SYSEX_CONTINUE:
		return "Continue"
	case midi.SYSEX_STOP:
		return "Stop"
	case midi.SYSEX_UNDEFINED_13:
		return "Undefined (13)"
	case midi.SYSEX_ACTIVE_SENSING:
		return "Active Sensing"
	case midi.SYSEX_RESET:
		return "Reset"
	}
}
