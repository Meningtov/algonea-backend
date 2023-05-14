package slice_test

import (
	"testing"

	"github.com/Meningtov/algonea-backend/lang/slice"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tests := map[string]struct {
		gotSlice     []uint64
		gotElement   uint64
		wantContains bool
	}{
		"contains": {
			gotSlice:     []uint64{1, 2, 3},
			gotElement:   3,
			wantContains: true,
		},
		"does not contain": {
			gotSlice:     []uint64{1, 2, 3},
			gotElement:   4,
			wantContains: false,
		},
		"does not contain with empty slice": {
			gotSlice:     []uint64{},
			gotElement:   4,
			wantContains: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			contains := slice.Contains(test.gotSlice, test.gotElement)
			assert.Equal(t, test.wantContains, contains)
		})
	}
}
