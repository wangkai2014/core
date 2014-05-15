package core

import (
	"mime/multipart"
	"net/url"
)

type Value struct {
	Form          url.Values
	PostForm      url.Values
	MultipartForm url.Values
}

func (v Value) Get(key string) string {
	var value string

	if v.Form == nil {
		goto postForm
	}

	value = v.Form.Get(key)
	if value != "" {
		return value
	}

postForm:

	if v.Form == nil {
		goto multi
	}

	value = v.PostForm.Get(key)
	if value != "" {
		return value
	}

multi:

	if v.MultipartForm == nil {
		return value
	}

	return v.MultipartForm.Get(key)
}

// Get All
func (v Value) All(key string) []string {
	var strings []string

	if v.Form == nil {
		goto postForm
	}

	if len(v.Form[key]) > 0 {
		strings = v.Form[key]
		delete(v.Form, key)
	} else {
		goto postForm
	}

	return strings

postForm:

	if v.PostForm == nil {
		goto multi
	}

	if len(v.PostForm[key]) > 0 {
		strings = v.PostForm[key]
		delete(v.PostForm, key)
	} else {
		goto multi
	}

	return strings

multi:

	if v.MultipartForm == nil {
		return strings
	}

	if len(v.MultipartForm[key]) > 0 {
		strings = v.MultipartForm[key]
		delete(v.MultipartForm, key)
	}

	return strings
}

func (v Value) Shift(key string) string {
	var str string

	if v.Form == nil {
		goto postForm
	}

	if len(v.Form[key]) == 0 {
		goto postForm
	}

	if len(v.Form[key]) <= 1 {
		if len(v.Form[key]) == 1 {
			str = v.Form[key][0]
		}
		delete(v.Form, key)
	} else {
		str, v.Form[key] = v.Form[key][0], v.Form[key][1:]
	}

	return str

postForm:

	if v.PostForm == nil {
		goto multi
	}

	if len(v.PostForm[key]) == 0 {
		goto multi
	}

	if len(v.PostForm[key]) <= 1 {
		if len(v.PostForm[key]) == 1 {
			str = v.PostForm[key][0]
		}
		delete(v.PostForm, key)
	} else {
		str, v.PostForm[key] = v.PostForm[key][0], v.PostForm[key][1:]
	}

	return str

multi:

	if v.MultipartForm == nil {
		return str
	}

	if len(v.MultipartForm[key]) <= 1 {
		if len(v.MultipartForm[key]) == 1 {
			str = v.MultipartForm[key][0]
		}
		delete(v.MultipartForm, key)
	} else {
		str, v.MultipartForm[key] = v.MultipartForm[key][0], v.MultipartForm[key][1:]
	}

	return str
}

type Form struct {
	Value Value
	File  map[string][]*multipart.FileHeader
}

func (f *Form) GetFile(key string) *multipart.FileHeader {
	if f.File == nil {
		return nil
	}
	if len(f.File[key]) == 0 {
		return nil
	}
	return f.File[key][0]
}

func (f *Form) ShiftFile(key string) *multipart.FileHeader {
	if f.File == nil {
		return nil
	}

	var file *multipart.FileHeader

	if len(f.File[key]) <= 1 {
		if len(f.File[key]) == 1 {
			file = f.File[key][0]
		}
		delete(f.File, key)
	} else {
		file, f.File[key] = f.File[key][0], f.File[key][1:]
	}

	return file
}

// Generate a new form.
func (c *Context) Form() *Form {
	if c.pri.form != nil {
		return c.pri.form
	}

	c.Req.ParseMultipartForm(c.App.FormMemoryLimit)
	c.Req.ParseForm()

	form := &Form{Value: Value{}}

	form.Value.Form = c.Req.Form
	form.Value.PostForm = c.Req.PostForm

	if c.Req.MultipartForm != nil {
		form.Value.MultipartForm = url.Values(c.Req.MultipartForm.Value)
		form.File = c.Req.MultipartForm.File
	}

	c.pri.form = form

	return form
}
