package main

import (
	"fmt"
	"log"

	"github.com/zokypesch/protoc-gen-react/lib"
	tpl "github.com/zokypesch/protoc-gen-react/template"
)

func main() {
	list := []lib.List{
		tpl.Basic,
		tpl.Wrapper,
		tpl.Service,
		tpl.ServiceMock,
		tpl.TypesRef,
		tpl.Types,
	}

	res, err := lib.NewMaster(list).Generate()

	if err != nil {
		log.Println(err)
	}

	for _, v := range res {
		pr := fmt.Sprintf("Execute %s has been successfull", v.Filename)
		log.Println(pr)
	}
}
