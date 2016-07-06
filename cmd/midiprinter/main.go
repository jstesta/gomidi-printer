package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"sam/midiprinter"

	"github.com/jstesta/gomidi"
	"github.com/jstesta/gomidi/cfg"
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

	//f, err := os.Create("out.log")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//printerLogger := log.New(f, "", 0)
	printerLogger := log.New(os.Stdout, "", 0)

	stdOutLogger := log.New(os.Stdout, "gomidi ", log.LUTC|log.LstdFlags)
	stdOutLogger.Printf("midiFile: %v", *midiFile)

	m, err := gomidi.ReadMidiFromFile(*midiFile, cfg.GomidiConfig{
		Log: stdOutLogger,
	})
	if err != nil {
		stdOutLogger.Fatal(err)
	}
	//stdOutLogger.Printf("midi: %v", m)

	colWidths := [3]int{15, 30, 30}
	cfg := midiprinter.NewPrinterConfig(colWidths[:], leftPad, rightPad, plane, intersection, column)

	var spacerRow = midiprinter.BuildSpacerRow(cfg)
	printerLogger.Println(spacerRow)

	var itemRow = midiprinter.BuildItemRow(cfg,
		m.Division().Type(),
		strconv.Itoa(m.NumberOfTracks())+" really really long tracks string",
		"this is a super duper really long string that probably can't fit in the row")
	printerLogger.Println(itemRow)
}
