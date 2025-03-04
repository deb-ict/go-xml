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

func TestNamespacePrefix(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	ctx.SetNamespacePrefix("foo", "http://foo.com")
	if ctx.GetNamespacePrefix("http://foo.com") != "foo" {
		t.Errorf("expected %v, got %v", "foo", ctx.GetNamespacePrefix("http://foo.com"))
	}
}

func TestNamespacePrefixNotFound(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	if ctx.GetNamespacePrefix("http://foo.com") != "http://foo.com" {
		t.Errorf("expected %v, got %v", "http://foo.com", ctx.GetNamespacePrefix("http://foo.com"))
	}
}

func TestNamespaceUri(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	ctx.SetNamespacePrefix("foo", "http://foo.com")
	if ctx.GetNamespaceUri("foo") != "http://foo.com" {
		t.Errorf("expected %v, got %v", "http://foo.com", ctx.GetNamespaceUri("foo"))
	}
}

func TestNamespaceUriNotFound(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	if ctx.GetNamespaceUri("foo") != "foo" {
		t.Errorf("expected %v, got %v", "foo", ctx.GetNamespaceUri("foo"))
	}
}

func TestTypeConstructor(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	ctor := func(ctx Context) (Node, error) {
		return nil, nil
	}
	ctx.RegisterTypeConstructor("http://foo.com", "bar", ctor)
	got, err := ctx.GetTypeConstructor("http://foo.com", "bar")
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if got == nil {
		t.Error("ctor is nil")
	}
}

func TestTypeConstructorNotFound(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	ctor := func(ctx Context) (Node, error) {
		return nil, nil
	}
	ctx.RegisterTypeConstructor("http://foo.com", "bar", ctor)
	_, err := ctx.GetTypeConstructor("http://bar.com", "foo")
	if err != ErrNoTypeConstructor {
		t.Errorf("expected %v, got %v", ErrNoTypeConstructor, err)
	}
}

func TestElementTypeConstructor(t *testing.T) {
	doc := etree.NewDocument()
	ctx := NewContext(doc)
	ctor := func(ctx Context) (Node, error) {
		return nil, nil
	}
	el := etree.NewElement("foo:bar")
	el.CreateAttr("xmlns:foo", "http://foo.com")

	ctx.RegisterTypeConstructor("http://foo.com", "bar", ctor)
	got, err := ctx.GetElementTypeConstructor(el)
	if err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}
	if got == nil {
		t.Error("ctor is nil")
	}
}
