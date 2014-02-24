package tfwf

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

var Templates map[string]*template.Template

func init() {
	Templates = make(map[string]*template.Template)
}

//抛出异常结束执行 0 0
func loadTemplate(path string, f os.FileInfo, err error) error {
	if f == nil {
		return err
	}
	if f.IsDir() {
		return nil
	}
	keyPath := path[len(Settings["template_dir"])+1:]
	t := template.Must(template.ParseFiles(path))
	Templates[keyPath] = t
	return nil
}

func LoadTemplate() error {
	//先遍历取到所有文件名（包含路径）
	err := filepath.Walk(Settings["template_dir"], loadTemplate)
	if err != nil {
		fmt.Printf("load template error:%v", err)
		return err
	}
	return nil
}
