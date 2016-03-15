package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func testDiffColor(t *testing.T, c Color) {
	p := GenericPlayer{Col:c}

	for i := 0; i < 28; i++ {
		pos := (i+startPos(c))%28
		assert.Equal(t, 28-i, p.diffToHomeStretch(pos))
	}
}

func TestDiff(t *testing.T) {
	testDiffColor(t, Red)
	testDiffColor(t, Purple)
	testDiffColor(t, White)
	testDiffColor(t, Blue)
}
