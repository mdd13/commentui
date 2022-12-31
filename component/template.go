package component

import (
	"math"
	"path/filepath"
	"sort"
	"strings"
)

type TemplateEntityType int

const (
	TplUnknown TemplateEntityType = iota
	TplLoopMap
	TplLoopString
	TplLoopEntity
	TplString
	TplComponent
	TplFile
)

type TemplateEntity struct {
	vt TemplateEntityType
	vn string
	dt string
}

var (
	TplPrefFile = "_file"
	TplPrefString = "_str"
	TplPrefComponent = "_cpn"
	TplPrefLoopString = "_loop_str"
	TplPrefLoopStringEnd = "end_loop_str"
	TplPrefLoopMap = "_loop_map"
	TplPrefLoopMapEnd = "_end_loop_map"
	TplPrefLoopElement = "_loop_element"
)

var Prefs = []*string {
	&TplPrefFile,
	&TplPrefString,
	&TplPrefComponent,
	&TplPrefLoopString,
	&TplPrefLoopMap,
}

func templateToTheEnd(s string) int {
	for i, c := range s {
		if c == '"' || c == '\'' || c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			return i
		}
	}

	return len(s)
}

func templateToTheEndBlock(s string, t string) int {
	return strings.Index(s, "_end" + t)
}

func templateCountEntity(s string) int {
	result := 0
	for _, p := range Prefs {
		result += strings.Count(s, *p)
	}
	return result
}

func templateFindNextEntity(s string) (*TemplateEntity, int) {
	idx := math.MaxInt32
	typ := ""
	for _, p := range Prefs {
		t := strings.Index(s, *p)
		if t >= 0 && t < idx {
			idx = t
			typ = *p
		}
	}

	s = s[idx:]
	
	switch typ {
	case TplPrefFile:
		end := templateToTheEnd(s)
		dt := s[:end]
		tpl := &TemplateEntity{
			dt: dt,
			vn: templateVarName(dt),
			vt: TplFile,
		}
		return tpl, end + idx
	case TplPrefString:
		end := templateToTheEnd(s)
		dt := s[:end]
		tpl := &TemplateEntity{
			dt: dt,
			vn: templateVarName(dt),
			vt: TplString,
		}
		return tpl, end + idx
	case TplPrefComponent:
		end := templateToTheEnd(s)
		dt := s[:end]
		tpl := &TemplateEntity{
			dt: dt,
			vn: templateVarName(dt),
			vt: TplComponent,
		}
		return tpl, end + idx
	case TplPrefLoopString:
		end := templateToTheEndBlock(s, TplPrefLoopString)
		if end > 0 {
			dt := s[:end]
			tpl := &TemplateEntity{
				dt: dt,
				vn: templateVarName(dt),
				vt: TplLoopString,
			}
			return tpl, end + idx
		}
		panic("Not found end character of " + TplPrefLoopString)
	case TplPrefLoopMap:
		end := templateToTheEndBlock(s, TplPrefLoopMap)
		if end > 0 {
			dt := s[:end]
			tpl := &TemplateEntity{
				dt: dt,
				vn: templateVarName(dt),
				vt: TplLoopMap,
			}
			return tpl, end + idx
		}
		panic("Not found end character of " + TplPrefLoopMap)
	}
	return nil, -1
}

func templateFindEntities(s string) []TemplateEntity {
	count := templateCountEntity(s)
	result := make([]TemplateEntity, count)

	for i := 0; i < count; i++ {
		entity, end := templateFindNextEntity(s)
		result[i] = *entity
		s = s[end:]
	}

	return result
}

func templateVarName(s string) string {
	s = s[strings.IndexByte(s, ':')+1:]
	return strings.TrimSpace(s)
}

func htmlPath(path string) string {
	return filepath.Join(globalConfig.htmlPath, path)
}

func TemplateRender(tpl string, m map[string]*Metadata) string {
	ents := templateFindEntities(tpl)
	
	sort.SliceStable(ents, func(i, j int) bool {
		return ents[i].vt < ents[j].vt
	})

	for _, e := range ents {
		switch e.vt {
		case TplString:
			metadata := m[e.vn]
			if metadata != nil {
				tpl = strings.ReplaceAll(tpl, e.dt, MetadataString(metadata))
			}
		case TplFile:
			file := e.vn
			tpl = strings.ReplaceAll(tpl, e.dt, readFile(htmlPath(file)))
			tpl = TemplateRender(tpl, m)
		case TplComponent:
			component := globalComponents[e.vn]
			if len(component) == 0 {
				RenderComponent(e.vn)
				component = globalComponents[e.vn]
			}
			tpl = strings.ReplaceAll(tpl, e.dt, component)
		case TplUnknown:
			tpl = strings.ReplaceAll(tpl, e.dt, "")
		}
	}

	return tpl
}
