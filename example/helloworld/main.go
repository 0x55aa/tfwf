package main

import (
	"fmt"
	"tfwf"
)

type index struct {
	tfwf.Context
}

func (i index) Get() {
	i.TemplateName = "index/home.html"
	i.TemplateArgs["title"] = "hello"
	i.TemplateArgs["name"] = i.Args["name"]
	i.Render()
	fmt.Fprintf(i.Response, "end")
}

func (i index) Hello(name string) {
	i.TemplateName = "index/home.html"
	i.TemplateArgs["title"] = "hello"
	i.TemplateArgs["name"] = name
	i.Render()
}

func main() {
	fmt.Printf("running\n")
	s := map[string]string{"name": "hello world"}
	tfwf.AddSettings(s)

	tfwf.HandleFunc(`^/(?P<name>\w+)/$`, &index{}, "", "home")

	i := &index{}
	i.AllowMethod = []string{"GET"}
	tfwf.HandleFunc(`^/hello/(?P<name>\w+)/$`, i, "Hello", "hello")

	tfwf.ListenAndServe(":8080", nil)
}
