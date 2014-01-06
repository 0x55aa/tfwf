package tfwf

import (
	"errors"
	"net/http"
)

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Args     map[string]string
	TemplateName string
	AllowMethod []string
}

var SUPPORTED_METHODS = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

func (c *Context) Init(){
	if len(c.AllowMethod) == 0 {
		c.AllowMethod = SUPPORTED_METHODS
	}
	if !inSlice(c.Request.Method, c.AllowMethod) {
		c.HttpMethodNotAllowed("method not allowed", "")
		return
	}
}

func (c *Context) Initialize() {

}
func (c *Context) Prepare() {

}
func (c *Context) Finish() {
}

func (c *Context) ExecTemplate(i interface{}) (err error) {
	if v, ok := Templates[c.TemplateName]; ok {
		err = v.Execute(c.Response, i)
	} else {
		err = errors.New("Template " + c.TemplateName + " Not Found")
	}
	return

}

func (c *Context) Render(i interface{}) (err error) {
	c.ExecTemplate(i)
	return
}

func (c *Context) Write(i interface{}) (err error) {
	c.ExecTemplate(i)
	return
}

func (c *Context) SetCookie() {
}

type HttpError struct {
	Code         int
	Desc         string
	TemplateName string
}

//405
func (c *Context) HttpMethodNotAllowed(desc, template string) {
	panic(HttpError{http.StatusMethodNotAllowed, desc, template})
}

//404
func (c *Context) HttpNotFound(desc, template string) {
	panic(HttpError{http.StatusNotFound, desc, template})
}
