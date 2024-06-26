# Text Clipper

## About

`Text Clipper` is a TUI-based text manager. It allows you to save texts that you use frequently and, through this tool, display them in a list. You can then select the text you want to copy to the clipboard.

## Environments

- WSL 22.04(LTS)
- Windows 11
- Go 1.22.1

## Windows Setup

To use SQLite within the application on Windows, please follow these steps:

1. **Install GCC**
   - Install a GCC compiler such as MinGW or TDM-GCC.
   - Verify the installation by running `gcc --version` in the command prompt to ensure it was successful.

2. **Enable CGO**
   - Set the environment variable `CGO_ENABLED` to `1`. Execute the following command in the command prompt:
     ```cmd
     set CGO_ENABLED=1
     ```

After completing these settings, proceed to install the application.

## Installation

```bash
go install github.com/o-kaisan/text-clipper@latest
```

## Dependency info and Libraries used

### [atotto/clipboard](https://github.com/atotto/clipboard) ※ for copy text to clipboard

- OSX
- Windows10 (probably work on other Windows)
- Linux, Unix (requires `xclip` or `xsel` command to be installed)
- WSL (The following settings are required)

#### For WSL Settings to Use Text Clipper

In a WSL environment, `atotto/clipboard` did not work. Here, settings for both `bash` and `fish` are prepared, so add them according to your environment. The mechanism executes the WSL clipboard functionality (write) using the `xclip` command. As a prerequisite, ensure that the clipboard functionality (write) is available with `clip.exe`.
Reference:[[wsl] 地味に便利なclip.exeでのテキストコピー](https://qiita.com/sasaki_hir/items/45885960b46f87226fd8)
*Note: If you are using `xclip` or `xsel` with other tools, check their impact before proceeding.

##### Common

- Remove `xclip` and `xsel`. If either is installed, it did not work.

    ```bash
    sudo apt-get remove xclip
    sudo apt-get remove xsel
    ```

##### Bash

- Open `.bashrc`

    ```bash
    vi ~/.bashrc
    ```

- Add a function to execute `cat | clip.exe` using the xclip command

    ```bash
    function xclip(){
        cat | clip.exe
    }
    ```

- Reload

    ```bash
    source ~/.bashrc
    ```

##### Fish

- Create a function

    ```bash
    vi ~/.config/fish/functions/xclip.fish
    ```

- Add a function to execute `cat | clip.exe` using the xclip command

    ```bash
    function xclip
        cat | clip.exe
    end
    ```

- Reload

    ```bash
    source ~/.config/fish/functions/xclip.fish
    ```

## Configuration File Location

By default, this application creates a .text-clipper directory under the user's home directory. If you want to change the location where this file is stored, please set a new path by specifying the TEXT_CLIPPER_PATH environment variable.

Bash

```bash
export TEXT_CLIPPER_PATH=/home/hoge
```

fish

```bash
set -x TEXT_CLIPPER_PATH /home/hoge
```

To make this setting permanent, add the above command to the appropriate shell configuration file (e.g., .bashrc, .bash_profile, or config.fish).

## Usage

- run

  ```bash
  text-clipper
  ```

### list view

- key binding

    | key | description |
    | --- | --- |
    | ↓/j | down |
    | ↑/k | up |
    | ctrl+a | add new item |
    | ctrl+d | delete item |
    | / | filter |
    | q/ctrl+c | quit |
    | ? | more help |
    | ... | ... |

### register view

- key binding

    | key | description |
    | --- | --- |
    | tab | move down |
    | shift+tab | move up |
    | Enter | enter over the submit button to register the item |
    | ctrl+c | back to list view |

- input form
  - title
    - 50 character limit
  - contents
    - no limit

