package configuration

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	e "github.com/liioan/faek/internal/errors"
	o "github.com/liioan/faek/internal/options"
)

const settingsFilePath = "/faek_settings.json"
const settingsDirectoryPath = "/.config/faek"

type Settings struct {
	OutputStyle string `json:"outputStyle"`
	FileName    string `json:"fileName"`
	Language    string `json:"lang"`
}

func getConfigDirectory() (string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return dirname + settingsDirectoryPath, nil
}
func getConfigFilePath() (string, error) {
	dirname, err := getConfigDirectory()
	if err != nil {
		return "", err
	}
	return dirname + settingsFilePath, nil
}

func SaveUserSettings(settings *Settings) error {
	if settings.FileName == "" {
		settings.FileName = "faekOutput.ts"
	}

	if settings.OutputStyle == "" {
		settings.OutputStyle = string(o.Terminal)
	}

	if settings.Language == "" {
		settings.Language = string(o.TypeScript)
	}

	settings.FileName = strings.Split(settings.FileName, ".")[0] + ".ts"

	bytes, err := json.Marshal(settings)
	if err != nil {
		return errors.New(e.CantMarshalJson)
	}

	filePath, err := getConfigDirectory()
	if err != nil {
		return errors.New(e.CantCreateConfigDirectory)
	}

	os.MkdirAll(filePath, 0755)
	filePath, err = getConfigFilePath()
	if err != nil {
		return errors.New(e.CantCreateConfigDirectory)
	}

	file, _ := os.Create(filePath)
	_, err = file.Write(bytes)

	if err != nil {
		return errors.New(e.CanSaveToFile)
	}
	return nil
}

func GetUserSettings() (Settings, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Settings{}, errors.New(e.FileDoesNotExists)
	}
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return Settings{}, errors.New(e.FileDoesNotExists)
	}
	s := Settings{}
	err = json.Unmarshal(fileBytes, &s)

	if err != nil {
		return Settings{}, errors.New(e.CantUnmarshalJson)
	}

	return s, nil
}
