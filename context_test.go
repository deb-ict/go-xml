package xml

import (
	"testing"

	"github.com/beevik/etree"
)

func TestNewContext(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	if ctx.GetDocument() != doc {
		t.Errorf("expected %v, got %v", doc, ctx.GetDocument())
	}
}
