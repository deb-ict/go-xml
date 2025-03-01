package xml

import (
	"errors"

	"github.com/beevik/etree"
)

var (
	ErrNoTypeConstructor = errors.New("no type constructor")
)

type XmlTypeConstructor func(resolver XmlContext) (XmlNode, error)

type XmlContext interface {
	GetDocument() *etree.Document
	SetNamespacePrefix(prefix string, uri string)
	GetNamespacePrefix(uri string) string
	GetNamespaceUri(prefix string) string
	RegisterTypeConstructor(uri string, tag string, ctor XmlTypeConstructor)
	GetTypeConstructor(uri string, tag string) (XmlTypeConstructor, error)
}

type xmlContext struct {
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

func NewXmlContext(doc *etree.Document) XmlContext {
	return &xmlContext{
		doc:              doc,
		uris:             make(map[string]string),
		prefixes:         make(map[string]string),
		typeConstructors: make([]*xmlTypeEntry, 0),
	}
}

func (resolver *xmlContext) GetDocument() *etree.Document {
	return resolver.doc
}

func (resolver *xmlContext) SetNamespacePrefix(prefix string, uri string) {
	resolver.prefixes[uri] = prefix
	resolver.uris[prefix] = uri
}

func (resolver *xmlContext) GetNamespacePrefix(uri string) string {
	prefix, found := resolver.prefixes[uri]
	if !found {
		return uri
	}
	return prefix
}

func (resolver *xmlContext) GetNamespaceUri(prefix string) string {
	namespaceUri, found := resolver.uris[prefix]
	if !found {
		return prefix
	}
	return namespaceUri
}

func (resolver *xmlContext) RegisterTypeConstructor(uri string, tag string, ctor XmlTypeConstructor) {
	entry, ok := resolver.getTypeConstructor(uri, tag)
	if !ok {
		entry = &xmlTypeEntry{
			uri: uri,
			tag: tag,
		}
		resolver.typeConstructors = append(resolver.typeConstructors, entry)
	}
	entry.constructor = ctor
}

func (resolver *xmlContext) GetTypeConstructor(uri string, tag string) (XmlTypeConstructor, error) {
	entry, ok := resolver.getTypeConstructor(uri, tag)
	if !ok {
		return nil, ErrNoTypeConstructor
	}
	return entry.constructor, nil
}

func (resolver *xmlContext) getTypeConstructor(uri string, tag string) (*xmlTypeEntry, bool) {
	for _, entry := range resolver.typeConstructors {
		if entry.uri == uri && entry.tag == tag {
			return entry, true
		}
	}
	return nil, false
}
