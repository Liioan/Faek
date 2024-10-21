package utils

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearConsole() {
	clear, available := clear[runtime.GOOS]
	if available {
		clear()
	}
}

func LogToFile(data string) {
	file, err := os.Create("./test.txt")
	if err != nil {
		log.Fatal("could not create file")
	}
	defer file.Close()
	file.Write([]byte(data))
}
