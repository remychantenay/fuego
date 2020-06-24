package internal

import (
	"testing"
)

func TestCalculateRequiredBatches(t *testing.T) {

	tests := []struct {
		description string
		with        int
		want        int
	}{
		{
			description: "With 0",
			with:        0,
			want:        0,
		},
		{
			description: "With 1",
			with:        1,
			want:        1,
		},
		{
			description: "With 500",
			with:        500,
			want:        1,
		},
		{
			description: "With 501",
			with:        501,
			want:        2,
		},
		{
			description: "With 3100",
			with:        3100,
			want:        7,
		},
	}

	for _, test := range tests {
		result := CalculateRequiredBatches(test.with)

		if result != test.want {
			t.Fatalf("%s -> Got %d but expected %d", test.description, result, test.want)
		}
	}
}
