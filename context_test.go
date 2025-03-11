package xml

import (
	"testing"

	"github.com/beevik/etree"
)

func TestNewContextWithDocumentLoadsNametable(t *testing.T) {
	xlm := `<?xml version="1.0"?><root xmlns="http://example.com" xmlns:foo="http://foo.com" xmlns:bar="http://bar.com"><foo:child/><bar:child/></root>`
	doc := etree.NewDocument()
	err := doc.ReadFromString(xlm)
	if err != nil {
		t.Fatal(err)
	}
	ctx := NewContextWithDocument(doc)

	testCases := []struct {
		prefix string
		uri    string
		found  bool
	}{
		{"", "http://example.com", true},
		{"foo", "http://foo.com", true},
		{"bar", "http://bar.com", true},
		{"baz", "", false},
	}
	for _, tc := range testCases {
		prefix, found := ctx.GetNamespacePrefix(tc.uri)
		if tc.found && prefix != tc.prefix {
			t.Fatalf("Expected prefix %s for namespace %s, got %s", tc.prefix, tc.uri, prefix)
		}
		if found != tc.found {
			t.Fatalf("Expected found=%t for namespace %s, got %t", tc.found, tc.uri, found)
		}

		uri, found := ctx.GetNamespaceUri(tc.prefix)
		if tc.found && uri != tc.uri {
			t.Fatalf("Expected uri %s for prefix %s, got %s", tc.uri, tc.prefix, uri)
		}
		if found != tc.found {
			t.Fatalf("Expected found=%t for prefix %s, got %t", tc.found, tc.prefix, found)
		}
	}
}
