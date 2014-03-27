package core

import (
	"bytes"
	"fmt"
	"html"
)

type html5Buffer struct {
	htmlAttr, title, head, bodyAttr, bodyHeader, bodyContent, bodyFooter, bodyJs *bytes.Buffer
}

func (h html5Buffer) HtmlAttr() *bytes.Buffer {
	return h.htmlAttr
}

func (h html5Buffer) Title() *bytes.Buffer {
	return h.title
}

func (h html5Buffer) Head() *bytes.Buffer {
	return h.head
}

func (h html5Buffer) BodyAttr() *bytes.Buffer {
	return h.bodyAttr
}

func (h html5Buffer) BodyHeader() *bytes.Buffer {
	return h.bodyHeader
}

func (h html5Buffer) BodyContent() *bytes.Buffer {
	return h.bodyContent
}

func (h html5Buffer) BodyFooter() *bytes.Buffer {
	return h.bodyFooter
}

func (h html5Buffer) BodyJs() *bytes.Buffer {
	return h.bodyJs
}

// A Restful HTML Template Controller
type MethodHtml5 struct {
	Method
	buffers      html5Buffer
	_init        bool
	onInitFunc   []func(HtmlPrinter, *Context)
	onFinishFunc []func(HtmlPrinter, *Context)
}

func (me *MethodHtml5) init() {
	if me._init {
		return
	}
	me._init = true

	me.buffers = html5Buffer{&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}}

	for _, fn := range me.onInitFunc {
		fn(me, me.C)
	}
}

func (me *MethodHtml5) RegOnInitFunc(fns ...func(HtmlPrinter, *Context)) {
	if me.onInitFunc == nil {
		me.onInitFunc = []func(HtmlPrinter, *Context){}
	}

	me.onInitFunc = append(me.onInitFunc, fns...)
}

func (me *MethodHtml5) RegOnFinishFunc(fns ...func(HtmlPrinter, *Context)) {
	if me.onFinishFunc == nil {
		me.onFinishFunc = []func(HtmlPrinter, *Context){}
	}

	me.onFinishFunc = append(me.onFinishFunc, fns...)
}

func (me *MethodHtml5) GetBuffer() HtmlBuffer {
	me.init()
	return me.buffers
}

func (me *MethodHtml5) HtmlAttr(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.htmlAttr, a...)
}

func (me *MethodHtml5) HtmlAttrF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.htmlAttr, format, a...)
}

func (me *MethodHtml5) HtmlAttrLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.htmlAttr, a...)
}

func (me *MethodHtml5) Title(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.title, a...)
}

func (me *MethodHtml5) TitleF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.title, format, a...)
}

func (me *MethodHtml5) TitleLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.title, a...)
}

func (me *MethodHtml5) Head(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.head, a...)
}

func (me *MethodHtml5) HeadF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.head, format, a...)
}

func (me *MethodHtml5) HeadLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.head, a...)
}

func (me *MethodHtml5) BodyAttr(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.bodyAttr, a...)
}

func (me *MethodHtml5) BodyAttrF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.bodyAttr, format, a...)
}

func (me *MethodHtml5) BodyAttrLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.bodyAttr, a...)
}

func (me *MethodHtml5) BodyHeader(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.bodyHeader, a...)
}

func (me *MethodHtml5) BodyHeaderF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.bodyHeader, format, a...)
}

func (me *MethodHtml5) BodyHeaderLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.bodyHeader, a...)
}

func (me *MethodHtml5) BodyContent(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.bodyContent, a...)
}

func (me *MethodHtml5) BodyContentF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.bodyContent, format, a...)
}

func (me *MethodHtml5) BodyContentLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.bodyContent, a...)
}

func (me *MethodHtml5) BodyFooter(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.bodyFooter, a...)
}

func (me *MethodHtml5) BodyFooterF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.bodyFooter, format, a...)
}

func (me *MethodHtml5) BodyFooterLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.bodyFooter, a...)
}

func (me *MethodHtml5) BodyJs(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprint(me.buffers.bodyJs, a...)
}

func (me *MethodHtml5) BodyJsF(format string, a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintf(me.buffers.bodyJs, format, a...)
}

func (me *MethodHtml5) BodyJsLn(a ...interface{}) (int, error) {
	me.init()
	return fmt.Fprintln(me.buffers.bodyJs, a...)
}

func (me *MethodHtml5) Finish() {
	if !me._init {
		return
	}

	for _, fn := range me.onFinishFunc {
		fn(me, me.C)
	}

	w := me.C.Pub.Writers["gzip"]
	if w == nil {
		w = me.C.Res
	}

	es := html.EscapeString

	fmt.Fprint(w, `<!DOCTYPE html>
<html `)

	me.buffers.htmlAttr.WriteTo(w)

	fmt.Fprint(w, `>
<head>
<title>`, es(me.buffers.title.String()), `</title>
`)
	me.buffers.head.WriteTo(w)

	fmt.Fprint(w, `
</head>
<body `)

	me.buffers.bodyAttr.WriteTo(w)

	fmt.Fprint(w, `>
`)

	me.buffers.bodyHeader.WriteTo(w)
	me.buffers.bodyContent.WriteTo(w)
	me.buffers.bodyFooter.WriteTo(w)
	me.buffers.bodyJs.WriteTo(w)

	fmt.Fprint(w, `
</body>
</html>`)
}

// Alais of MethodHtml5
type VerbHtml5 MethodHtml5

func HtmlAttrLang(code string) func(HtmlPrinter, *Context) {
	code = fmt.Sprintf(`lang="%s" `, html.EscapeString(code))
	return func(h HtmlPrinter, c *Context) {
		h.HtmlAttr(code)
	}
}

func HtmlAttrDir(dir string) func(HtmlPrinter, *Context) {
	dir = fmt.Sprintf(`dir="%s" `, html.EscapeString(dir))
	return func(h HtmlPrinter, c *Context) {
		h.HtmlAttr(dir)
	}
}
