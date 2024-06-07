package main

import (
	"faek/internal/constance"
	"faek/internal/generator"
	"faek/internal/types"
	"faek/internal/utils"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

//+ Ideas:
//. add date type

func main() {
	utils.ClearConsole()
	steps := []types.Step{
		{Instruction: "what will the array be called? (default: arr)", Placeholder: "E.g. userPosts"},
		{Instruction: "write your field (to continue press enter without input)", IsRepeating: true, Placeholder: "E.g. email string"},
		{Instruction: "Create type for your object? (default: no type, input: type name)", Placeholder: "E.g. Post"},
		{Instruction: "how many items will be in this array (default 5)", Placeholder: "E.g. 5"},
	}

	model := generator.New(steps)

	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer file.Close()

	if utils.FileExists(constance.OUTPUT_FILEPATH) {
		err = os.Remove(constance.OUTPUT_FILEPATH)
	}
	if err != nil {
		log.Fatal(err)
	}

	program := tea.NewProgram(*model)
	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}
}
