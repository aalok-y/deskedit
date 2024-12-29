# Deskedit

DeskEditDeskEdit is a command-line tool for managing and editing `.desktop` files on your system. It allows you to 
list, search, and open `.desktop` files in your preferred text editor.

**Note**: The programs uses ```EDITOR``` environment variable to determine the default text editor.If it is not set then Nano editor is used

## Installation

Clone the repository and build the project:

```sh
git clone https://github.com/aalok-deskedit.git
cdeskedit
go build -deskedit main.go
```

## Usage

### Show usage guide
```sh
deskedit --help
```

### List all desktop files

```sh
deskedit --get
```

### List system-wide desktop files

```sh
deskedit --get -s
```

### List user-specific desktop files

```sh
deskedit --get -u
```

### Search desktop files by name

```sh
deskedit --search <term>
```

#### Flags

- `--get`:  
  List all desktop files.

- `--get -s`:  
  List and edit system-wide desktop files.

- `--get -u`:  
  List and edit user-specific desktop files.

- `--search <term>`:  
  Search desktop files by name.

- `--help`:  
  Show usage guide.
