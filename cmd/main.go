package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/liioan/faek/internal/generator"
	"github.com/liioan/faek/internal/utils"
)

func main() {
	utils.ClearConsole()
	steps := []generator.Step{
		*generator.NewTextStep("Array name:", "e.g. users", false),
		*generator.NewTextStep("fields", "e.g. name string", true),
		*generator.NewTextStep("len", "e.g. 5", false),
	}

	model := generator.NewModel(steps, false)

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
