package midiprinter

import (
	"errors"
	"fmt"
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

func BuildItemRow(cfg *PrinterConfig, a ...interface{}) (string, error) {

	return BuildItemRowJustified(cfg, "", a...)
}

func BuildItemRowJustified(cfg *PrinterConfig, justify string, a ...interface{}) (string, error) {

	return buildItemRowExtended(cfg, false, justify, a...)
}

func buildItemRowExtended(cfg *PrinterConfig, extended bool, justify string, a ...interface{}) (string, error) {

	// todo this is inefficient, move out somewhere
	var f = buildItemFormatString(cfg.colWidths, cfg.leftPad, cfg.rightPad, cfg.column, justify)

	var q = make([]interface{}, len(cfg.colWidths))
	var toExtend = false

	if len(cfg.colWidths) < len(a) {
		return "", errors.New("more data columns than columns in configuration")
	}

	if len(cfg.colWidths) > len(a) {
		for i := 0; i < len(cfg.colWidths) - len(a); i++ {
			a = append(a, "")
		}
	}

	for idx, i := range a {
		var s = fmt.Sprintf("%v", i)

		if len(s) > cfg.colWidths[idx] {
			a[idx] = s[:cfg.colWidths[idx]]
			q[idx] = s[cfg.colWidths[idx]:]
			toExtend = true
		} else {
			q[idx] = ""
		}
	}

	if toExtend {
		ext, err := buildItemRowExtended(cfg, true, "-", q...)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(f, a...) + "\n" + ext, nil
	} else {
		return fmt.Sprintf(f, a...), nil
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
