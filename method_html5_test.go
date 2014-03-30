package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MethodHtml5Dummy struct {
	MethodHtml5
}

func (me *MethodHtml5Dummy) Prepare() {
	me.RegOnInitFunc(me.C.App.Data("initFunc").(func(HtmlPrinter, *Context)))
	me.RegOnFinishFunc(me.C.App.Data("finishFunc").(func(HtmlPrinter, *Context)))
}

func (me *MethodHtml5Dummy) Get() {
	me.Title("title")
	me.TitleF("%s", "A")
	me.TitleLn()
	me.BodyHeader("bodyHeader")
	me.BodyHeaderF("%s", "A")
	me.BodyHeaderLn()
	me.BodyContent("bodyContent")
	me.BodyContentF("%s", "A")
	me.BodyContentLn()
	me.BodyFooter("bodyFooter")
	me.BodyFooterF("%s", "A")
	me.BodyFooterLn()
	me.BodyJs("bodyJs")
	me.BodyJsF("%s", "A")
	me.BodyJsLn()
}

func TestMethodHtml5(t *testing.T) {
	App := NewApp()

	App.Debug = true

	method := &MethodHtml5Dummy{}

	App.DataSet("initFunc", func(h HtmlPrinter, c *Context) {
		h.HtmlAttr(`lang="en`)
		h.HtmlAttrF("%s", `"`)
		h.HtmlAttrLn()
		h.Head("head")
		h.HeadF("%s", "A")
		h.HeadLn()
		h.BodyAttr(`bgcolor="blue`)
		h.BodyAttrF("%s", `"`)
		h.BodyAttrLn()
	})

	App.DataSet("finishFunc", func(h HtmlPrinter, c *Context) {
		buf := h.GetBuffer()

		ln := `
`
		if buf.HtmlAttr().String() != `lang="en"`+ln {
			t.Fail()
		}

		if buf.Title().String() != "titleA"+ln {
			t.Fail()
		}

		if buf.Head().String() != "headA"+ln {
			t.Fail()
		}

		if buf.BodyAttr().String() != `bgcolor="blue"`+ln {
			t.Fail()
		}

		if buf.BodyHeader().String() != "bodyHeaderA"+ln {
			t.Fail()
		}

		if buf.BodyContent().String() != "bodyContentA"+ln {
			t.Fail()
		}

		if buf.BodyFooter().String() != "bodyFooterA"+ln {
			t.Fail()
		}

		if buf.BodyJs().String() != "bodyJsA"+ln {
			t.Fail()
		}
	})

	App.TestView = method

	ts := httptest.NewServer(App)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	Check(err)

	expectedOutput := `<!DOCTYPE html>
<html lang="en"
>
<head>
<title>titleA
</title>
headA

</head>
<body bgcolor="blue"
>
bodyHeaderA
bodyContentA
bodyFooterA
bodyJsA

</body>
</html>`

	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if string(b) != expectedOutput {
		t.Fail()
	}
}
