package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liioan/faek/internal/configuration"
	m "github.com/liioan/faek/internal/model"
	"github.com/liioan/faek/internal/utils"
	v "github.com/liioan/faek/internal/variants"
)

func main() {
	utils.ClearConsole()

	var configMode bool
	flag.BoolVar(&configMode, "c", false, "enter configuration mode")

	var debugMode bool
	var userTemplate bool
	var typesTemplate bool
	var length int

	flag.BoolVar(&debugMode, "d", false, "enter debug mode")
	flag.BoolVar(&userTemplate, "user", false, "create user template")
	flag.BoolVar(&typesTemplate, "types", true, "create types template")
	flag.IntVar(&length, "lenght", 5, "add length")
	flag.Parse()

	if debugMode {
		settings, err := configuration.GetUserSettings()
		if err != nil {
			fmt.Println("Fatal: ", err)
			os.Exit(1)
		}
		fmt.Println(settings)
	}

	_, err := configuration.GetUserSettings()

	if err != nil {
		configMode = true
	}

	var steps []m.Step

	if configMode {
		steps = []m.Step{
			*m.NewListStep("Choose your default output style", "Output options:", false, v.OutputSet),
			*m.NewListStep("Choose your preferred language (default: TypeScript)", "Language options:", false, v.LanguageSet),
			*m.NewTextStep("Choose filename for output file (default: faekOutput.ts)", "e.g. output.ts", false),
		}

	} else {
		steps = []m.Step{
			*m.NewTextStep("What will the array be called? (default: arr)", "e.g. users", false),
			*m.NewTextStep("Write your field (to continue press enter without input)", "e.g. name string", true),
			*m.NewTextStep("Create type for your object? (default: no type, input: type name)", "e.g. Post", false),
			*m.NewTextStep("How many items will be in this array (default 5)", "e.g. 5", false),
		}
	}

	model := m.NewModel(steps, configMode)

	if debugMode {
		model = m.NewDebugModel(steps, userTemplate, typesTemplate, length)
	}

	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("Fatal:", err)
		os.Exit(1)
	}
	defer file.Close()

	program := tea.NewProgram(*model)
	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
