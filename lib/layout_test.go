package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlignGrid(t *testing.T) {
	grid := AlignGrid([]Box{
		{X: 0, Y: 0, H: 1, W: 1},
		{X: 0, Y: 1, H: 1, W: 1},
		{X: 1, Y: 0, H: 2, W: 2},
		{X: 0, Y: 2, H: 1, W: 2},
		{X: 2, Y: 2, H: 1, W: 1},
	})
	t.Log(grid)
}

func TestFindGroupSeparators(t *testing.T) {
	separators := FindGroupSeparators([]Box{
		{X: 0, Y: 0, H: 1, W: 1},
		{X: 0, Y: 1, H: 1, W: 1},
		{X: 1, Y: 0, H: 2, W: 2},
		{X: 0, Y: 2, H: 1, W: 2},
		{X: 2, Y: 2, H: 1, W: 1},
	}, func(b Box) (int, int) { return b.Y, b.Y + b.H })
	require.Equal(t, separators, []int{0, 2, 3})
}

func TestSeparateGrid(t *testing.T) {
	grid := []Box{
		{X: 0, Y: 0, H: 1, W: 1},
		{X: 0, Y: 1, H: 1, W: 1},
		{X: 1, Y: 0, H: 2, W: 2},
		{X: 0, Y: 2, H: 1, W: 2},
		{X: 2, Y: 2, H: 1, W: 1},
	}
	separators := FindGroupSeparators(grid, func(b Box) (int, int) { return b.Y, b.Y + b.H })
	groups := SeparateGrid(grid, separators)
	t.Log(groups)
	require.Equal(t, groups, [][]Box{
		{{X: 0, Y: 0, H: 1, W: 1}, {X: 1, Y: 0, H: 2, W: 2}, {X: 0, Y: 1, H: 1, W: 1}},
		{{X: 0, Y: 2, H: 1, W: 2}, {X: 2, Y: 2, H: 1, W: 1}},
	})
}

func TestAlignGroup(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		grid := []Box{
			{Id: 0, X: 0, Y: 0, H: 1, W: 1},
			{Id: 1, X: 0, Y: 1, H: 1, W: 1},
			{Id: 2, X: 1, Y: 0, H: 2, W: 2},
			{Id: 3, X: 0, Y: 2, H: 1, W: 2},
			{Id: 4, X: 2, Y: 2, H: 1, W: 1},
		}
		separators := FindGroupSeparators(grid, func(b Box) (int, int) { return b.Y, b.Y + b.H })
		groups := SeparateGrid(grid, separators)
		layout := AlignGroup(groups[0])
		require.Equal(t, layout, GroupLayout{
			Height:    2,
			Widths:    []int{1, 2},
			Panels:    []Box{{Id: 0, X: 0, Y: 0, H: 1, W: 1}, {Id: 1, X: 0, Y: 1, H: 1, W: 1}, {Id: 2, X: 1, Y: 0, H: 2, W: 2}},
			Overrides: map[int]LayoutOverride{0: {0, 1}, 1: {0, 1}},
		})
		t.Log(layout)
	})
	t.Run("hard", func(t *testing.T) {
		grid := []Box{
			{Id: 0, X: 0, Y: 0, H: 1, W: 1},
			{Id: 1, X: 1, Y: 0, H: 1, W: 1},
			{Id: 2, X: 0, Y: 0, H: 2, W: 2},
		}
		separators := FindGroupSeparators(grid, func(b Box) (int, int) { return b.Y, b.Y + b.H })
		groups := SeparateGrid(grid, separators)
		layout := AlignGroup(groups[0])
		require.Equal(t, layout, GroupLayout{
			Height: 24,
			Widths: []int{24},
			Panels: []Box{{Id: 0, X: 0, Y: 0, H: 1, W: 1}, {Id: 2, X: 0, Y: 0, H: 2, W: 2}, {Id: 1, X: 1, Y: 0, H: 1, W: 1}},
		})
		t.Log(layout)
	})
}
