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
	SetNamespacePrefix(prefix string, uri string)
	GetNamespacePrefix(uri string) string
	GetNamespaceUri(prefix string) string
	RegisterTypeConstructor(uri string, tag string, ctor XmlTypeConstructor)
	GetTypeConstructor(uri string, tag string) (XmlTypeConstructor, error)
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

func NewContext(doc *etree.Document) Context {
	return &context{
		doc:              doc,
		uris:             make(map[string]string),
		prefixes:         make(map[string]string),
		typeConstructors: make([]*xmlTypeEntry, 0),
	}
}

func (ctx *context) GetDocument() *etree.Document {
	return ctx.doc
}

func (ctx *context) SetNamespacePrefix(prefix string, uri string) {
	ctx.prefixes[uri] = prefix
	ctx.uris[prefix] = uri
}

func (ctx *context) GetNamespacePrefix(uri string) string {
	prefix, found := ctx.prefixes[uri]
	if !found {
		return uri
	}
	return prefix
}

func (ctx *context) GetNamespaceUri(prefix string) string {
	namespaceUri, found := ctx.uris[prefix]
	if !found {
		return prefix
	}
	return namespaceUri
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

func (ctx *context) getTypeConstructor(uri string, tag string) (*xmlTypeEntry, bool) {
	for _, entry := range ctx.typeConstructors {
		if entry.uri == uri && entry.tag == tag {
			return entry, true
		}
	}
	return nil, false
}
