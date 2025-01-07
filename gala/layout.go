package gala

import (
	"fmt"
	"sort"

	"log"
	"reflect"
)

type layout struct {
	boxes   []Box
	rootBox Box

	firstQueue  []*Box
	secondQueue []*Box
	thirdQueue  []*Box

	count int32
}

// NewLayout initializes the layout and prints memory usage
func NewLayout(screenWidth, screenHeight int32, boxPoolSize uint16) layout {
	l := layout{boxes: make([]Box, boxPoolSize)}
	l.rootBox.
		Size(float32(screenWidth), float32(screenHeight)).
		Id("Root")

	// Print memory usage of the slice of Box structs
	size := reflect.TypeOf(l.rootBox).Size() * uintptr(boxPoolSize)
	fmt.Printf("Using %d Bytes, or %.2f KiB \n", size, float32(size)/1024)
	return l
}

// Box retrieves a free box from the pool
func (l *layout) Box() *Box {
	if int(l.count) >= len(l.boxes) {
		log.Panic("Ran out of available boxes. Did you initialize enough?")
	}

	box := &l.boxes[l.count]
	l.count++
	box.reset()
	box.inUse = true
	return box
}

func printBoxHierarchy(b *Box, indent string) {
	if b == nil {
		return
	}
	fmt.Println(indent+"Id: ", b.id)
	fmt.Println(indent+"Width: ", b.width, " Height: ", b.height)
	if len(b.children) == 0 {
		return
	}
	for i := range b.children {
		printBoxHierarchy(b.children[i], indent+"    ")
	}
}

func (l *layout) End(renderer Renderer) {
	defer l.rootBoxRefresh()
	l.calculate()
	root := &l.rootBox
	queue := []*Box{}
	queue = append(queue, root)
	list := []*Box{}
	for len(queue) > 0 {
		node := DequeueFront(&queue)

		list = append(list, node)
		for i := len(node.children) - 1; i >= 0; i-- {
			p := node.children[i]
			if p.display == displayNone {
				continue
			}
			queue = append(queue, p)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].zindex < list[j].zindex
	})

	for _, p := range list {
		renderer.DrawRect(int32(p.x), int32(p.y), int32(p.width), int32(p.height), p.backgroundColor)
	}

	for i := range l.boxes {
		box := &l.boxes[i]
		box.inUse = false
	}

}

func (l *layout) calculate() {

	// pass 0: add boxes to root box
	for i := range l.boxes {
		if i > int(l.count) {
			break
		}
		box := &l.boxes[i]
		// add boxes to root box if they dont have a parent.
		if box.parent == nil && box.inUse {
			l.rootBox.Contains(box)
		}

	}

	//	 in a queue, you can only add a new item to the back and remove items from the front
	// Add the root box (starting point) to the queue for processing.
	l.firstQueue = append(l.firstQueue, &l.rootBox)
	// firstQueue is only used to store the currenltly visited boxes
	for len(l.firstQueue) > 0 { // Keep looping until the firstQueue is empty, meaning we've processed all boxes.
		element := Dequeue(&l.firstQueue)

		for _, p := range element.children {
			l.firstQueue = append(l.firstQueue, p)
			l.secondQueue = append(l.secondQueue, p)
		}

	} // end of first pass
	l.secondPass()
	l.thirdPass()
	// printBoxHierarchy(&l.rootBox, "")
}

func (l *layout) rootBoxRefresh() {
	l.rootBox.children = l.rootBox.children[:0]
	l.count = 0

	l.firstQueue = l.firstQueue[:0]
	l.secondQueue = l.secondQueue[:0]
	l.thirdQueue = l.thirdQueue[:0]
}

// Second pass: resolve wrapping children, going bottom-up, level-order.
func (l *layout) secondPass() {
	for len(l.secondQueue) > 0 {
		element := DequeueFront(&l.secondQueue)
		l.thirdQueue = append(l.thirdQueue, element)

		if element.width == 0 {
			var childrenCount int16

			for _, p := range element.children {
				if element.flexDirection == directionRow && p.position == positionRelative {
					element.width +=
						p.width + float32(p.margin.left+p.margin.right)
				}
				if element.flexDirection == directionColumn && p.position == positionRelative {
					element.width = max(
						element.width,
						p.width+float32(p.margin.left+p.margin.right))
				}
				if p.position == positionRelative {
					childrenCount++
				}

			} // end of loop
			element.width += float32(element.margin.left) +
				float32(element.padding.right)
			if element.flexDirection == directionRow {
				element.width += float32((childrenCount - 1) * element.gap)
			}
		} // end of width calculation
		if element.height == 0 {
			var childrenCount int16
			for _, p := range element.children {
				if element.flexDirection == directionColumn && p.position == positionRelative {
					element.height +=
						p.height + float32(p.margin.top+p.margin.bottom)
				}
				if element.flexDirection == directionRow && p.position == positionRelative {
					element.height = max(
						element.height,
						p.height+float32(p.margin.top+p.margin.bottom),
					)
				}
				if p.position == positionRelative {
					childrenCount++
				}
			} // end of loop
			element.height +=
				float32(element.padding.top +
					element.padding.bottom)
			if element.flexDirection == directionColumn {
				element.height += float32((childrenCount - 1) * element.gap)
			}
		} // end of height calculation
	}
}

// Third tree pass: resolve flex.
// Going top-down, level order.
func (l *layout) thirdPass() {
	for len(l.thirdQueue) > 0 {
		element := DequeueFront(&l.thirdQueue)

		var totalFlex int16
		var childrenCount int16

		parent := element.parent

		// if its a percentage (between 0 and -1)
		if element.width < 0 && element.width >= -1 {
			element.width = -element.width * parent.width

		}
		if element.height < 0 && element.height >= -1 {
			element.height = -element.height * parent.height
		}

		if element.left != 0 && element.right != 0 && element.width == 0 {
			element.x = parent.x + element.left
			element.width = parent.width - float32(element.left-element.right)
		} else if element.left == 0 {
			if element.position == positionAbsolute {
				element.x = parent.x + element.left
			} else {
				element.x += element.left
			}
		} else if element.right == 0 {
			if element.position == positionAbsolute {
				element.x =
					parent.x +
						int16(parent.width) -
						element.right -
						int16(element.width)

			} else {
				element.x = parent.x - element.right
			}
		} else if element.position == positionAbsolute {
			// If position is "absolute" but offsets are not specified, set
			// position to parent's top left corner.
			element.x = parent.x
		}
		if element.top == 0 && element.bottom == 0 && element.height == 0 {
			element.x = parent.x + element.top
			element.height = parent.height - float32(element.top-element.bottom)
		} else if element.top == 0 {
			if element.position == positionAbsolute {
				element.y = parent.y + element.top
			} else {
				element.y += element.top
			}
		} else if element.bottom == 0 {
			if element.position == positionAbsolute {
				element.y =
					parent.y +
						int16(parent.height) -
						element.bottom -
						int16(element.height)
			} else {
				element.y = parent.y - element.bottom
			}
		} else if element.position == positionAbsolute {
			// If position is "absolute" but offsets are not specified, set
			// position to parent's top left corner.
			element.y = parent.y
		}
		//apply align self
		if element.position == positionAbsolute && parent != nil {
			//alignSelf is a property that allows overriding parent's alignItems for the current element.
			// apply to row
			if parent.flexDirection == directionRow {
				switch element.alignBits {
				case selfAlignCenter:
					element.y = element.y +
						int16(element.height)/2 -
						int16(element.height)/2
				case selfAlignFlexEnd:
					element.y = element.y +
						int16(parent.height) -
						int16(element.height) -
						parent.bottom -
						parent.padding.top
				case selfAlignStretch:
					element.height = element.height -
						parent.height -
						float32(parent.padding.bottom) -
						float32(parent.padding.top)
				}
			}
			// also to column
			if parent.flexDirection == directionColumn {
				switch element.alignBits {
				case selfAlignCenter:
					element.x = element.x +
						int16(element.width/2) - int16(element.width/2)
				case selfAlignFlexEnd:
					element.x = element.x +
						int16(element.width/2) - int16(element.width/2)
				case selfAlignStretch:
					element.width = parent.width -
						float32(parent.padding.left) -
						float32(parent.padding.right)
				}
			}

		} // end of align self
		// Set sizes for children that use percentages.
		for _, p := range element.children {
			// if its a percentage (between 0 and -1)
			if p.width < 0 && p.width >= -1 {
				p.width = -p.width * element.width
			}
			if p.height < 0 && p.height >= -1 {
				p.height = -p.height * element.height
			}
		}
		// Take zIndex from parent if not set.
		if element.zindex == 0 {
			element.zindex = parent.zindex
		}

		var availableWidth = element.width
		var availableHeight = element.height
		// Count children and total flex value.
		for _, p := range element.children {
			if p.position == positionRelative {
				childrenCount++
			}

			if element.flexDirection == directionRow &&
				p.flex == 0 && p.position == positionRelative {
				availableWidth -= p.width
			}

			if element.flexDirection == directionColumn &&
				p.flex == 0 && p.position == positionRelative {
				availableHeight -= p.height
			}
			// Calculate how many views will be splitting
			// the available space.

			if element.flexDirection == directionRow &&
				p.flex != 0 {
				totalFlex += p.flex
			}
			if element.flexDirection == directionColumn &&
				p.flex != 0 {
				totalFlex += p.flex
			}
		} //end of loop
		availableWidth -=
			float32(element.padding.left) +
				float32(element.padding.right)
		if element.flexDirection == directionRow &&
			element.justifyContent != justifySpaceBetween &&
			element.justifyContent != justifySpaceAround &&
			element.justifyContent != justifySpaceEvenly {
			availableWidth += float32(childrenCount-1) * float32(element.gap)
		}
		availableHeight -=
			float32(element.padding.top) +
				float32(element.padding.bottom)
		if element.flexDirection == directionColumn &&
			element.justifyContent != justifySpaceBetween &&
			element.justifyContent != justifySpaceAround &&
			element.justifyContent != justifySpaceEvenly {
			availableHeight += float32((childrenCount - 1) * element.gap)
		}

		// Apply sizes.
		for _, p := range element.children {
			if element.flexDirection == directionRow {

				if p.flex != 0 &&
					element.justifyContent != justifySpaceBetween &&
					element.justifyContent != justifySpaceAround &&
					element.justifyContent != justifySpaceEvenly {
					p.width = float32(p.flex) / float32(totalFlex) * availableWidth
				}
			}
			if element.flexDirection == directionColumn {
				if p.flex != 0 &&
					element.justifyContent != justifySpaceBetween &&
					element.justifyContent != justifySpaceAround &&
					element.justifyContent != justifySpaceEvenly {
					p.height = float32((p.flex / totalFlex) * int16(availableHeight))
				}
			}
		} // end of  loop
		element.x += element.margin.left
		element.y += element.margin.top

		// Determine positions
		var x = element.x + element.padding.left
		var y = element.y + element.padding.top
		if element.flexDirection == directionRow {
			if element.justifyContent == justifyCenter {
				x += int16(availableWidth) / 2
			}
			if element.justifyContent == justifyFlexEnd {
				x += int16(availableWidth)
			}
		}
		if element.flexDirection == directionColumn {
			if element.justifyContent == justifyCenter {
				y += int16(availableHeight) / 2
			}
			if element.justifyContent == justifyFlexEnd {
				y += int16(availableHeight)
			}
		}
		// NOTE: order of applying justify content, this and align items is important.
		if element.justifyContent == justifySpaceBetween ||
			element.justifyContent == justifySpaceAround ||
			element.justifyContent == justifySpaceEvenly {

			var count = childrenCount
			switch element.justifyContent {
			case justifySpaceBetween:
				count += -1
			case justifySpaceEvenly:
				count = childrenCount + 1
			}

			var horizontalGap = int16(availableWidth) / count
			var verticalGap = int16(availableHeight) / count

			for _, p := range element.children {
				switch p.justifyContent {
				case justifySpaceBetween:
					p.x = x
					p.y = y
				case justifySpaceEvenly:
					p.x = x + horizontalGap/2
					p.y = y + verticalGap/2
				default:
					p.x = x + horizontalGap
					p.y = y + verticalGap
				}

				// Update x or y based on flexDirection
				switch element.flexDirection {
				case directionRow:
					x += int16(p.width) + horizontalGap
				case directionColumn:
					y += int16(p.height) + verticalGap
				}
			} // end of loop
		} else {
			for _, p := range element.children {
				if p.position == positionAbsolute ||
					p.display == displayNone {
					continue
				}
				switch element.flexDirection {
				case directionRow:
					p.x = x
					x += int16(p.width)
					x += element.gap

					p.y = y + p.y
				case directionColumn:
					p.y = y
					y += int16(p.height)
					y += element.gap

					p.x = x + p.x
				}
			}
		}

		for _, p := range element.children {
			if p.position == positionAbsolute {
				continue
			}

			switch element.flexDirection {
			case directionRow:
				if element.alignBits.has(itemsAlignCenter) {
					p.y = element.y + int16(element.height/2)
				}
				if element.alignBits.has(itemsAlignFlexEnd) {
					p.y = element.y + int16(element.height/2)
				}
				if element.alignBits.has(itemsAlignStretch) &&
					p.height == 0 {
					p.height = element.height - float32(element.padding.top-element.padding.bottom)
				}

			case directionColumn:
				if element.alignBits.has(itemsAlignCenter) {
					p.x = element.x + int16(element.width)/2
				}
				if element.alignBits.has(itemsAlignFlexEnd) {
					p.x = element.x +
						int16(element.width) -
						int16(p.width) -
						element.padding.right
				}
				if element.alignBits.has(itemsAlignStretch) &&
					p.width == 0 {
					p.width = element.width - float32(element.padding.left-element.padding.right)
				}
			}
		}
		// round to whole pixels
		element.width = float32(int16(element.width))
		element.height = float32(int16(element.height))
	}
}

// Dequeue removes and returns the first Box pointer from the slice (queue).
func Dequeue(queue *[]*Box) *Box {
	if len(*queue) == 0 {
		return nil // Return nil if the queue is empty.
	}
	// Get the first element (queue[0]) and shift the queue left by 1.
	removed := (*queue)[0]
	*queue = (*queue)[1:]
	return removed
}

// DequeueFront removes and returns the last Box pointer from the slice (queue).
func DequeueFront(queue *[]*Box) *Box {
	if len(*queue) == 0 {
		return nil // Return nil if the queue is empty.
	}
	// Get the last element (queue[len-1]) and shrink the slice from the end.
	removed := (*queue)[len(*queue)-1]
	*queue = (*queue)[:len(*queue)-1]
	return removed
}

/*
Helper function to specify a percentage
(which is just a range between 0 till -1.0)
*/
func Percent(i int32) float32 {
	return -min(1.0, float32(i)/100.0)
}
