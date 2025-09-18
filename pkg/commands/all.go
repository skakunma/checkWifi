package commands

type HandleFunction interface {
	Init() error
}

func HandleCommands(command string) HandleFunction {
	switch command {
	case "start":
		return CreateDefaultStartHandler()
	}
}
