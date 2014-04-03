package core

import (
	"mime/multipart"
	"net/url"
)

type Form struct {
	Value url.Values
	File  map[string][]*multipart.FileHeader
}

// Generate a new form with the key prefix trimed out, if the key does not have the prefix, it will get ignore.
func (f *Form) TrimPrefix(prefix string) *Form {
	form := &Form{Value: url.Values{}}
	prefixLen := len(prefix)
	for key, value := range f.Value {
		if len(key) <= prefixLen {
			continue
		}
		if prefix != key[:prefixLen] {
			continue
		}
		form.Value[key[prefixLen:]] = value
	}
	if f.File == nil {
		return form
	}
	form.File = map[string][]*multipart.FileHeader{}
	for key, value := range f.File {
		if len(key) <= prefixLen {
			continue
		}
		if prefix != key[:prefixLen] {
			continue
		}
		form.File[key[prefixLen:]] = value
	}
	return form
}

// Generate a new form, retaining the chosen slot number!
func (f *Form) Slot(slot int) *Form {
	form := &Form{Value: url.Values{}}
	for key, value := range f.Value {
		if len(value) > slot {
			form.Value[key] = append(form.Value[key], value[slot])
		}
	}
	if f.File == nil {
		return form
	}
	form.File = map[string][]*multipart.FileHeader{}
	for key, value := range f.File {
		if len(value) > slot {
			form.File[key] = append(form.File[key], value[slot])
		}
	}
	return form
}

func inSlice(needle string, haystack []string) bool {
	if haystack == nil {
		return false
	}
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}
	return false
}

// Generate a new form.  With only the allowed fields.
func (f *Form) Allow(fields ...string) *Form {
	form := &Form{Value: url.Values{}}
	for key, value := range f.Value {
		if inSlice(key, fields) {
			form.Value[key] = value
		}
	}
	if f.File == nil {
		return form
	}
	form.File = map[string][]*multipart.FileHeader{}
	for key, value := range f.File {
		if inSlice(key, fields) {
			form.File[key] = value
		}
	}
	return form
}

// Generate a new form, while filtering out the denied fields.
func (f *Form) Deny(fields ...string) *Form {
	form := &Form{Value: url.Values{}}
	for key, value := range f.Value {
		if inSlice(key, fields) {
			continue
		}
		form.Value[key] = value
	}
	if f.File == nil {
		return form
	}
	form.File = map[string][]*multipart.FileHeader{}
	for key, value := range f.File {
		if inSlice(key, fields) {
			continue
		}
		form.File[key] = value
	}
	return form
}

// Generate a new form.
func (c *Context) Form() *Form {
	c.Req.ParseMultipartForm(c.App.FormMemoryLimit)

	form := &Form{}
	if c.Req.MultipartForm != nil {
		form.Value = url.Values(c.Req.MultipartForm.Value)
		form.File = c.Req.MultipartForm.File
	} else if c.Req.PostForm != nil {
		form.Value = c.Req.PostForm
		form.File = nil
	} else {
		form.Value = c.Req.Form
		form.File = nil
	}
	return form
}
