package main

import (
	"fmt"
	"tfwf"
	//"os"
)

type index struct {
	tfwf.Context
}

func (i index) Get() {
	i.TemplateName = "index/home.html"
	a := map[string]string{"title": "timu", "date": "1900-2-1"}
	err := i.Render(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(i.Response, "Hello world!%s,%s", i.Args["second"], i.Args["secon"])
    tfwf.Logger.Error("dfdfd")
}

func (i index) Test(second string) {
	i.TemplateName = "index/home.html"
	a := map[string]string{"title": "timu", "date": "1900-2-1"}
	err := i.Render(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(i.Response, "00000Hello world!%s", second)
}

func main() {
	fmt.Printf("running\n")
	m := map[string]string{"host": "www.go.org", "name": "gogogo"}
	//设置在上边
	tfwf.AddSettings(m)
	tfwf.PrintSettings()
	//load 在下边
	err := tfwf.LoadTemplate()
	if err != nil {
		fmt.Printf("%s", err)
	}
	//tfwf.HandleFunc("^/d$", &index{}, "home")
	tfwf.HandleFunc(`^/aa/(?P<second>\d+)/(\d+)/(?P<secon>\d+)`, &index{}, "", "home")
	i := &index{}
	i.AllowMethod = []string{"POST", "GET"}
	tfwf.HandleFunc(`^/test/(?P<second>\d+)/$`, i, "Test", "home2")
	tfwf.ListenAndServe(":8080", nil)

}
