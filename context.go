package xml

import (
	"errors"

	"github.com/beevik/etree"
)

var (
	ErrNoTypeConstructor = errors.New("no type constructor")
)

type XmlTypeConstructor func(ctx Context) (Node, error)

type Context interface {
	GetDocument() *etree.Document
	WriteToString() (string, error)
	SetNamespacePrefix(prefix string, uri string)
	GetNamespacePrefix(uri string) (string, bool)
	GetNamespaceUri(prefix string) (string, bool)
	RegisterTypeConstructor(uri string, tag string, ctor XmlTypeConstructor)
	GetTypeConstructor(uri string, tag string) (XmlTypeConstructor, error)
	GetElementTypeConstructor(el *etree.Element) (XmlTypeConstructor, error)
}

type context struct {
	doc              *etree.Document
	uris             map[string]string
	prefixes         map[string]string
	typeConstructors []*xmlTypeEntry
}

type xmlTypeEntry struct {
	uri         string
	tag         string
	constructor XmlTypeConstructor
}

func NewContext() Context {
	return NewContextWithDocument(etree.NewDocument())
}

func NewContextWithDocument(doc *etree.Document) Context {
	ctx := &context{
		doc:              doc,
		uris:             make(map[string]string),
		prefixes:         make(map[string]string),
		typeConstructors: make([]*xmlTypeEntry, 0),
	}
	ctx.loadNametable(doc.Root())

	return ctx
}

func (ctx *context) GetDocument() *etree.Document {
	return ctx.doc
}

func (ctx *context) WriteToString() (string, error) {
	ctx.writeNametable(ctx.doc.Root())
	return ctx.doc.WriteToString()
}

func (ctx *context) SetNamespacePrefix(prefix string, uri string) {
	ctx.prefixes[uri] = prefix
	ctx.uris[prefix] = uri
}

func (ctx *context) GetNamespacePrefix(uri string) (string, bool) {
	prefix, found := ctx.prefixes[uri]
	if !found {
		return "", false
	}
	return prefix, true
}

func (ctx *context) GetNamespaceUri(prefix string) (string, bool) {
	namespaceUri, found := ctx.uris[prefix]
	if !found {
		return "", false
	}
	return namespaceUri, true
}

func (ctx *context) RegisterTypeConstructor(uri string, tag string, ctor XmlTypeConstructor) {
	entry, ok := ctx.getTypeConstructor(uri, tag)
	if !ok {
		entry = &xmlTypeEntry{
			uri: uri,
			tag: tag,
		}
		ctx.typeConstructors = append(ctx.typeConstructors, entry)
	}
	entry.constructor = ctor
}

func (ctx *context) GetTypeConstructor(uri string, tag string) (XmlTypeConstructor, error) {
	entry, ok := ctx.getTypeConstructor(uri, tag)
	if !ok {
		return nil, ErrNoTypeConstructor
	}
	return entry.constructor, nil
}

func (ctx *context) GetElementTypeConstructor(el *etree.Element) (XmlTypeConstructor, error) {
	return ctx.GetTypeConstructor(el.NamespaceURI(), el.Tag)
}

func (ctx *context) getTypeConstructor(uri string, tag string) (*xmlTypeEntry, bool) {
	for _, entry := range ctx.typeConstructors {
		if entry.uri == uri && entry.tag == tag {
			return entry, true
		}
	}
	return nil, false
}

func (ctx *context) loadNametable(el *etree.Element) {
	if el == nil {
		return
	}
	for _, attr := range el.Attr {
		if attr.Space == "xmlns" {
			ctx.SetNamespacePrefix(attr.Key, attr.Value)
		}
		if attr.Space == "" && attr.Key == "xmlns" {
			ctx.SetNamespacePrefix("", attr.Value)
		}
	}

	for _, child := range el.ChildElements() {
		ctx.loadNametable(child)
	}
}

func (ctx *context) writeNametable(el *etree.Element) {
	if el == nil {
		return
	}
	for prefix, uri := range ctx.uris {
		if prefix == "" {
			el.CreateAttr("xmlns", uri)
		} else {
			el.CreateAttr("xmlns:"+prefix, uri)
		}
	}
	el.SortAttrs()
}
