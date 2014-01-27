package tfwf

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type ServeMux struct {
	routes         []*route
	named_handlers map[string]*route
}

type route struct {
	name    string
	regex   *regexp.Regexp
	handler HandlerInterface
	f       string // 方法
}

type HandlerInterface interface {
	Initialize()
	Prepare()
	Finish()
}

func NewServerMux() *ServeMux {
	return &ServeMux{
		named_handlers: make(map[string]*route),
	}
}

var DefaultServeMux = NewServerMux()

//先这样写，还没想好接口的样子
func HandleFunc(r string, handler HandlerInterface, f string, name string) {
	regex, err := regexp.Compile(r)
	if err != nil {
		fmt.Printf("Error in route regex %q\n", r)
		return
	}
	DefaultServeMux.HandleFunc(regex, handler, f, name)
}

func Error(w http.ResponseWriter, http_error HttpError) (err error) {

	w.WriteHeader(http_error.Code)
	if http_error.TemplateName != "" {
		if v, ok := Templates[http_error.TemplateName]; ok {
			err = v.Execute(w, "")
		} else {
			err = errors.New("Template " + http_error.TemplateName + " Not Found")
		}
		return
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintln(w, http_error.Desc)
	}
	return nil
}

func (mux *ServeMux) HandleFunc(regex *regexp.Regexp, handler HandlerInterface, f string, name string) {
	r := route{name: name, regex: regex, handler: handler, f: f}
	mux.routes = append(mux.routes, &r)
	if _, ok := mux.named_handlers[name]; ok {
		panic("handler had name " + name)
	} else {
		mux.named_handlers[name] = &r
	}
}

//进行匹配，返回handler
func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//用来判定是否返回panic，做到提前return
	//重写http.Error
	defer func() {
		if x := recover(); x != nil {
			switch x.(type) {
			case HttpError:
				httpError := reflect.ValueOf(x)
				code := int(httpError.FieldByName("Code").Int())
				desc := httpError.FieldByName("Desc").String()
				templateName := httpError.FieldByName("TemplateName").String()
				http_error := HttpError{code, desc, templateName}
				err := Error(w, http_error)
				if err != nil {
					panic(err)
				}
				return
			default:
				panic(x)
			}
		}
	}()

	//request_middleware
	h, f, name := mux.Handler(r)
	//没有匹配到,返回404
	if h == nil {
		http.NotFound(w, r)
		return
	}
	h_value := reflect.ValueOf(h)
	var h_value_ptr reflect.Value
	if !h_value.CanSet() {
		h_value_ptr = h_value.Elem()
	}
	h_value_ptr.FieldByName("Request").Set(reflect.ValueOf(r))
	h_value_ptr.FieldByName("Response").Set(reflect.ValueOf(w))

	in := make([]reflect.Value, 0)
	init := h_value.MethodByName("Init")
	init.Call(in)
	initialize := h_value.MethodByName("Initialize")
	initialize.Call(in)
	prepare := h_value.MethodByName("Prepare")
	prepare.Call(in)

	//添加取不到方法的判断
	//有方法用提供的
	if f != "" {
		method := h_value.MethodByName(f)
		if !method.IsValid() {
			err := "no method " + f
			Logger.Error(err)
			panic(err)
		}
		h_value_ptr.FieldByName("Args").Set(reflect.ValueOf(name))
		var inn []reflect.Value
		for _, v := range name {
			inn = append(inn, reflect.ValueOf(v))
		}
		method.Call(inn)
	} else {
		//转成合适的名字格式
		methodName := r.Method[0:1] + strings.ToLower(r.Method[1:])
		method := h_value.MethodByName(methodName)
		if !method.IsValid() {
			err := "no method " + methodName
			Logger.Error(err)
			panic(HttpError{http.StatusMethodNotAllowed, "method no allowed!", ""})
		}
		h_value_ptr.FieldByName("Args").Set(reflect.ValueOf(name))
		method.Call(in)
	}

	finish := h_value.MethodByName("Finish")
	finish.Call(in)
}

func (mux *ServeMux) Handler(r *http.Request) (h HandlerInterface, f string, name map[string]string) {
	return mux.handler(r.URL.Path)
}

func (mux *ServeMux) handler(path string) (h HandlerInterface, f string, name map[string]string) {
	h, f, name = mux.match(path)
	return
}

func (mux *ServeMux) match(path string) (HandlerInterface, string, map[string]string) {
	for _, r := range mux.routes {
		b, name := pathMatch(r, path)
		if !b {
			continue
		}
		h := r.handler
		f := r.f
		return h, f, name
	}
	return nil, "", nil
}

func (mux *ServeMux) Reverse() {
}

//看路由和url是否匹配
//后面改成返回匹配的串
func pathMatch(route *route, path string) (bool, map[string]string) {
	if !route.regex.MatchString(path) {
		return false, nil
	}
	match := route.regex.FindStringSubmatch(path)
	args := make(map[string]string)
	for i, name := range route.regex.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		args[name] = match[i]
	}
	return true, args
}

func ListenAndServe(addr string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: DefaultServeMux}
	if handler != nil {
		server.Handler = handler
	}
	return server.ListenAndServe()
}
