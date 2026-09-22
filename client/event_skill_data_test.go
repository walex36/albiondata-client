package client

import (
	"testing"
)

func TestParseFame(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"[[12345]]", 12345},
		{"[500]", 500},
		{"0", 0},
		{"  [[9999]]  ", 9999},
		{"", 0},
		{"invalid", 0},
		{"[[0]]", 0},
	}

	for _, tt := range tests {
		got := parseFame(tt.input)
		if got != tt.want {
			t.Errorf("parseFame(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestEventSkillData_Process(t *testing.T) {
	event := eventSkillData{
		SkillIds:    []int{10, 20},
		Levels:      []int{1, 2},
		Percentages: []float64{0.5, 0.75},
		Fame:        []string{"[[100]]", "[[250]]"},
	}

	state := &albionState{
		CharacterName: "TestPlayer",
	}

	// Make sure Process executes without error or panic
	event.Process(state)
}
