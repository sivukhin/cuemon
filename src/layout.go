package src

import (
	"sort"
	"strconv"
	"strings"
)

type Grid struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type LayoutOverride struct {
	Width, Height, X, Y int
}

type Layout struct {
	Columns   []int
	Heights   []int
	Order     []int
	Overrides map[int]LayoutOverride
}

type Indexed struct {
	Index int
	Grid  Grid
}

func createWidthDescriptor(row []Indexed) (string, []int) {
	widthString := make([]string, 0)
	width := make([]int, 0)
	for _, cell := range row {
		widthString = append(widthString, strconv.Itoa(cell.Grid.W))
		width = append(width, cell.Grid.W)
	}
	return strings.Join(widthString, ","), width
}

func getUnique[T comparable](a []T) []T {
	used := make(map[T]struct{}, 0)
	unique := make([]T, 0)
	for _, element := range a {
		if _, ok := used[element]; !ok {
			used[element] = struct{}{}
			unique = append(unique, element)
		}
	}
	return unique
}

func getHeights(row []Indexed) []int {
	heights := make([]int, 0)
	for _, cell := range row {
		heights = append(heights, cell.Grid.H)
	}
	return getUnique(heights)
}

func AnalyzeGrid(grid []Grid) Layout {
	indexed := make([]Indexed, 0, len(grid))
	for i, element := range grid {
		indexed = append(indexed, Indexed{Index: i, Grid: element})
	}
	sort.Slice(indexed, func(i, j int) bool {
		if indexed[i].Grid.Y != indexed[j].Grid.Y {
			return indexed[i].Grid.Y < indexed[j].Grid.Y
		}
		return indexed[i].Grid.X < indexed[j].Grid.X
	})

	stat := make(map[string]int, 0)
	rows := make([][]Indexed, 0)
	for _, element := range indexed {
		if len(rows) == 0 || rows[len(rows)-1][0].Grid.Y != element.Grid.Y {
			rows = append(rows, make([]Indexed, 0))
		}
		rows[len(rows)-1] = append(rows[len(rows)-1], element)
	}
	commonWidthString, commonWidth := "", make([]int, 0)
	for _, row := range rows {
		widthString, width := createWidthDescriptor(row)
		if _, ok := stat[widthString]; !ok {
			stat[widthString] = 0
		}
		stat[widthString]++
		if stat[widthString] > stat[commonWidthString] {
			commonWidthString = widthString
			commonWidth = width
		}
	}
	order := make([]int, 0)
	heights := make([]int, 0)
	overrides := make(map[int]LayoutOverride)
	for _, row := range rows {
		rowWidthString, _ := createWidthDescriptor(row)
		rowHeights := getUnique(getHeights(row))
		for _, cell := range row {
			order = append(order, cell.Index)
			if rowWidthString != commonWidthString || len(rowHeights) > 1 {
				overrides[cell.Index] = LayoutOverride{Width: cell.Grid.W, Height: cell.Grid.H}
			}
		}
		if rowWidthString == commonWidthString && len(rowHeights) == 1 {
			heights = append(heights, rowHeights[0])
		}
	}
	if len(getUnique(heights)) == 1 {
		heights = []int{heights[0]}
	}
	return Layout{
		Columns:   commonWidth,
		Heights:   heights,
		Order:     order,
		Overrides: overrides,
	}
}
