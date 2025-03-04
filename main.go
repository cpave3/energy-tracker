package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Entry struct {
	Timestamp string `json:"timestamp"`
	Mood      int    `json:"mood"`
	Energy    int    `json:"energy"`
	Activity  string `json:"activity"`
	Notes     string `json:"notes"`
}

func main() {

	timestamp := time.Now().Format(time.RFC3339)

	if !confirm("Do you want to log your energy level?") {
		fmt.Println("Aborted.")
		return
	}

	mood, err := prompt("Mood (1-10):")
	if err != nil {
		fmt.Println("Error capturing mood:", err)
		return
	}

	energy, err := prompt("Energy level (1-10):")
	if err != nil {
		fmt.Println("Error capturing energy:", err)
		return
	}

	activity, err := prompt("What have you been doing?")
	if err != nil {
		fmt.Println("Error capturing activity:", err)
		return
	}

	notes, err := prompt("Notes:")
	if err != nil {
		fmt.Println("Error capturing notes:", err)
		return
	}

	entry := Entry{
		Timestamp: timestamp,
		Mood:      toInt(mood),
		Energy:    toInt(energy),
		Activity:  activity,
		Notes:     notes,
	}

	saveEntry(entry)
}

func confirm(question string) bool {
	cmd := exec.Command("zenity", "--question", "--text", question)
	err := cmd.Run()
	return err == nil
}

func prompt(question string) (string, error) {
	cmd := exec.Command("zenity", "--entry", "--text", question)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func saveEntry(entry Entry) {
	dir := filepath.Join(os.Getenv("HOME"), ".tracker")
	filePath := filepath.Join(dir, "responses.jsonl")

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	data, _ := json.Marshal(entry)
	file.WriteString(string(data) + "\n")
}
