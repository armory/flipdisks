package snake

import (
	"testing"

	"flipdisks/pkg/virtualboard"
	"github.com/stretchr/testify/assert"
)

func assertGameBoard(t *testing.T, expectedGameBoard, gotGameBoard *virtualboard.VirtualBoard, ) bool {
	return assert.Equalf(t, expectedGameBoard, gotGameBoard, "Expected:\n%s\nGot:\n%s", expectedGameBoard, gotGameBoard)
}
