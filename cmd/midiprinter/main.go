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

	itemRowHeader, err := midiprinter.BuildItemRowJustified(cfg,
		"-",
		"Division",
		"Number of Tracks",
		"Note")
	if err != nil {
		return err
	}
	logr.Println(itemRowHeader)
	logr.Println(spacerRow)

	itemRow, err := midiprinter.BuildItemRow(cfg,
		m.Division().Type(),
		strconv.Itoa(m.NumberOfTracks()),
		"")
	if err != nil {
		return err
	}
	logr.Println(itemRow)
	logr.Println(spacerRow)

	for n, t := range m.Tracks() {
		header = midiprinter.BuildHeader(
			cfg,
			fmt.Sprintf("TRACK #"+strconv.Itoa(n+1)))
		logr.Println(header)
		logr.Println(spacerRow)

		itemRowHeader, err := midiprinter.BuildItemRowJustified(cfg,
			"-",
			"Delta Time",
			"Length (bytes)",
			"Data")
		if err != nil {
			return err
		}
		logr.Println(itemRowHeader)
		logr.Println(spacerRow)

		events := t.Events()
		sort.Sort(midiprinter.ByDeltaTime(events))

		for _, e := range events {

			var dataFormat string
			switch e.(type) {
			default:
				dataFormat = "%X"
			case *midi.MetaEvent:
				dataFormat = "%s"
			}

			itemRow, err := midiprinter.BuildItemRow(cfg,
				e.DeltaTime(),
				e.Length(),
				fmt.Sprintf(dataFormat, e.Data()))
			if err != nil {
				return err
			}
			logr.Println(itemRow)
		}
		logr.Println(spacerRow)
	}

	return nil
}
