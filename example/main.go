package main

import (
	"fmt"

	"github.com/mdd13/commentui/component"
)

// commentui:component
// _file:root.html
func Root() string {
	m := component.Metadatas{
		"C": component.NewMetadataString("HELLO"),
		"XX": component.NewMetadataString("ZZZ"),
	}
	return component.PutData(m)
}

// commentui:component
// <div class="_str:Title">_str:Title</div>
func Header() string {
	m := component.Metadatas{
		"Title": component.NewMetadataString("HOLA"),
	}
	return component.PutData(m)
}

// commentui:component
// <div>Xin chao the gioi</div>
func Body() string {
 	return component.PutData(nil)
}

func main() {
	component.InitConfig(
		"./",
		"./html",
	)

	root := Root()
	Header()
	Body()

	fmt.Println(component.RenderComponent(root))
}
