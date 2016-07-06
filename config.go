package midiprinter

import (
	"fmt"
	"strconv"
)

type PrinterConfig struct {
	colWidths          []int
	leftPad            string
	rightPad           string
	plane              string
	intersection       string
	column             string
	itemFormat         string
	extendedItemFormat string
}

func (c *PrinterConfig) String() string {
	return fmt.Sprintf("PrinterConfig [colWidths=%v, leftPad=%v, rightPad=%v, plane=%v, intersection=%v, column=%v]",
		c.colWidths,
		c.leftPad,
		c.rightPad,
		c.plane,
		c.intersection,
		c.column,
	)
}

// TODO make a builder
func NewPrinterConfig(colWidths []int, leftPad string, rightPad string, plane string, intersection string, column string) *PrinterConfig {
	return &PrinterConfig{
		colWidths,
		leftPad,
		rightPad,
		plane,
		intersection,
		column,
		buildItemFormatString(colWidths, leftPad, rightPad, column, ""),
		buildItemFormatString(colWidths, leftPad, rightPad, column, "-"),
	}
}

func buildItemFormatString(colWidths []int, leftPad string, rightPad string, column string, justify string) string {

	var s = column
	for _, w := range colWidths {
		s += leftPad + "%" + justify + strconv.Itoa(w) + "v" + rightPad + column
	}
	return s
}
