package midiprinter

import (
	"fmt"
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

	return buildItemRowExtended(cfg, false, a...)
}

func buildItemRowExtended(cfg *PrinterConfig, extended bool, a ...interface{}) string {

	// todo this is inefficient, move out somewhere
	var f = buildItemFormatString(cfg.colWidths, cfg.leftPad, cfg.rightPad, cfg.column, "")

	var q = make([]interface{}, len(cfg.colWidths))
	var toExtend = false

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

	if extended {
		f = buildItemFormatString(cfg.colWidths, cfg.leftPad, cfg.rightPad, cfg.column, "-")
	}

	if toExtend {
		return fmt.Sprintf(f, a...) + "\n" +
			buildItemRowExtended(cfg, true, q...)
	} else {
		return fmt.Sprintf(f, a...)
	}
}
