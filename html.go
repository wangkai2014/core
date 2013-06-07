package core

import (
	"bytes"
	html "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
)

type HtmlMiddleware struct {
	Middleware
}

func (ht *HtmlMiddleware) Html() {
	c := ht.C
	// HTML Marksafe
	c.Pub.HtmlFunc["html"] = func(str string) html.HTML {
		return html.HTML(str)
	}

	// HTML Attr MarkSafe
	c.Pub.HtmlFunc["htmlattr"] = func(str string) html.HTMLAttr {
		return html.HTMLAttr(str)
	}

	// JS Marksafe
	c.Pub.HtmlFunc["js"] = func(str string) html.JS {
		return html.JS(str)
	}

	// JS String Marksafe
	c.Pub.HtmlFunc["jsstr"] = func(str string) html.JSStr {
		return html.JSStr(str)
	}

	// CSS Marksafe
	c.Pub.HtmlFunc["css"] = func(str string) html.CSS {
		return html.CSS(str)
	}
}

func init() {
	MainMiddlewares.Register(&HtmlMiddleware{})
}

type Html struct {
	c *Core
}

func (c *Core) Html() Html {
	return Html{c}
}

func (h Html) render(htmlstr string, value_map interface{}, buf io.Writer) {
	c := h.c
	if buf == nil {
		// To prevent headers from being sent too early.
		buf = &bytes.Buffer{}
		defer func(buf *bytes.Buffer) {
			io.Copy(c, buf)
			buf.Reset()
		}(buf.(*bytes.Buffer))
	}
	t := html.Must(html.New("html").Funcs(c.Pub.HtmlFunc).Parse(htmlstr))
	err := t.Execute(buf, value_map)
	c.Check(err)
}

// Render HTML
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
func (h Html) Render(htmlstr string, value_map interface{}) string {
	buf := &bytes.Buffer{}
	defer buf.Reset()
	h.render(htmlstr, value_map, buf)
	return buf.String()
}

// Render HTML and Send Response to Client
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
func (h Html) RenderSend(htmlstr string, value_map interface{}) {
	h.render(htmlstr, value_map, nil)
}

type htmlFileCacheStruct struct {
	content string
	expire  time.Time
}

// Get HTML File
//
// Note: Can also be used to get other kind of files. DO NOT USE THIS WITH LARGE FILES.
func (h Html) GetFile(htmlfile string) string {
	var content string
	var content_in_byte []byte
	var err error

	htmlFileCache.Lock()
	defer htmlFileCache.Unlock()

	switch t := htmlFileCache.m[htmlfile].(type) {
	case htmlFileCacheStruct:
		if time.Now().Unix() > t.expire.Unix() {
			goto getfile_and_cache
		}
		content = t.content
		goto return_content
	}

getfile_and_cache:
	content_in_byte, err = ioutil.ReadFile(htmlfile)
	if err != nil {
		return err.Error()
	}
	content = string(content_in_byte)
	if !DEBUG {
		htmlFileCache.m[htmlfile] = htmlFileCacheStruct{content, time.Now().Add(HtmlTemplateCacheExpire)}
	}

return_content:
	return content
}

// Render HTML File
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
// DO NOT USE THIS WITH LARGE FILES.
func (h Html) RenderFile(htmlfile string, value_map interface{}) string {
	return h.Render(h.GetFile(htmlfile), value_map)
}

// Render HTML File and Send Response to Client
//
// Note: Marksafe functions/filters avaliable are 'html', 'htmlattr', 'js' and 'jsattr'.
// DO NOT USE THIS WITH LARGE FILES.
func (h Html) RenderFileSend(htmlfile string, value_map interface{}) {
	h.RenderSend(h.GetFile(htmlfile), value_map)
}

func (h Html) ParseFiles(filenames ...string) *html.Template {
	t := html.New("html").Funcs(h.c.Pub.HtmlFunc)

	for _, filename := range filenames {
		html.Must(t.Parse(h.GetFile(filename)))
	}

	return t
}

func (h Html) ParseGlob(pattern string) *html.Template {
	htmlGlobLocker.Lock()
	defer htmlGlobLocker.Unlock()

	if len(htmlGlobLocker.filenames[pattern]) > 0 {
		return h.ParseFiles(htmlGlobLocker.filenames[pattern]...)
	}

	filenames, err := filepath.Glob(pattern)
	h.c.Check(err)

	if !DEBUG {
		htmlGlobLocker.filenames[pattern] = filenames
	}

	return h.ParseFiles(filenames...)
}

type htmlDefault struct {
	filenames []string
	pattern   string
	template  *html.Template
}

func (h Html) init() {
	if h.c.pri.html == nil {
		h.c.pri.html = &htmlDefault{}
	}
}

func (h Html) SetDefaultFiles(filenames ...string) {
	h.init()
	h.c.pri.html.filenames = filenames
}

func (h Html) SetDefaultGlob(pattern string) {
	h.init()
	h.c.pri.html.pattern = pattern
}

func (h Html) Default() *html.Template {
	if h.c.pri.html == nil {
		panic(ErrorStr("HTML: Default Template is not set!"))
	}

	if h.c.pri.html.template != nil {
		goto gotoreturn
	}

	if len(h.c.pri.html.filenames) > 0 {
		h.c.pri.html.template = h.ParseFiles(h.c.pri.html.filenames...)
		goto gotoreturn
	}

	if h.c.pri.html.pattern != "" {
		h.c.pri.html.template = h.ParseGlob(h.c.pri.html.pattern)
		goto gotoreturn
	}

gotoreturn:
	return h.c.pri.html.template
}

func (h Html) DefaultRender(name string, data interface{}) string {
	b := &bytes.Buffer{}
	defer b.Reset()
	h.Default().ExecuteTemplate(b, name, data)
	return b.String()
}

func (h Html) DefaultRenderSend(name string, data interface{}) {
	b := &bytes.Buffer{}
	defer func(buf *bytes.Buffer) {
		io.Copy(h.c, buf)
		buf.Reset()
	}(b)
	h.Default().ExecuteTemplate(b, name, data)
}