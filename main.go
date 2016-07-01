package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jstesta/gomidi"
	"github.com/jstesta/gomidi/cfg"
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

	colWidths := [3]int{15, 30, 30}

	var spacerRow = buildSpacerRow(colWidths[:], leftPad, rightPad, plane, intersection)
	fmt.Println(spacerRow)

	var itemRow = buildItemRow(colWidths[:], leftPad, rightPad, column,
		m.Division().Type(),
		strconv.Itoa(m.NumberOfTracks())+" really really long tracks string",
		"this is a super duper really long string that probably can't fit in the row")
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

func buildItemFormatString(colWidths []int, leftPad string, rightPad string, column string, justify string) string {

	var s = column
	for _, w := range colWidths {
		s += leftPad + "%" + justify + strconv.Itoa(w) + "v" + rightPad + column
	}
	return s
}

func buildItemRow(colWidths []int, leftPad string, rightPad string, column string, a ...interface{}) string {
	return buildItemRowExtended(colWidths, leftPad, rightPad, column, false, a...)
}

func buildItemRowExtended(colWidths []int, leftPad string, rightPad string, column string, extended bool, a ...interface{}) string {

	// todo this is inefficient, move out somewhere
	var f = buildItemFormatString(colWidths, leftPad, rightPad, column, "")

	var q = make([]interface{}, len(colWidths))
	var toExtend = false

	for idx, i := range a {
		var s = fmt.Sprintf("%v", i)

		if len(s) > colWidths[idx] {
			a[idx] = s[:colWidths[idx]]
			q[idx] = s[colWidths[idx]:]
			toExtend = true
		} else {
			q[idx] = ""
		}
	}

	if extended {
		f = buildItemFormatString(colWidths, leftPad, rightPad, column, "-")
	}

	if toExtend {
		return fmt.Sprintf(f, a...) + "\n" +
			buildItemRowExtended(colWidths, leftPad, rightPad, column, true, q...)
	} else {
		return fmt.Sprintf(f, a...)
	}
}
