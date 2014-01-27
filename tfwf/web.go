package tfwf

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

type Context struct {
	Request      *http.Request
	Response     http.ResponseWriter
	Args         map[string]string      // url中的参数
	TemplateName string                 // 模板名称
	TemplateArgs map[string]interface{} // 模板参数
	AllowMethod  []string
}

var SUPPORTED_METHODS = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}

func (c *Context) Init() {
	if len(c.AllowMethod) == 0 {
		c.AllowMethod = SUPPORTED_METHODS
	}
	if !inSlice(c.Request.Method, c.AllowMethod) {
		c.HttpMethodNotAllowed("method not allowed", "")
		return
	}
	c.TemplateArgs = make(map[string]interface{})
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

func (c *Context) Render() {
	err := c.ExecTemplate(c.TemplateArgs)
	if err != nil {
		panic(err)
	}
	return
}

func (c *Context) Write(i interface{}) {
	err := c.ExecTemplate(i)
	if err != nil {
		panic(err)
	}
	return
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
}

//from request
func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	cookie, err := c.Request.Cookie(name)
	return cookie, err
}
func (c *Context) ClearCookie(name string) (err error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return
	}
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	return
}
func (c *Context) ClearAllCookies() {
	for _, cookie := range c.Request.Cookies() {
		cookie.MaxAge = -1
		c.SetCookie(cookie)
	}
}

func (c *Context) SetSecureCookie(cookie *http.Cookie) {
	cookie.Value = create_signed_value(cookie.Name, cookie.Value)
	c.SetCookie(cookie)
}
func (c *Context) GetSecureCookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	v := strings.Split(cookie.Value, "|")
	if len(v) != 3 {
		return "", errors.New("Invalid cookie value!")
	}
	signature := cookie_signature(name, v[0], v[1])
	if v[2] != signature {
		msg := "Invalid cookie signature!"
		Logger.Error(msg)
		return "", errors.New(msg)
	}
	value, err := base64.StdEncoding.DecodeString(v[0])
	if err != nil {
		return "", err
	}
	return string(value), nil
}

type HttpError struct {
	Code         int    //http状态码
	Desc         string // 返回的文字
	TemplateName string //返回的template
}

//405
func (c *Context) HttpMethodNotAllowed(desc, template string) {
	panic(HttpError{http.StatusMethodNotAllowed, desc, template})
}

//404
func (c *Context) HttpNotFound(desc, template string) {
	panic(HttpError{http.StatusNotFound, desc, template})
}
