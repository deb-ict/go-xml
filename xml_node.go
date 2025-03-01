package xml

import (
	"github.com/beevik/etree"
)

type XmlNode interface {
	LoadXml(resolver XmlContext, el *etree.Element) error
	GetXml(resolver XmlContext) (*etree.Element, error)
}
