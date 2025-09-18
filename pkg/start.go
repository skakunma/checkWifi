package pkg

import (
	"checkWifi/pkg/commands"
	"fmt"
)

func InitApplication() {
	LoadInscription()
	for {
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			panic(err)
		}
		f := commands.HandleCommands(input)
		if err = f.Init(); err != nil {
			panic(err)
		}
	}
}
