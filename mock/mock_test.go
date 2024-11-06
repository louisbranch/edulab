package mock

import (
	"testing"

	"github.com/louisbranch/edulab/web"
)

func TestMockSatisfiesInterfaces(t *testing.T) {
	var _ web.Template = &Template{}
}
