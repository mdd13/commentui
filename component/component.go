package component

import (
	"os"
	"runtime"
	"strings"
	"unsafe"
)

type Config struct {
	goPath   string
	htmlPath string
}

var globalConfig Config

func InitConfig(goPath, htmlPath string) {
	globalConfig.goPath = goPath
	globalConfig.htmlPath = htmlPath
}

var globalTemplates = make(map[string]string)
var globalComponents = make(map[string]string)
var globalMetadatas = make(map[string]Metadatas)

func funcName(s string) string {
	idx := strings.LastIndexByte(s, '.')
	return s[idx+1:]
}

func readFile(path string) string {
	buf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return *(*string)(unsafe.Pointer(&buf))
}

func readFileLines(path string) []string {
	return strings.Split(readFile(path), "\n")
}

const cmtPrefix = "//"
const beginComponent = "commentui:component"

func extractComment(lines []string) string {
	n := len(lines) - 1
	
	i := n
	for ; i > 0; i-- {
		if !strings.HasPrefix(lines[i], cmtPrefix) {
			continue
		}

		v := strings.TrimSpace(lines[i][len(cmtPrefix):])
		if v == beginComponent {
			i++
			break
		}
	}

	sb := strings.Builder{}
	for ; i < n; i++ {
		if !strings.HasPrefix(lines[i], cmtPrefix) {
			continue
		}
		v := strings.TrimSpace(lines[i][len(cmtPrefix):])
		sb.WriteString(v)
		sb.WriteString("\n")
	}

	result := sb.String()
	if len(result) == 0 {
		return ""
	}
	return result[:len(result)-1]
}

var GlobalTpl = make(map[string]string)

type Metadatas = map[string]*Metadata

func PutData(m Metadatas) string {
	pc, file, n, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)

	lines := readFileLines(file)
	tpl := extractComment(lines[:n])

	fnName := funcName(fn.Name())
	globalTemplates[fnName] = tpl
	globalMetadatas[fnName] = m
	return fnName
}

func RenderComponent(k string) string {
	tpl := globalTemplates[k]
	mds := globalMetadatas[k]
	result := TemplateRender(tpl, mds)
	globalComponents[k] = result
	return result
}

func RenderAll() map[string]string {
	for k := range globalTemplates {
		RenderComponent(k)
	}
	return globalComponents
}
