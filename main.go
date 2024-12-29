package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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


func listFilesInDir() ([]string, []string, error) {
	// System-wide directories
	systemDirs := []string{
		"/usr/share/applications/",
		"/usr/local/share/applications/",
	}

	// User-specific directory (considering home directory expansion)
	userDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")

	// Arrays to store file names
	var systemFiles []string
	var userFiles []string

	// Helper function to read files from a directory
	readDirectory := func(dir string) ([]string, error) {
		// Open the directory
		directory, err := os.Open(dir)
		if err != nil {
			return nil, err
		}
		defer directory.Close()

		// Get the list of files and directories in the specified directory
		entries, err := directory.Readdir(-1) // -1 means to read all files
		if err != nil {
			return nil, err
		}

		// Initialize a slice to store the file names
		var files []string

		// Iterate over the entries and add file names to the slice
		for _, entry := range entries {
			// Check if it's a file (not a directory)
			if !entry.IsDir() {
				files = append(files, entry.Name()) // Append the file name
			}
		}
		return files, nil
	}

	// Process system-wide directories
	for _, dir := range systemDirs {
		files, err := readDirectory(dir)
		if err != nil {
			fmt.Println("Error reading system directory", dir, ":", err)
			continue
		}
		systemFiles = append(systemFiles, files...)
	}

	// Process user-specific directory
	userFiles, err := readDirectory(userDir)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading user directory %s: %v", userDir, err)
	}

	// Return both arrays
	return systemFiles, userFiles, nil
}



func main() {
	// fileName := "example.txt"
	// openInEditor(fileName)
	systemFiles, userFiles, err := listFilesInDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print system-wide application files
	fmt.Println("System-wide application files:")
	for _, file := range systemFiles {
		fmt.Println(file)
	}

	// Print user-specific application files
	fmt.Println("\nUser-specific application files:")
	for _, file := range userFiles {
		fmt.Println(file)
	}
}
