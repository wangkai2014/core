/*
A Basic Web Framework, the Core of Gorail.

It's built on top of the standard package 'net/http'.

Installation:

	go get github.com/gorail/core

HtmlFunc:
	cookie:	Get Cookie (input: string | output: *http.Cookie)
	css:	CSS Marksafe (input: string | output: template.CSS)
	html:	HTML Marksafe (input: string | output: template.HTML)
	htmlattr:	HTML Attribute Marksafe (input: string | output: template.HTMLAttr)
	js:	Javascript Marksafe (input: string | output: template.JS)
	jsstr:	Javascript String Marksafe (input: string | output: template.JSStr)
	url: URL Reverse (input: string, ...interface{} | output: string)
	time:	Convert to Default Timezone (input: time.Time | output: time.Time)
	timeZone: Convert to Timezone (input: string, time.Time | output: time.Time)
	timeFormat: Format time (input: string, time.Time | output: time.Time)
*/
package core

/*
	Assert Serial Numbers
	asn_Core_0001 : Method
	asn_Core_0002 : Protocol
*/
