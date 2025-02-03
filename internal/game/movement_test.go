package game

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	//"math"
	//"cyber/internal/models"
)

//func (h Hex) Cost(toPosition Hex) float64 {
//	return math.Abs(h.Q-toPosition.Q) + math.Abs(h.R-toPosition.R) + math.Abs(h.Q+h.R-toPosition.Q-toPosition.R)
//}
//
//func (h Hex) Heuristic(toPosition Hex) float64 {
//	return (math.Abs(h.Q-toPosition.Q) + math.Abs(h.R-toPosition.R) + math.Abs(h.Q+h.R-toPosition.Q-toPosition.R)) / 2
//}

func TestHexCost(t *testing.T) {
	tests := []struct {
		name     string
		start    Hex
		end      Hex
		expected float64
	}{
		{
			name:     "Simple case",
			start:    Hex{Q: 0, R: 0},
			end:      Hex{Q: 3, R: 4},
			expected: 14.0, // |0-3| + |0-4| + |0+0-3-4| = 3 + 4 + 7 = 14
		},
		{
			name:     "Same position",
			start:    Hex{Q: 1, R: 1},
			end:      Hex{Q: 1, R: 1},
			expected: 0.0, // |1-1| + |1-1| + |1+1-1-1| = 0 + 0 + 0 = 0
		},
		{
			name:     "Negative coordinates",
			start:    Hex{Q: -1, R: -1},
			end:      Hex{Q: 1, R: 1},
			expected: 8.0, // |-1-1| + |-1-1| + |-1-1-1-1| = 2 + 2 + 4 = 8
		},
		{
			name:     "Large coordinates",
			start:    Hex{Q: 254, R: 400},
			end:      Hex{Q: 200, R: 200},
			expected: 508.0, // |254-200| + |400-200| + |254+400-200-200| = 54 + 200 + 254 = 508
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.Cost(tt.end)
			assert.Equal(t, tt.expected, result, "Cost from %v to %v", tt.start, tt.end)
		})
	}
}

func TestHexHeuristic(t *testing.T) {
	tests := []struct {
		name     string
		start    Hex
		end      Hex
		expected float64
	}{
		{
			name:     "Simple case",
			start:    Hex{Q: 0, R: 0},
			end:      Hex{Q: 3, R: 4},
			expected: 7.0, // |0-3| + |0-4| + |0+0-3-4| = 3 + 4 + 7 = 14 / 2 = 7
		},
		{
			name:     "Same position",
			start:    Hex{Q: 1, R: 1},
			end:      Hex{Q: 1, R: 1},
			expected: 0.0, // |1-1| + |1-1| + |1+1-1-1| = 0 + 0 + 0 = 0 / 2 = 0
		},
		{
			name:     "Negative coordinates",
			start:    Hex{Q: -1, R: -1},
			end:      Hex{Q: 1, R: 1},
			expected: 4.0, // |-1-1| + |-1-1| + |-1-1-1-1| = 2 + 2 + 4 = 8 / 2 = 4
		},
		{
			name:     "Large coordinates",
			start:    Hex{Q: 254, R: 400},
			end:      Hex{Q: 200, R: 200},
			expected: 254.0, // |254-200| + |400-200| + |254+400-200-200| = 54 + 200 + 254 = 508 / 2 = 254
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.start.Heuristic(tt.end)
			assert.Equal(t, tt.expected, result, "Heuristic from %v to %v", tt.start, tt.end)
		})
	}
}

func Test_Neighbours(t *testing.T) {
	tests := []struct {
		name     string
		hex      Hex
		expected []Hex
	}{
		{
			name: "Center hex",
			hex:  Hex{Q: 0, R: 0},
			expected: []Hex{
				{1, 0},  // Сосед справа
				{-1, 0}, // Сосед слева
				{0, 1},  // Сосед сверху-справа
				{0, -1}, // Сосед снизу-слева
				{1, -1}, // Сосед снизу-справа
				{-1, 1}, // Сосед сверху-слева
			},
		},
		{
			name: "Hex with positive coordinates",
			hex:  Hex{Q: 2, R: 3},
			expected: []Hex{
				{3, 3}, // Сосед справа
				{1, 3}, // Сосед слева
				{2, 4}, // Сосед сверху-справа
				{2, 2}, // Сосед снизу-слева
				{3, 2}, // Сосед снизу-справа
				{1, 4}, // Сосед сверху-слева
			},
		},
		{
			name: "Hex with negative coordinates",
			hex:  Hex{Q: -2, R: -3},
			expected: []Hex{
				{-1, -3}, // Сосед справа
				{-3, -3}, // Сосед слева
				{-2, -2}, // Сосед сверху-справа
				{-2, -4}, // Сосед снизу-слева
				{-1, -4}, // Сосед снизу-справа
				{-3, -2}, // Сосед сверху-слева
			},
		},
		{
			name: "Hex with mixed coordinates",
			hex:  Hex{Q: -1, R: 2},
			expected: []Hex{
				{0, 2},  // Сосед справа
				{-2, 2}, // Сосед слева
				{-1, 3}, // Сосед сверху-справа
				{-1, 1}, // Сосед снизу-слева
				{0, 1},  // Сосед снизу-справа
				{-2, 3}, // Сосед сверху-слева
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hex.Neighbours()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Neighbors of %v: expected %v, got %v", tt.hex, tt.expected, result)
			}
		})
	}
}
