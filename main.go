package main

import (
	"fmt"
	"os"
	"os/exec"
)

func chooseEditor() string {
	editors := []string{"nano", "vi", "pico"}

	fmt.Println("No default text editor set in the environment variables.")
	fmt.Println("Please choose an editor from the following options:")

	for i, editor := range editors {
		fmt.Printf("%d. %s\n", i+1, editor)
	}

	var choice int
	fmt.Print("Enter the number of your choice: ")
	_, err := fmt.Scan(&choice)
	if err != nil || choice < 1 || choice > len(editors) {
		fmt.Println("Invalid choice. Using 'nano' as the default editor.")
		return "nano" 
	}

	return editors[choice-1]
}

func isWritable(fileName string) bool {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

func openInEditor(fileName string) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = chooseEditor()
		os.Setenv("EDITOR", editor)
	}

	if !isWritable(fileName) {
		fmt.Printf("The file '%s' is either not writable or does not exist. Opening with sudo privileges...\n", fileName)

		cmd := exec.Command("sudo", editor, fileName)

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error opening the editor:", err)
			return
		}
		fmt.Println("Successfully opened the file with the editor:", editor)
		return
	}

	cmd := exec.Command(editor, fileName)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error opening the editor:", err)
	} else {
		fmt.Println("Successfully opened the file with the editor:", editor)
	}
}

func main() {
	fileName := "example.txt"

	openInEditor(fileName)
}
