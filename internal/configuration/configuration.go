package configuration

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	e "github.com/liioan/faek/internal/errors"
	v "github.com/liioan/faek/internal/variants"
)

const settingsFilePath = "/faek_settings.json"
const settingsDirectoryPath = "/.config/faek"

type Settings struct {
	FileName    string    `json:"fileName"`
	Language    v.Variant `json:"lang"`
	OutputStyle v.Variant `json:"outputStyle"`
	Indent      string    `json:"indent"`
}

var defaultSettings = Settings{
	OutputStyle: v.Terminal,
	FileName:    "faekOutput.ts",
	Language:    v.TypeScript,
	Indent:      "2",
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

func SaveUserSettings(settings *Settings) error {

	previousSettings, err := GetUserSettings()
	if err != nil {
		previousSettings = defaultSettings
	}

	if settings.FileName == "" {
		settings.FileName = previousSettings.FileName
	}
	if settings.OutputStyle == "" {
		settings.OutputStyle = previousSettings.OutputStyle
	}
	if settings.Language == "" {
		settings.Language = previousSettings.Language
	}
	if settings.Indent == "" {
		settings.Indent = previousSettings.Indent
	}

	settings.FileName = strings.Split(settings.FileName, ".")[0]

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

func (s Settings) GetFullFileName() string {
	res := s.FileName
	switch s.Language {
	case v.TypeScript:
		res += ".ts"
	case v.JavaScript:
		res += ".js"
	case v.JSON:
		res += ".json"
	}

	return res
}
