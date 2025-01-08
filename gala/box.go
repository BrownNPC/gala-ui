package gala

import (
	"image/color"
)

func newBox() Box {
	return Box{}
}

type Box struct {
	inUse bool
	id    string

	x, y int16

	parent *Box

	children []*Box

	onHover func(box *Box)
	baseStyle
}

// Appends a child to the parent's list of children
func (b *Box) Contains(children ...*Box) *Box {
	for _, child := range children {
		child.parent = b
		b.children = append(b.children, child)
	}
	return b
}
func (b *Box) reset() {

	b.inUse = false
	b.id = ""
	b.x = 0
	b.y = 0
	b.parent = nil
	b.children = b.children[:0]

	// reset baseStyle
	b.Size(0, 0)
	b.flex = 0
	b.gap = 0
	b.zindex = 0
	b.left = 0
	b.right = 0
	b.top = 0
	b.bottom = 0
	b.alignBits = 0
	b.Padding(0).
		Margin(0).
		Position_Relative().
		Display_Flex().
		JustifyContent_FlexStart().
		AlignItems_Stretch()

}

// apply padding to all sides
func (b *Box) Padding(i int16) *Box {
	return b.PaddingBottom(i).
		PaddingTop(i).
		PaddingLeft(i).
		PaddingRight(i)
}

// apply padding to the left side
func (b *Box) PaddingLeft(i int16) *Box {
	b.padding.left = i
	return b
}

// apply padding to the right side
func (b *Box) PaddingRight(i int16) *Box {
	b.padding.right = i
	return b
}

// apply padding to the top side
func (b *Box) PaddingTop(i int16) *Box {
	b.padding.top = i
	return b
}

// apply padding to the bottom side
func (b *Box) PaddingBottom(i int16) *Box {
	b.padding.bottom = i
	return b
}

// apply margin to all sides
func (b *Box) Margin(i int16) *Box {

	return b.MarginBottom(i).
		MarginTop(i).
		MarginLeft(i).
		MarginRight(i)
}

// apply margin to the left side
func (b *Box) MarginLeft(i int16) *Box {
	b.margin.left = i
	return b
}

// apply margin to the right side
func (b *Box) MarginRight(i int16) *Box {
	b.margin.right = i
	return b
}

// apply margin to the top side
func (b *Box) MarginTop(i int16) *Box {
	b.margin.top = i
	return b
}

// apply margin to the bottom side
func (b *Box) MarginBottom(i int16) *Box {
	b.margin.bottom = i
	return b
}

// Position

func (b *Box) Position_Relative() *Box {
	b.position = positionRelative
	return b
}
func (b *Box) Position_Absolute() *Box {
	b.position = positionAbsolute
	return b
}

// Display

func (b *Box) Display_Flex() *Box {
	b.display = displayFlex
	return b
}
func (b *Box) Display_None() *Box {
	b.display = displayNone
	return b
}

//Align Self

func (b *Box) AlignSelf_FlexStart() *Box {
	b.alignSelf(alignFlexStart)
	return b
}

func (b *Box) AlignSelf_Center() *Box {
	b.alignSelf(alignCenter)
	return b
}

func (b *Box) AlignSelf_FlexEnd() *Box {
	b.alignSelf(alignFlexEnd)
	return b
}

func (b *Box) AlignSelf_Stretch() *Box {
	b.alignSelf(alignStretch)
	return b
}

// Align Items

func (b *Box) AlignItems_FlexStart() *Box {
	b.alignSelf(alignFlexStart)
	return b
}

func (b *Box) AlignItems_Center() *Box {
	b.alignSelf(alignCenter)
	return b
}

func (b *Box) AlignItems_FlexEnd() *Box {
	b.alignSelf(alignFlexEnd)
	return b
}

func (b *Box) AlignItems_Stretch() *Box {
	b.alignSelf(alignStretch)
	return b
}

// Justify Content

func (b *Box) JustifyContent_FlexStart() *Box {
	b.justifyContent = justifyFlexStart
	return b
}

func (b *Box) JustifyContent_Center() *Box {
	b.justifyContent = justifyCenter
	return b
}

func (b *Box) JustifyContent_FlexEnd() *Box {
	b.justifyContent = justifyFlexEnd
	return b
}

func (b *Box) JustifyContent_SpaceBetween() *Box {
	b.justifyContent = justifySpaceBetween
	return b
}

func (b *Box) JustifyContent_SpaceAround() *Box {
	b.justifyContent = justifySpaceAround
	return b
}

func (b *Box) JustifyContent_SpaceEvenly() *Box {
	b.justifyContent = justifySpaceEvenly
	return b
}

/*
options:

	FlexStart; Center; FlexEnd; Stretch
*/
func (b *Box) alignSelf(align flexAlign) *Box {

	var k flexAlignKey
	switch align {

	case alignCenter:
		k = selfAlignCenter

	case alignFlexStart:
		k = selfAlignFlexStart
	case alignFlexEnd:
		k = selfAlignFlexEnd

	case alignStretch:
		k = selfAlignStretch
	}
	//unset any previous
	b.alignBits.
		unset(selfAlignCenter).
		unset(selfAlignFlexEnd).
		unset(selfAlignFlexStart).
		unset(selfAlignStretch)
	b.alignBits.set(k)
	return b

}

/*
options:

	FlexStart; Center; FlexEnd; Stretch
*/
func (b *Box) alignItems(align flexAlign) *Box {

	var k flexAlignKey
	switch align {

	case alignCenter:
		k = itemsAlignCenter
	case alignFlexStart:
		k = itemsAlignFlexStart
	case alignFlexEnd:
		k = itemsAlignFlexEnd
	case alignStretch:
		k = itemsAlignStretch
	}
	//unset any previous
	b.alignBits.
		unset(itemsAlignCenter).
		unset(itemsAlignFlexEnd).
		unset(itemsAlignFlexStart).
		unset(itemsAlignStretch)
	b.alignBits.set(k)
	return b
}

func (b *Box) FlexDirection_Row() *Box {
	b.flexDirection = directionRow
	return b
}
func (b *Box) FlexDirection_Column() *Box {
	b.flexDirection = directionColumn
	return b
}

/*
If the Box is positioned absolutely, those properties are used to define its position relative to the parent view.

If the position is relative, an offset is calculated from the resolved layout position and then applied to the view, without affecting sibling Boxes.

If width is not defined and both left and right are, then the Box will stretch to fill the space between the two offsets. The same applies to height and top/bottom.
*/
func (b *Box) Left(i int16) *Box {
	b.left = i
	return b
}

/*
If the Box is positioned absolutely, those properties are used to define its position relative to the parent view.

If the position is relative, an offset is calculated from the resolved layout position and then applied to the view, without affecting sibling Boxes.

If width is not defined and both left and right are, then the Box will stretch to fill the space between the two offsets. The same applies to height and top/bottom.
*/
func (b *Box) Right(i int16) *Box {
	b.right = i
	return b
}

/*
If the Box is positioned absolutely, those properties are used to define its position relative to the parent view.

If the position is relative, an offset is calculated from the resolved layout position and then applied to the view, without affecting sibling Boxes.

If width is not defined and both left and right are, then the Box will stretch to fill the space between the two offsets. The same applies to height and top/bottom.
*/
func (b *Box) Top(i int16) *Box {
	b.top = i
	return b
}

/*
If the Box is positioned absolutely, those properties are used to define its position relative to the parent view.

If the position is relative, an offset is calculated from the resolved layout position and then applied to the view, without affecting sibling Boxes.

If width is not defined and both left and right are, then the Box will stretch to fill the space between the two offsets. The same applies to height and top/bottom.
*/
func (b *Box) Bottom(i int16) *Box {
	b.bottom = i
	return b
}

/*
negative numbers between 0 and 1 are used for percentage

	example: -0.5 = 50%
*/
func (b *Box) Width(i float32) *Box {
	b.width = max(0, i)
	return b
}

/*
negative numbers between 0 and 1 are used for percentage

	example: -0.5 = 50%
	use gala.Percent() as a helper function
*/
func (b *Box) Height(i float32) *Box {
	b.height = max(0, i)
	return b
}

/*
negative numbers between 0 and 1 are used for percentage

	example: -0.5 = 50%
	use gala.Percent() as a helper function
*/
func (b *Box) Size(w, h float32) *Box {
	b.width, b.height = w, h
	return b
}

func (b *Box) BackgroundColor(col color.RGBA) *Box {
	b.backgroundColor = col
	return b
}

func (b *Box) ZIndex(i int16) *Box {
	b.zindex = i
	return b
}

func (b *Box) Id(i string) *Box {
	b.id = i
	return b
}

func (b *Box) Flex(i int16) *Box {
	b.flex = i
	return b
}

func (b *Box) Hovered(onHover func(box *Box)) *Box {
	b.onHover = onHover
	return b
}

func (b *Box) pointIsInside(x, y int32) bool {
	return b.x <= int16(x) && b.x+int16(b.width) >= int16(x) &&
		b.y <= int16(y) && b.y+int16(b.height) >= int16(y)

}
