package sqlite

import (
	"testing"

	"github.com/louisbranch/edulab"
)

func TestDBInterface(t *testing.T) {
	var _ edulab.Database = &DB{}

}
