package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"sam/midiprinter"

	"github.com/jstesta/gomidi"
	"github.com/jstesta/gomidi/cfg"
	"github.com/jstesta/gomidi/midi"
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
	colWidths := []int{30, 30, 30, 30, 15}
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
			"Note",
			"Status")
		logr.Println(itemRowHeader)
		logr.Println(spacerRow)

		events := t.Events()
		sort.Sort(midiprinter.ByEvent(events))

		for _, e := range events {

			dataFormat := "%X"
			eventType := "default"
			note := ""
			var status byte
			switch t := e.(type) {
			case *midi.MidiEvent:
				eventType = "Midi Event"
				note = midiTypeNote(t.Status())
				status = t.Status()
			case *midi.SysexEvent:
				eventType = "Sysex Event"
				note = sysexTypeNote(t.Status())
				status = t.Status()
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
				status = t.MetaType()
			}

			itemRow := midiprinter.BuildItemRow(cfg,
				e.DeltaTime(),
				eventType,
				fmt.Sprintf(dataFormat, e.Data()),
				note,
				fmt.Sprintf("%X", status))
			logr.Println(itemRow)
		}
		logr.Println(spacerRow)
	}

	return nil
}

func metaTypeNote(b byte) string {

	if elem, ok := midiprinter.META_STATUS_STRING[b]; ok {
		return elem
	}
	return "default"
}

func midiTypeNote(b byte) string {

	if elem, ok := midiprinter.MIDI_STATUS_STRING[(b&0xF0)>>4]; ok {
		return elem
	}
	return "default"
}

func sysexTypeNote(b byte) string {

	if elem, ok := midiprinter.SYSEX_STATUS_STRING[b]; ok {
		return elem
	}
	return "default"
}
