package xml

import (
	"github.com/beevik/etree"
)

type Node interface {
	LoadXml(ctx Context, el *etree.Element) error
	GetXml(ctx Context) (*etree.Element, error)
}
