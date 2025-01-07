package gala

import (
	"image/color"
)

type flexDirection int8

const (
	directionRow flexDirection = iota
	directionColumn
)

type justifyContent int8

const (
	justifyFlexStart justifyContent = iota
	justifyCenter
	justifyFlexEnd
	justifySpaceBetween
	justifySpaceAround
	justifySpaceEvenly
)

type position int8

const (
	positionAbsolute position = iota
	positionRelative
)

// represent alignmet with a single 8 bit int
type alignProperties uint8
type flexAlignKey uint8

const (
	itemsAlignFlexStart = 1 << iota
	itemsAlignCenter
	itemsAlignFlexEnd
	itemsAlignStretch

	selfAlignFlexStart
	selfAlignCenter
	selfAlignFlexEnd
	selfAlignStretch
)

// not stored in the baseStyle, but its used by the end users
// for setting alignment
type flexAlign int8

const (
	alignFlexStart flexAlign = iota
	alignCenter
	alignFlexEnd
	alignStretch
)

type display int8

const (
	displayFlex display = iota
	displayNone
)

type padding struct {
	all, horizontal, vertical int16
	left, right, top, bottom  int16
}

type margin struct {
	all, horizontal, vertical int16
	left, right, top, bottom  int16
}

type baseStyle struct {
	width, height            float32 // 0.n floats represent percentage
	flex, gap, zindex        int16
	left, right, top, bottom int16
	padding                  padding
	margin                   margin

	// backgroundColor uint32 // Packed RGBA as uint32

	position        position
	display         display
	flexDirection   flexDirection
	justifyContent  justifyContent
	alignBits       alignProperties // Combined alignItems and alignSelf
	backgroundColor color.RGBA
}

// set sets a property.
func (p alignProperties) set(k flexAlignKey) alignProperties {
	return p | alignProperties(k)

}

// has checks if a property is set.
func (p alignProperties) has(k flexAlignKey) bool {
	return p&alignProperties(k) != 0
}

func (p alignProperties) unset(k flexAlignKey) alignProperties {
	return p &^ alignProperties(k) // Use bitwise AND NOT to unset the property
}
