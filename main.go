package main

import (
	"fmt"
	"os"

	"github.com/ikasoba/voich/pipewire"
	"github.com/k0kubun/pp"
)

func main() {
	fmt.Println("Hello, World!")

	pw, err := pipewire.LoadPipewire()
	if err != nil {
		panic(err)
	}

	pp.Println("version:", pw.GetLibraryVersion())

	pp.Println("init:", pw.Init(os.Args))
	pp.Println("application name:", pw.GetApplicationName())
	pp.Println("program name:", pw.GetPrgName())
	pp.Println("user name:", pw.GetUserName())
	pp.Println("host name:", pw.GetHostName())
	pp.Println("client name:", pw.GetClientName())

	fmt.Printf(
		"\npw.DirectionReverse(SpaDirectionInput) = %s\n\n",
		pw.DirectionReverse(pipewire.SpaDirectionInput),
	)

	pp.Println("domain:", pw.GetDomain())

	supports := make([]pipewire.SpaSupport, 64)

	supports = supports[:pw.GetSupport(supports)]
	for i, s := range supports {
		pp.Println(i, s)
	}

	defer pw.Deinit()
}
