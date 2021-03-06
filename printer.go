package midiprinter

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func BuildSpacerRow(cfg *PrinterConfig) string {

	var s string = cfg.intersection
	for _, w := range cfg.colWidths {
		s += strings.Repeat(cfg.plane, w+len(cfg.leftPad)+len(cfg.rightPad))
		s += cfg.intersection
	}
	return s
}

func BuildItemRow(cfg *PrinterConfig, a ...interface{}) string {

	return BuildItemRowJustified(cfg, "", a...)
}

func BuildItemRowJustified(cfg *PrinterConfig, justify string, a ...interface{}) string {

	return buildItemRowExtended(cfg, false, justify, a...)
}

func buildItemRowExtended(cfg *PrinterConfig, extended bool, justify string, a ...interface{}) string {

	// todo this is inefficient, move out somewhere
	var f = buildItemFormatString(cfg.colWidths, cfg.leftPad, cfg.rightPad, cfg.column, justify)

	var q = make([]interface{}, len(cfg.colWidths))
	var toExtend = false

	// it is an error if there is more data than columns
	if len(cfg.colWidths) < len(a) {
		log.Fatal("more data columns than columns in configuration")
	}

	// add empty data if there are more columns than data
	if len(cfg.colWidths) > len(a) {
		rem := len(cfg.colWidths) - len(a)
		for i := 0; i < rem; i++ {
			a = append(a, "")
		}
	}

	// check string length.  if it doesn't fit in the column, extend the
	// value to a new row
	for idx, i := range a {
		var s = fmt.Sprintf("%v", i)

		if nlIdx := strings.Index(s, "\n"); nlIdx != -1 {
			if nlIdx < cfg.colWidths[idx] {
				a[idx] = s[:nlIdx]
				q[idx] = s[nlIdx+1:]
				toExtend = true
			} else if len(s) > cfg.colWidths[idx] {
				a[idx] = s[:cfg.colWidths[idx]]
				q[idx] = s[cfg.colWidths[idx]:]
				toExtend = true
			}
		} else if len(s) > cfg.colWidths[idx] {
			a[idx] = s[:cfg.colWidths[idx]]
			q[idx] = s[cfg.colWidths[idx]:]
			toExtend = true
		} else {
			q[idx] = ""
		}
	}

	if toExtend {
		return fmt.Sprintf(f, a...) + "\n" + buildItemRowExtended(cfg, true, "", q...)
	} else {
		return fmt.Sprintf(f, a...)
	}
}

func buildItemFormatString(colWidths []int, leftPad string, rightPad string, column string, justify string) string {

	var s = column
	for _, w := range colWidths {
		s += leftPad + "%" + justify + strconv.Itoa(w) + "v" + rightPad + column
	}
	return s
}

func BuildHeader(cfg *PrinterConfig, text string) string {
	var s string
	s += BuildHeaderItemRow(cfg, text)
	return s
}

func BuildHeaderSpacerRow(cfg *PrinterConfig) string {
	var s string
	s += cfg.intersection
	for _, w := range cfg.colWidths {
		s += strings.Repeat(cfg.plane, w+len(cfg.leftPad)+len(cfg.rightPad))
	}
	s += strings.Repeat(cfg.plane, (len(cfg.colWidths)-1)*len(cfg.column))
	s += cfg.intersection
	return s
}

func BuildHeaderItemRow(cfg *PrinterConfig, text string) string {

	var n int
	for _, w := range cfg.colWidths {
		n += w
	}
	n += (len(cfg.colWidths) - 1) * (len(cfg.leftPad) + len(cfg.rightPad))
	n += (len(cfg.colWidths) - 1) * len(cfg.column)

	var textPart = text
	var remaining string
	if len(text) > n {
		textPart = text[:n]
		remaining = text[n:]
	}

	var s string = cfg.column
	s += cfg.leftPad
	s += fmt.Sprintf("%-"+strconv.Itoa(n)+"v", textPart)
	s += cfg.rightPad
	s += cfg.column

	if len(remaining) > 0 {
		s += "\n" + BuildHeaderItemRow(cfg, remaining)
	}

	return s
}
