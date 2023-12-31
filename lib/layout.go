package lib

import (
	"fmt"
	"sort"
)

type LayoutOverride struct {
	Width  int
	Height int
}

type Layout struct {
	Columns   []int
	Heights   []int
	Order     []int
	Overrides map[int]LayoutOverride
}

func extractWidths(row []Box) []int {
	widths := make([]int, 0)
	for _, box := range row {
		widths = append(widths, box.W)
	}
	return widths
}

func extractHeights(row []Box) []int {
	heights := make([]int, 0)
	for _, box := range row {
		heights = append(heights, box.H)
	}
	return heights
}

func isPrefix[T comparable](a, b []T) bool {
	if len(a) > len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func AnalyzeGrid(grid []Box) Layout {
	sort.Slice(grid, func(i, j int) bool {
		if grid[i].Y != grid[j].Y {
			return grid[i].Y < grid[j].Y
		}
		return grid[i].X < grid[j].X
	})

	stat := make(map[string]int, 0)
	rows := make([][]Box, 0)
	for _, element := range grid {
		if len(rows) == 0 || rows[len(rows)-1][0].Y != element.Y {
			rows = append(rows, make([]Box, 0))
		}
		rows[len(rows)-1] = append(rows[len(rows)-1], element)
	}
	columnsString, columns := "", make([]int, 0)
	for _, row := range rows {
		width := extractWidths(row)
		widthString := fmt.Sprintf("%v", width)
		if _, ok := stat[widthString]; !ok {
			stat[widthString] = 0
		}
		stat[widthString] += len(row)
		if stat[widthString] > stat[columnsString] {
			columnsString = widthString
			columns = width
		}
	}
	layout := Layout{Columns: columns, Heights: make([]int, 0), Order: make([]int, 0), Overrides: make(map[int]LayoutOverride)}
	for _, row := range rows {
		width := extractWidths(row)
		heights := Unique(extractHeights(row))
		if !isPrefix(width, columns) || len(heights) > 1 {
			for _, cell := range row {
				h := cell.H
				if heights[0] == cell.H {
					h = 0
				}
				layout.Overrides[cell.Id] = LayoutOverride{Width: cell.W, Height: h}
			}
		}
		layout.Heights = append(layout.Heights, heights[0])
		for _, cell := range row {
			layout.Order = append(layout.Order, cell.Id)
		}
	}
	if len(Unique(layout.Heights)) == 1 {
		layout.Heights = []int{layout.Heights[0]}
	}
	return layout
}
