package snake

import (
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/stretchr/testify/assert"
)

func assertGameBoard(t *testing.T, expectedGameBoard, gotGameboard *virtualboard.VirtualBoard, ) bool {
	return assert.Equalf(t, expectedGameBoard, gotGameboard, "Expected:\n%s\nGot:\n%s", expectedGameBoard, gotGameboard)
}
