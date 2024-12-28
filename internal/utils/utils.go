package utils

import (
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
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

func LogToDebug(data string) {
	file, err := os.Create("./debug.log")
	if err != nil {
		log.Fatal("could not create file")
	}
	defer file.Close()
	file.Write([]byte(data))
}

func Random(min, max int) int {
	return rand.Intn(max+1-min) + min
}

func ParseInt(s string, defVal int) int {
	res := defVal
	num, err := strconv.Atoi(s)
	if err == nil {
		res = num
	}
	return res
}
