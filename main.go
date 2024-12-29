package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/manifoldco/promptui"
)

func listFilesInDir(dir string) ([]string, error) {
    directory, err := os.Open(dir)
    if err != nil {
        return nil, err
    }
    defer directory.Close()

    entries, err := directory.Readdir(-1)
    if err != nil {
        return nil, err
    }

    var files []string
    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".desktop") {
            files = append(files, filepath.Join(dir, entry.Name()))
        }
    }
    return files, nil
}

func getDesktopFiles() ([]string, []string) {
    systemDirs := []string{
        "/usr/share/applications/",
        "/usr/local/share/applications/",
    }
    userDir := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")

    var systemFiles, userFiles []string

    for _, dir := range systemDirs {
        files, err := listFilesInDir(dir)
        if err == nil {
            systemFiles = append(systemFiles, files...)
        }
    }

    files, err := listFilesInDir(userDir)
    if err == nil {
        userFiles = append(userFiles, files...)
    }

    return systemFiles, userFiles
}

func searchFiles(files []string, term string) []string {
    var matches []string
    for _, file := range files {
        if strings.Contains(strings.ToLower(filepath.Base(file)), strings.ToLower(term)) {
            matches = append(matches, file)
        }
    }
    return matches
}

func promptForFile(files []string) (string, error) {
    names := make([]string, len(files))
    for i, file := range files {
        names[i] = filepath.Base(file)
    }

    prompt := promptui.Select{
        Label: "Select a file",
        Items: names,
        Size: 15,
    }

    index, _, err := prompt.Run()
    if err != nil {
        return "", err
    }

    return files[index], nil
}

func isWritable(fileName string) bool {
    file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
    if err != nil {
        return false
    }
    defer file.Close()
    return true
}

func openFileInEditor(fileName string) {
    editor := os.Getenv("EDITOR")
    if editor == "" {
        editor = "nano"
    }

    if !isWritable(fileName) {
        fmt.Printf("The file '%s' requires elevated privileges. Please provide your password.\n", fileName)
        cmd := exec.Command("sudo", editor, fileName)
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            fmt.Println("Error opening file with sudo privileges:", err)
        }
        return
    }

    cmd := exec.Command(editor, fileName)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if err := cmd.Run(); err != nil {
        fmt.Println("Error opening file in editor:", err)
    }
}

func main() {
    getFlag := flag.Bool("get", false, "List desktop files")
    searchFlag := flag.String("search", "", "Search desktop files by name")
    userFlag := flag.Bool("u", false, "Show user-specific desktop files only")
    systemFlag := flag.Bool("s", false, "Show system-wide desktop files only")
    helpFlag := flag.Bool("help", false, "Show usage guide")

    flag.Parse()

    if *helpFlag {
        showUsage()
        return
    }

    if !*getFlag && *searchFlag == "" {
        fmt.Println("Error: You must either use --get or --search flag.")
        showUsage()
        os.Exit(1)
    }

    if *getFlag && *searchFlag != "" {
        fmt.Println("Error: You cannot use --get and --search flags together.")
        showUsage()
        os.Exit(1)
    }

    if *getFlag && *userFlag && *systemFlag {
        fmt.Println("Error: You cannot use both --u and --s flags together.")
        showUsage()
        os.Exit(1)
    }

    systemFiles, userFiles := getDesktopFiles()

    if *getFlag {
        if *userFlag {
            fmt.Println("User-specific desktop files:")
            file, err := promptForFile(userFiles)
            if err == nil {
                openFileInEditor(file)
            }
        } else if *systemFlag {
            fmt.Println("System-wide desktop files:")
            file, err := promptForFile(systemFiles)
            if err == nil {
                openFileInEditor(file)
            }
        } else {
            fmt.Println("System-wide desktop files:")
            for _, file := range systemFiles {
                fmt.Println("  ", filepath.Base(file))
            }
            fmt.Println("\nUser-specific desktop files:")
            for _, file := range userFiles {
                fmt.Println("  ", filepath.Base(file))
            }
        }
    } else if *searchFlag != "" {
        allFiles := append(systemFiles, userFiles...)
        matches := searchFiles(allFiles, *searchFlag)
        fmt.Println("Matching desktop files:")
        file, err := promptForFile(matches)
        if err == nil {
            openFileInEditor(file)
        }
    } else {
        showUsage()
    }
}

func showUsage() {
    fmt.Println("Usage:")
    fmt.Println("  deskedit --get        List all desktop files")
    fmt.Println("  deskedit --get -s     List system-wide desktop files")
    fmt.Println("  deskedit --get -u     List user-specific desktop files")
    fmt.Println("  deskedit --search     Search desktop files by name")
    fmt.Println("  deskedit --help       Show this usage guide")
}