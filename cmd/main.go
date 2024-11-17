package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liioan/faek/internal/model"
	"github.com/liioan/faek/internal/utils"
)

func main() {
	utils.ClearConsole()

	var config bool
	flag.BoolVar(&config, "c", false, "enter configuration mode")
	flag.Parse()

	steps := []model.Step{

	}

	if config {
		steps = []model.Step{
			*model.NewListStep("Choose your default output style", false, {});
		}
	} else {
	steps = []model.Step{
		*model.NewTextStep("What will the array be called? (default: arr)", "e.g. users", false),
		*model.NewTextStep("Write your field (to continue press enter without input)", "e.g. name string", true),
		*model.NewTextStep("Create type for your object? (default: no type, input: type name)", "e.g. Post", false),
		*model.NewTextStep("How many items will be in this array (default 5)", "e.g. 5", false),
	}
	}


	model := model.NewModel(steps, false)

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
