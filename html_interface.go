package core

import (
	"bytes"
)

// Html Buffer Interface
type HtmlBuffer interface {
	HtmlAttr() *bytes.Buffer
	Title() *bytes.Buffer
	Head() *bytes.Buffer
	BodyAttr() *bytes.Buffer
	BodyHeader() *bytes.Buffer
	BodyContent() *bytes.Buffer
	BodyFooter() *bytes.Buffer
	BodyJs() *bytes.Buffer
}

// Html Printer Interface
type HtmlPrinter interface {
	GetBuffer() HtmlBuffer
	HtmlAttr(a ...interface{}) (int, error)
	HtmlAttrF(format string, a ...interface{}) (int, error)
	HtmlAttrLn(a ...interface{}) (int, error)
	Title(a ...interface{}) (int, error)
	TitleF(format string, a ...interface{}) (int, error)
	TitleLn(a ...interface{}) (int, error)
	Head(a ...interface{}) (int, error)
	HeadF(format string, a ...interface{}) (int, error)
	HeadLn(a ...interface{}) (int, error)
	BodyAttr(a ...interface{}) (int, error)
	BodyAttrF(format string, a ...interface{}) (int, error)
	BodyAttrLn(a ...interface{}) (int, error)
	BodyHeader(a ...interface{}) (int, error)
	BodyHeaderF(format string, a ...interface{}) (int, error)
	BodyHeaderLn(a ...interface{}) (int, error)
	BodyContent(a ...interface{}) (int, error)
	BodyContentF(format string, a ...interface{}) (int, error)
	BodyContentLn(a ...interface{}) (int, error)
	BodyFooter(a ...interface{}) (int, error)
	BodyFooterF(format string, a ...interface{}) (int, error)
	BodyFooterLn(a ...interface{}) (int, error)
	BodyJs(a ...interface{}) (int, error)
	BodyJsF(format string, a ...interface{}) (int, error)
	BodyJsLn(a ...interface{}) (int, error)
}
