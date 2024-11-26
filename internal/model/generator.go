package model

import "github.com/liioan/faek/internal/utils"

func generateOutput(m Model) {
	output := ""
	test := m.Steps
	for _, step := range test {
		output += "\n"
		if len(step.Answer.fields) > 0 {
			output += "fields: \n"
			for _, field := range step.Answer.fields {
				output += field.name + " " + field.fieldType
				if len(field.variant) > 0 {
					output += " " + string(field.variant)
				}

				if len(field.options) > 0 {
					for _, option := range field.options {
						output += " " + option
					}
				}

				output += "\n"
			}
			output += "\n"
			continue
		}
		output += step.Instruction + ": \n"
		output += step.Answer.text + "\n"
	}

	utils.LogToDebug(output)
}
