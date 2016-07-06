package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"sam/midiprinter"

	"github.com/jstesta/gomidi"
	"github.com/jstesta/gomidi/cfg"
	"fmt"
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

	f, err := os.Create("out.log")
	if err != nil {
		log.Fatal(err)
	}
	printerLogger := log.New(f, "", 0)
	//printerLogger := log.New(os.Stdout, "", 0)

	stdOutLogger := log.New(os.Stdout, "gomidi ", log.LUTC|log.LstdFlags)
	stdOutLogger.Printf("midiFile: %v", *midiFile)

	m, err := gomidi.ReadMidiFromFile(*midiFile, cfg.GomidiConfig{
		Log: stdOutLogger,
	})
	if err != nil {
		stdOutLogger.Fatal(err)
	}
	//stdOutLogger.Printf("midi: %v", m)

	colWidths := [3]int{30, 30, 30}
	cfg := midiprinter.NewPrinterConfig(colWidths[:], leftPad, rightPad, plane, intersection, column)

	var headerSpacerRow = midiprinter.BuildHeaderSpacerRow(cfg)
	printerLogger.Println(headerSpacerRow)

	var header = midiprinter.BuildHeaderItemRow(cfg, "HEADER")
	printerLogger.Println(header)

	var spacerRow = midiprinter.BuildSpacerRow(cfg)
	printerLogger.Println(spacerRow)

	var itemRowHeader = midiprinter.BuildItemRowJustified(cfg,
		"-",
		"Division",
		"Number of Tracks",
		"Note")
	printerLogger.Println(itemRowHeader)
	printerLogger.Println(spacerRow)

	var itemRow = midiprinter.BuildItemRow(cfg,
		m.Division().Type(),
		strconv.Itoa(m.NumberOfTracks()),
		"")
	printerLogger.Println(itemRow)
	printerLogger.Println(spacerRow)

	for n, t := range m.Tracks() {
		header = midiprinter.BuildHeader(
			cfg,
			fmt.Sprintf("TRACK #"+ strconv.Itoa(n+1)))
		printerLogger.Println(header)
		printerLogger.Println(spacerRow)

		itemRowHeader =  midiprinter.BuildItemRowJustified(cfg,
			"-",
			"Delta Time",
			"Length (bytes)",
			"Data")
		printerLogger.Println(itemRowHeader)
		printerLogger.Println(spacerRow)

		for _, e := range t.Events() {
			itemRow = midiprinter.BuildItemRow(cfg,
				e.DeltaTime(),
				e.Length(),
				fmt.Sprintf("%X", e.Data()))
			printerLogger.Println(itemRow)
		}
		printerLogger.Println(spacerRow)
	}
}
