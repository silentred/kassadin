package rotator

import (
	"io"
)

type RotatorWriter interface {
	io.Writer
}
