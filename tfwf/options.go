package tfwf

import (
	"fmt"
)

//后面添加文件配置的读取

var settings map[string]string

//初始化Settings
// template_dir 模板目录 最后不加斜杠
// static_dir 静态文件目录
// secret_key 密钥，加密用的0 0
func init() {
	settings = map[string]string{"template_dir": "templates", "static_dir": "static"}
}

//添加的接口
func AddSettings(m map[string]string) {
	for k, v := range m {
		settings[k] = v
	}
}

//test
func PrintSettings() {
	for k, v := range settings {
		fmt.Printf("%s: %s\n", k, v)
	}
}
