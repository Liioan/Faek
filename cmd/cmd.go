package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/liioan/faek/internal/configuration"
	"github.com/liioan/faek/internal/help"
	m "github.com/liioan/faek/internal/model"
	"github.com/liioan/faek/internal/styles"
	"github.com/liioan/faek/internal/utils"
	v "github.com/liioan/faek/internal/variants"
)

type RuntimeFlags struct {
	helpMode   bool
	configMode bool
	debugMode  bool
	template   string
	length     int

	language v.Variant
	output   v.Variant
	export   v.Variant
}

func Execute() {
	utils.ClearConsole()

	flags := parseFlags()

	if flags.helpMode {
		help.ShowHelpScreen()
		return
	}

	// enter configuration mode if config is not found
	_, err := configuration.GetUserSettings()
	if err != nil {
		flags.configMode = true
	}

	var steps []m.Step
	if flags.configMode {
		steps = []m.Step{
			*m.CreateListStep("Choose your default output style", "Output options:", v.OutputVariants),
			*m.CreateListStep("Choose your preferred language", "Language options:", v.LanguageVariants),
			*m.CreateTextStep("Choose filename for output file (default: faekOutput.ts)", "e.g. output.ts"),
			*m.CreateTextStep("Choose indent size (default: 2)", "e.g. 4"),
			*m.CreateListStep("Choose preferred export", "Export options:", v.ExportVariants),
		}
	} else {
		propStep := m.CreatePropsStep()
		steps = []m.Step{
			*m.CreateTextStep("What will the array be called? (default: arr)", "e.g. users"),
			*propStep,
			*m.CreateEditStep(propStep),
			*m.CreateOptionalStep("Type", "Create type for your object?", []string{"no", "yes"}, "Write your type name", "e.g. Post"),
			*m.CreateTextStep("How many items will be in this array (default 5)", "e.g. 5"),
		}

	}

	overrideFlags := m.Override{Language: flags.language, Output: flags.output, Export: flags.export}

	model, err := m.NewModel(steps, flags.configMode, overrideFlags)
	if err != nil {
		log.Fatal(err)
	}

	if flags.debugMode {
		text := styles.TitleStyle.Render("----- Debug mode -----\n")

		settings, err := configuration.GetUserSettings()
		if err != nil {
			fmt.Println("Fatal: ", err)
			os.Exit(1)
		}
		text += styles.OutputStyle.Render(fmt.Sprintf("user settings:\n%v\n\n", settings))

		text += styles.OutputStyle.Render(fmt.Sprintf("runtime flags:\n%v\n\n", flags))

		fmt.Print(text)
		model = m.NewDebugModel(steps, flags.template, flags.length, overrideFlags)
		file, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("Fatal:", err)
			os.Exit(1)
		}
		defer file.Close()
	}

	program := tea.NewProgram(*model)
	_, err = program.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags() RuntimeFlags {
	flags := RuntimeFlags{}

	//  help mode
	flag.BoolVar(&flags.helpMode, "h", false, "display help")

	//  config mode
	flag.BoolVar(&flags.configMode, "c", false, "enter configuration mode")

	//  config mode
	flag.BoolVar(&flags.debugMode, "d", false, "enter debug mode")
	flag.StringVar(&flags.template, "template", "types", "create types template")
	flag.IntVar(&flags.length, "len", 5, "add length")

	// language override
	var tsFlag bool
	var jsFlag bool
	var jsonFlag bool
	flag.BoolVar(&tsFlag, "ts", false, "overrides configuration - changes language to ts")
	flag.BoolVar(&jsFlag, "js", false, "overrides configuration - changes language to js")
	flag.BoolVar(&jsonFlag, "json", false, "overrides configuration - changes language to json")

	// output override
	var fileFlag bool
	var terminalFlag bool
	flag.BoolVar(&fileFlag, "file", false, "overrides configuration - changes output to file")
	flag.BoolVar(&terminalFlag, "terminal", false, "overrides configuration - changes output to terminal")

	flag.BoolFunc("export", "overrides configuration - sets export", func(s string) error {
		if len(s) == 0 {
			flags.export = v.Config
		} else {
			flags.export = v.Variant(s)
		}

		return nil
	})

	flag.Parse()

	flags.language = getLangOverride(tsFlag, jsFlag, jsonFlag)
	flags.output = getOutputOverride(fileFlag, terminalFlag)

	return flags
}

func getOutputOverride(fileFlag, terminalFlag bool) v.Variant {
	if fileFlag {
		return v.File
	} else if terminalFlag {
		return v.Terminal
	}
	return v.Config
}

func getLangOverride(tsFlag, jsFlag, jsonFlag bool) v.Variant {
	if tsFlag {
		return v.TypeScript
	} else if jsFlag {
		return v.JavaScript
	} else if jsonFlag {
		return v.JSON
	}
	return v.Config
}
