# text-clipper


## About

`text-clipper`はTUIベースのテキストマネージャーです。
高頻度で使うテキストを保存しておき、このツールを通して一覧表示し、選択したテキストをクリップボードに呼び出すことができます。

## 環境

- WSL

## Dependexy info andlibraries used

### [atotto/clipboard](https://github.com/atotto/clipboard) ※ for copy text to clipboard

- OSX
- Windows7 (probably work on other Windows)
- Linux, Unix (requires `xclip` or `xsel` command to be installed)
- WSL (※以下の設定が必要)


#### for WSL setting to use clipboard

WSL環境の場合、`atotto/clipboard`が機能しなかった。
ここでは`bash`と`fish`用の設定を用意したので環境に合わせて設定を追加する。
仕組みは、`xcel`コマンドでWSLのクリップボード機能(書き込み)を実行させる。
前提として、`clip.exe`でクリップボード機能(書き込み)が利用できることを確認する
参考: [[wsl] 地味に便利なclip.exeでのテキストコピー](https://qiita.com/sasaki_hir/items/45885960b46f87226fd8)
※注意）他のツール等で`xclip`や`xsel`を利用している場合は影響を確認した上で実施してください。

##### 共通

-  `xclip` と`xsel`を削除する。
※なぜかいずれかが入っていたら動作しなかった

    ```
    sudo apt-get remove xclip
    sudo apt-get remove xsel
    ```

##### bash

- .bashrcを開く

    ```bash
    vi ~/.bashrc
    ```

- xclipコマンドで`cat | clip.exe`が実行されるようにfunctionを追加

    ```bash
    function xclip(){
    cat | clip.exe
    }
    ```

- 再読み込み

    ```bash
    source ~/.bashrc
    ```

##### fish

- functionを作成する

    ```bash
    vi ~/.config/fish/functions/xclip.fish
    ```

- xclipコマンドで`cat | clip.exe`が実行されるようにfunctionを追加

    ```bash
    function xclip
        cat | clip.exe
    end
    ```

  - 再読み込み

  ```bash
  source ~/.config/fish/functions/xclip.fish
  ```