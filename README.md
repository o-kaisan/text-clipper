# Text Clipper

## About

`Text Clipper` is a TUI-based text manager. It allows you to save texts that you use frequently and, through this tool, display them in a list. You can then select the text you want to copy to the clipboard.

## Environments

- WSL 22.04(LTS)
- Go 1.22.1

## installation

```
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

## This App Uses SQLite

This application utilizes SQLite for data storage. The default database path is `$HOME/.text-clipper/text-clipper.db`. If you wish to specify a different path, please set the environment variable `TEXT_CLIPPER_DB_PATH`.

Bash

```bash
export TEXT_CLIPPER_DB_PATH=/home/hoge/fuga.db
```

fish

```bash
set -x TEXT_CLIPPER_DB_PATH /home/hoge/fuga.db
```

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
