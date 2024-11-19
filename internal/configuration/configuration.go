package configuration

import (
	"encoding/json"
	"errors"
	"os"

	e "github.com/liioan/faek/internal/errors"
)

const settingsFilePath = "/faek_settings.json"
const settingsDirectoryPath = "/.config/faek"

type Settings struct {
	OutputStyle string `json:"outputStyle"`
	FileName    string `json:"fileName"`
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

func SaveUserSettings(outputStyle, fileName string) error {
	settings := Settings{OutputStyle: outputStyle, FileName: fileName}

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
	defaultSettings := Settings{
		OutputStyle: "terminal",
		FileName:    "faekOutput.ts",
	}

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

	if s.FileName == "" {
		s.FileName = defaultSettings.FileName
	}

	if s.OutputStyle == "" {
		s.OutputStyle = defaultSettings.OutputStyle
	}

	return s, nil
}
