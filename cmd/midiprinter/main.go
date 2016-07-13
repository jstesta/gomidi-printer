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
)

var colWidths = []int{10, 25, 20, 100}
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

	cfg := midiprinter.NewPrinterConfig(colWidths[:], leftPad, rightPad, plane, intersection, column)

	var headerSpacerRow = midiprinter.BuildHeaderSpacerRow(cfg)
	logr.Println(headerSpacerRow)

	var header = midiprinter.BuildHeaderItemRow(cfg, "HEADER")
	logr.Println(header)

	var spacerRow = midiprinter.BuildSpacerRow(cfg)
	logr.Println(spacerRow)

	itemRowHeader := midiprinter.BuildItemRowJustified(cfg,
		"-",
		"Format",
		"Division",
		"Tracks")
	logr.Println(itemRowHeader)
	logr.Println(spacerRow)

	var divisionNote string
	switch division := m.Division().(type) {
	default:
		divisionNote = "default"
	case *midi.MetricalDivision:
		divisionNote = fmt.Sprintf("Metrical Division [%v]",
			division.Resolution())
	case *midi.TimeCodeBasedDivision:
		divisionNote = fmt.Sprintf("Time-code Based [Format: %v, Resolution: %v]",
			division.Format(),
			division.Resolution())
	}

	itemRow := midiprinter.BuildItemRow(cfg,
		m.Header().Format(),
		divisionNote,
		strconv.Itoa(m.NumberOfTracks()))
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
			"Description",
			"Data")
		logr.Println(itemRowHeader)
		logr.Println(spacerRow)

		events := t.Events()

		for _, e := range events {

			data := ""
			eventType := "default"
			note := ""
			switch t := e.(type) {
			case *midi.MidiEvent:
				eventType = "Midi Event"
				note = midiTypeNote(t.Status())
				data = midiprinter.ParseMidiEventData(t)
			case *midi.SysexEvent:
				eventType = "Sysex Event"
				note = sysexTypeNote(t.Status())
				data = fmt.Sprintf("%X", t.Data())
			case *midi.MetaEvent:
				eventType = "Meta Event"
				data = midiprinter.ParseMetaEventData(t)
				note = metaTypeNote(t.MetaType())
			}

			itemRow := midiprinter.BuildItemRow(cfg,
				e.DeltaTime(),
				eventType,
				note,
				data)
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
