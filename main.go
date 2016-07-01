package main

import (
	"flag"
	"log"
	"os"

	"fmt"
	"github.com/jstesta/gomidi"
	"github.com/jstesta/gomidi/cfg"
	"strconv"
	"strings"
)

const (
	column       = "|"
	plane        = "-"
	intersection = "+"
)

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
	//logger := log.New(f, "gomidi ", log.LUTC|log.LstdFlags)
	logger := log.New(os.Stdout, "gomidi ", log.LUTC|log.LstdFlags)
	logger.Printf("midiFile: %v", *midiFile)

	m, err := gomidi.ReadMidiFromFile(*midiFile, cfg.GomidiConfig{
		Log: logger,
	})
	if err != nil {
		logger.Fatal(err)
	}
	//logger.Printf("midi: %v", m)

	colWidths := [3]int{15, 15, 30}

	var spacerRow = buildSpacerRow(colWidths[:], leftPad, rightPad, plane, intersection)
	fmt.Println(spacerRow)

	var itemRow = buildItemRow(colWidths[:], leftPad, rightPad, column, m.Division().Type(), m.NumberOfTracks(), "")
	fmt.Println(itemRow)
}

func buildSpacerRow(colWidths []int, leftPad string, rightPad string, plane string, intersection string) string {

	var s string = intersection
	for _, w := range colWidths {
		s += strings.Repeat(plane, w+len(leftPad)+len(rightPad))
		s += intersection
	}
	return s
}

func buildItemFormatString(colWidths []int, leftPad string, rightPad string, column string) string {

	var s = column
	for _, w := range colWidths {
		s += leftPad + "%" + strconv.Itoa(w) + "v" + rightPad + column
	}
	return s
}

func buildItemRow(colWidths []int, leftPad string, rightPad string, column string, a ...interface{}) string {

	var f = buildItemFormatString(colWidths, leftPad, rightPad, column)

	for _, i := range a {
		switch i.(type) {
		case string:
			fmt.Println(i)
		default:
			fmt.Println("some", i)
		}
	}

	return fmt.Sprintf(f, a...)
}
