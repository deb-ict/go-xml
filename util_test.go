package xml

import (
	"testing"

	"github.com/beevik/etree"
)

func TestCreateElement(t *testing.T) {
	testCases := []struct {
		tag           string
		namespaceUri  string
		expectedSpace string
		expectedError error
	}{
		{
			tag:           "test",
			namespaceUri:  "",
			expectedSpace: "",
			expectedError: nil,
		},
		{
			tag:           "test",
			namespaceUri:  "http://example.com",
			expectedSpace: "foo",
			expectedError: nil,
		},
	}

	context := NewContext(etree.NewDocument())
	context.SetNamespacePrefix("foo", "http://example.com")
	for _, tc := range testCases {
		el := CreateElement(context, tc.tag, tc.namespaceUri)
		if el.Space != tc.expectedSpace {
			t.Errorf("expected %v, got %v", tc.expectedSpace, el.Space)
		}
	}
}

func TestValidateElement(t *testing.T) {
	testCases := []struct {
		el            *etree.Element
		tag           string
		namespaceUri  string
		expectedError error
	}{
		{
			el:            nil,
			tag:           "test",
			namespaceUri:  "",
			expectedError: ErrElementIsNil,
		},
		{
			el:            etree.NewElement("test"),
			tag:           "test",
			namespaceUri:  "",
			expectedError: nil,
		},
		{
			el:            etree.NewElement("test"),
			tag:           "test",
			namespaceUri:  "http://example.com",
			expectedError: ErrInvalidElementTag,
		},
		{
			el:            etree.NewElement("test"),
			tag:           "test2",
			namespaceUri:  "",
			expectedError: ErrInvalidElementTag,
		},
	}
	for _, tc := range testCases {
		err := ValidateElement(tc.el, tc.tag, tc.namespaceUri)
		if err != tc.expectedError {
			t.Errorf("expected %v, got %v", tc.expectedError, err)
		}
	}
}

func TestGetSingleChildElement(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")
	child := root.CreateElement("child")
	child.CreateElement("subchild")
	child.CreateElement("subchild")

	testCases := []struct {
		el            *etree.Element
		tag           string
		namespaceUri  string
		expectedError error
	}{
		{
			el:            root,
			tag:           "child",
			namespaceUri:  "",
			expectedError: nil,
		},
		{
			el:            root,
			tag:           "child",
			namespaceUri:  "http://example.com",
			expectedError: ErrChildElementNotFound,
		},
		{
			el:            root,
			tag:           "child2",
			namespaceUri:  "",
			expectedError: ErrChildElementNotFound,
		},
		{
			el:            child,
			tag:           "subchild",
			namespaceUri:  "",
			expectedError: ErrMultipleChildElementsFound,
		},
	}
	for _, tc := range testCases {
		_, err := GetSingleChildElement(tc.el, tc.tag, tc.namespaceUri)
		if err != tc.expectedError {
			t.Errorf("expected %v, got %v", tc.expectedError, err)
		}
	}
}

func TestGetOptionalSingleChildElement(t *testing.T) {
	doc := etree.NewDocument()
	root := doc.CreateElement("root")
	child := root.CreateElement("child")
	child.CreateElement("subchild")
	child.CreateElement("subchild")

	testCases := []struct {
		el            *etree.Element
		tag           string
		namespaceUri  string
		expectedError error
	}{
		{
			el:            root,
			tag:           "child",
			namespaceUri:  "",
			expectedError: nil,
		},
		{
			el:            root,
			tag:           "child",
			namespaceUri:  "http://example.com",
			expectedError: nil,
		},
		{
			el:            root,
			tag:           "child2",
			namespaceUri:  "",
			expectedError: nil,
		},
		{
			el:            child,
			tag:           "subchild",
			namespaceUri:  "",
			expectedError: ErrMultipleChildElementsFound,
		},
	}
	for _, tc := range testCases {
		_, err := GetOptionalSingleChildElement(tc.el, tc.tag, tc.namespaceUri)
		if err != tc.expectedError {
			t.Errorf("expected %v, got %v", tc.expectedError, err)
		}
	}
}
