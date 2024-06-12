package lib

import (
	"encoding/json"
	"fmt"
	"sort"

	_ "go.uber.org/zap"
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

type GridLayout struct {
	Height int
	Widths []int
	Groups []GroupLayout
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

func FindGroupSeparators(grid []Box, rangeF func(Box) (int, int)) []int {
	type Endpoint struct{ position, sign int }
	endpoints := make([]Endpoint, 0)
	for _, box := range grid {
		l, r := rangeF(box)
		endpoints = append(endpoints, Endpoint{l, +1}, Endpoint{r, -1})
	}
	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].position != endpoints[j].position {
			return endpoints[i].position < endpoints[j].position
		}
		return endpoints[i].sign < endpoints[j].sign
	})

	separators := []int{0}
	balance := 0
	for _, endpoint := range endpoints {
		balance += endpoint.sign
		if balance == 0 {
			separators = append(separators, endpoint.position)
		}
	}
	Logger.Infof("found group separators: %v", separators)
	return separators
}

func SeparateGrid(grid []Box, separators []int) [][]Box {
	sort.Slice(grid, func(i, j int) bool { return grid[i].Y < grid[j].Y })
	groups := [][]Box{{}}
	separator := 0
	for _, box := range grid {
		for separators[separator+1] < box.Y+box.H {
			separator++
			groups = append(groups, []Box{})
		}
		groups[separator] = append(groups[separator], box)
	}
	return groups
}

type GroupLayout struct {
	Height    int
	Widths    []int
	Panels    []Box
	Overrides map[int]LayoutOverride
}

func ExtractPeriod[T comparable](sequence []T) []T {
	for period := 1; period < len(sequence); period++ {
		ok := true
		for i := period; ok && i < len(sequence); i++ {
			ok = ok && (sequence[i] == sequence[i-period])
		}
		if ok {
			return sequence[:period]
		}
	}
	return sequence
}

func AlignGroup(group []Box) GroupLayout {
	minY, maxY := group[0].Y, group[0].Y
	for _, box := range group {
		if minY > box.Y {
			minY = box.Y
		}
		if maxY < box.Y+box.H {
			maxY = box.Y + box.H
		}
	}
	height := maxY - minY

	xSeparators := FindGroupSeparators(group, func(b Box) (int, int) { return b.X, b.X + b.W })
	sort.Slice(group, func(i, j int) bool {
		if group[i].X != group[j].X {
			return group[i].X < group[j].X
		}
		return group[i].Y < group[j].X
	})

	columns := make([]int, 0)
	for i := 0; i < len(xSeparators)-1; i++ {
		columns = append(columns, xSeparators[i+1]-xSeparators[i])
	}

	overrides := make(map[int]LayoutOverride)
	separator := 0
	for _, box := range group {
		for xSeparators[separator+1] < box.X+box.W {
			separator++
		}
		if box.X != xSeparators[separator] || box.X+box.W != xSeparators[separator+1] {
			Logger.Warnf("failed to align group into grid, fallback to dummy layout")
			return GroupLayout{Height: 8 * len(group), Widths: []int{24}, Panels: group}
		}
		if box.H != height {
			overrides[box.Id] = LayoutOverride{Height: box.H}
		}
	}

	return GroupLayout{
		Height:    height,
		Widths:    columns,
		Panels:    group,
		Overrides: overrides,
	}
}

type FrequencyStat[T any] struct {
	Instances map[string]T
	Frequency map[string]int
}

func NewFrequencyStat[T any]() *FrequencyStat[T] {
	return &FrequencyStat[T]{Instances: make(map[string]T), Frequency: make(map[string]int)}
}

func (s *FrequencyStat[T]) Add(value T) {
	serialized, err := json.Marshal(value)
	if err != nil {
		Logger.Fatalf("failed to serialize instance %v", value)
	}
	s.Instances[string(serialized)] = value
	s.Frequency[string(serialized)]++
}

func (s *FrequencyStat[T]) GetTop() T {
	bestInstance, bestF := *new(T), 0
	keys := make([]string, 0, len(s.Frequency))
	for key := range s.Frequency {
		keys = append(keys, key)
	}
	// note: guarantee deterministic Top calculation
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, key := range keys {
		f := s.Frequency[key]
		if bestF < f {
			bestF = f
			bestInstance = s.Instances[key]
		}
	}
	return bestInstance
}

func AlignGrid(grid []Box) GridLayout {
	ySeparators := FindGroupSeparators(grid, func(b Box) (int, int) { return b.Y, b.Y + b.H })
	groups := SeparateGrid(grid, ySeparators)
	groupLayouts := make([]GroupLayout, 0)

	heightsFrequency := NewFrequencyStat[int]()
	columnsFrequency := NewFrequencyStat[[]int]()
	for _, group := range groups {
		groupLayout := AlignGroup(group)
		groupLayouts = append(groupLayouts, groupLayout)
		heightsFrequency.Add(groupLayout.Height)
		columnsFrequency.Add(groupLayout.Widths)
	}

	return GridLayout{
		Height: heightsFrequency.GetTop(),
		Widths: columnsFrequency.GetTop(),
		Groups: groupLayouts,
	}
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
