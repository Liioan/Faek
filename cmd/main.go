package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liioan/faek/internal/configuration"
	"github.com/liioan/faek/internal/model"
	o "github.com/liioan/faek/internal/options"
	"github.com/liioan/faek/internal/utils"
)

func main() {
	utils.ClearConsole()

	var configMode bool
	flag.BoolVar(&configMode, "c", false, "enter configuration mode")

	var debugMode bool
	flag.BoolVar(&debugMode, "d", false, "enter debug mode")
	flag.Parse()

	if debugMode {
		settings, err := configuration.GetUserSettings()
		if err != nil {
			fmt.Println("Fatal: ", err)
			os.Exit(1)
		}
		fmt.Println(settings)
		return
	}

	_, err := configuration.GetUserSettings()

	if err != nil {
		configMode = true
	}

	var steps []model.Step

	if configMode {
		steps = []model.Step{
			*model.NewListStep("Choose your default output style", "Output options:", false, o.OutputSet),
			*model.NewListStep("Choose your preferred language (default: TypeScript)", "Language options:", false, o.LanguageSet),
			*model.NewTextStep("Choose filename for output file (default: faekOutput.ts)", "e.g. output.ts", false),
		}

	} else {
		steps = []model.Step{
			*model.NewTextStep("What will the array be called? (default: arr)", "e.g. users", false),
			*model.NewTextStep("Write your field (to continue press enter without input)", "e.g. name string", true),
			*model.NewTextStep("Create type for your object? (default: no type, input: type name)", "e.g. Post", false),
			*model.NewTextStep("How many items will be in this array (default 5)", "e.g. 5", false),
		}
	}

	model := model.NewModel(steps, configMode)

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
