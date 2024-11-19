package model

import "github.com/liioan/faek/internal/utils"

func saveSettings(m Model) {
	output := ""
	test := m.Steps
	for _, step := range test {
		output += "\n"
		output += step.Answer.text
	}
	utils.LogToDebug(output)

}
