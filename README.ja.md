# [wsldl](https://github.com/yuk7/wsldl)
汎用WSLディストリビューション インストーラ/ランチャー

![screenshot](https://raw.githubusercontent.com/wiki/yuk7/wsldl/img/Arch_Alpine_Ubuntu.png)

[![Travis (.org)](https://img.shields.io/travis/yuk7/wsldl.svg?logo=Travis&style=flat-square)](https://travis-ci.org/yuk7/wsldl)
[![AppVeyor](https://img.shields.io/appveyor/ci/yuk7/wsldl.svg?logo=AppVeyor&style=flat-square)](https://ci.appveyor.com/project/yuk7/wsldl)
[![Github All Releases](https://img.shields.io/github/downloads/yuk7/wsldl/total.svg?style=flat-square)](https://github.com/yuk7/wsldl/releases/latest)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
![License](https://img.shields.io/github/license/yuk7/wsldl.svg?style=flat-square)



## 💻システム要件
* Windows 10 1709 Fall Creators Update x64以降
* WSL機能が有効であること

## 📦既にパッケージされているものをインストール(外部のプロジェクト)
#### 1. zipをダウンロード
[Alpine Linux](https://github.com/yuk7/AlpineWSL)

[Alpine Linux with Git-LFS and Sphinx](https://github.com/binarylandscapes/AlpineWSL) (by [binarylandscapes](https://github.com/binarylandscapes))

[Arch Linux](https://github.com/yuk7/ArchWSL)

[Artix Linux](https://github.com/hdk5/ArtixWSL) (by [hdk5](https://github.com/hdk5))

[CentOS](https://github.com/fbigun/WSL-Distro-Rootfs) (by [fbigun](https://github.com/fbigun))

[Void Linux (glibc)](https://github.com/am11/VoidWSL) (by [am11](https://github.com/am11))

[Void Linux (musl-libc)](https://github.com/am11/VoidMuslWSL) (by [am11](https://github.com/am11))

#### 2. zip内のファイルをすべて同じ場所に展開

#### 3. exeを実行し、インストールを開始
exeのファイル名はインストール名に使用されます。
リネームすることでご自由な名前でインストールすることが出来、複数インストールも可能です。


## 🔧あなたのrootfsでインストール
#### 1. [Launcher.exeをダウンロード](https://github.com/yuk7/wsldl/releases/latest)
#### 2. Launcher.exeのファイル名を登録するインスタンス名に変更
(例:Arch.exeにリネームすると「Arch」という名前でインストールされます。
#### 3. あなたが用意したrootfs.tar.gzをexeと同じ場所に配置します。
#### 4. exeを実行するとインストールされます。
※あなたのrootfs.tar.gzとexeをパッケージし、自由に再配布することもできます📦


## 📝使い方(インストール後)
#### exe Usage
```cmd
Usage :
    <引数なし>
      - デフォルト設定で新しいシェルを起動します

    run <command line>
      - 与えられたコマンドラインをインスタンス内で実行します。 カレントディレクトリが引き継がれます。

    runp <command line (windowsのパスを含む)>
      - 与えられたコマンドラインのパスを変換した上でインスタンス内で実行します。

    config [setting [value]]
      - `--default-user <user>`: インスタンスのデフォルトユーザーを<user>に設定します。
      - `--default-uid <uid>`: インスタンスのデフォルトユーザーのuidを<uid>に設定します。
      - `--append-path <on|off>`: Windows側のPATH設定をLinux側に引き継ぐ機能のon/offを設定します。
      - `--mount-drive <on|off>`: Windowsのドライブをマウントする機能のon/offを設定します。

    get [setting]
      - `--default-uid`: インスタンスのデフォルトユーザーのuidを取得します。
      - `--append-path`: Windows側のPATH設定をLinux側に引き継ぐ機能のon/offを確認します。
      - `--mount-drive`: Windowsのドライブをマウントする機能のon/offを確認します。
      - `--lxuid`: システム内部で使用されているLxUIDを取得します。

    backup [contents]
      - `--tgz`: tarを使用し、カレントディレクトリにbackup.tar.gzを出力します。
      - `--reg`: 設定のレジストリファイルをbackup.regとしてカレントディレクトリに出力します。
      
    clean
      - インスタンスをアンインストールします。

    help
      - helpを表示します。
```


#### exeをそのまま実行
カレントディレクトリが引き継がれず、実行されます。
```cmd
>{インスタンス名}.exe
[root@PC-NAME ~]#
```

#### 任意のコマンドを実行させる
```cmd
>{インスタンス名}.exe run uname -r
4.4.0-43-Microsoft
```

#### 任意のコマンドをパス変換し実行させる
```cmd
>{インスタンス名}.exe runp echo C:\Windows\cmd.exe
/mnt/c/Windows/cmd.exe

```

#### デフォルトユーザー変更(ディストリビューションにidコマンドが必要です)
```cmd
>{インスタンス名}.exe config --default-user user

>{インスタンス名}.exe
[user@PC-NAME dir]$
```

#### インスタンスをアンインストール
```cmd
>{インスタンス名}.exe clean

```

## 📄ライセンス
[MIT](https://github.com/yuk7/wsldl/blob/master/LICENSES.md)

Copyright (c) 2017-2020 yuk7
