# wsldl
General Purpose WSL Distribution Installer & Launcher


![screenshot](https://raw.githubusercontent.com/wiki/yuk7/wsldl/img/Arch_Alpine_Cent.png)

[![Travis (.org)](https://img.shields.io/travis/yuk7/wsldl.svg?logo=Travis&style=flat-square)](https://travis-ci.org/yuk7/wsldl)
[![AppVeyor](https://img.shields.io/appveyor/ci/yuk7/wsldl.svg?logo=AppVeyor&style=flat-square)](https://ci.appveyor.com/project/yuk7/wsldl)
[![Github All Releases](https://img.shields.io/github/downloads/yuk7/wsldl/total.svg?style=flat-square)](https://github.com/yuk7/wsldl/releases/latest)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
![License](https://img.shields.io/github/license/yuk7/wsldl.svg?style=flat-square)



## 💻Requirements
* Windows 10 1709 Fall Creators Update 64bit or later.
* Windows Subsystem for Linux feature is enabled.

## 📦Install with prebuilt rootfs
#### 1. Download installer zip
[Alpine Linux](https://github.com/yuk7/AlpineWSL)

[Alpine Linux with Git-LFS and Sphinx](https://github.com/binarylandscapes/AlpineWSL) (by [binarylandscapes](https://github.com/binarylandscapes))

[Arch Linux](https://github.com/yuk7/ArchWSL)

[Artix Linux](https://github.com/hdk5/ArtixWSL) (by [hdk5](https://github.com/hdk5))

[CentOS](https://github.com/yuk7/CentWSL)

[CentOS](https://github.com/fbigun/WSL-Distro-Rootfs) (by [fbigun](https://github.com/fbigun))

[Void Linux (glibc)](https://github.com/am11/VoidWSL) (by [am11](https://github.com/am11))

[Void Linux (musl-libc)](https://github.com/am11/VoidMuslWSL) (by [am11](https://github.com/am11))

#### 2. Extract all files in zip file to same directory

#### 3.Run exe to Extract rootfs and Register to WSL
Exe filename is using to the instance name to register.
If you rename it, you can register with a different name.


## 🔧Install with any rootfs
#### 1. [Download Launcher.exe](https://github.com/yuk7/wsldl/releases/latest)
#### 2. Rename it for distribution name to register.
(Ex:Rename to Arch.exe if you want to use "Arch" for the Instance name)
#### 3. Put your rootfs.tar.gz in same directory as exe (Installation directory)
#### 4. Run exe to install. This process may take a few minutes.


## 📝How-to-Use(for Installed Instance)
#### exe Usage
```cmd
Usage :
    <no args>
      - Open a new shell with your default settings.

    run <command line>
      - Run the given command line in that distro. Inherit current directory.

    config [setting [value]]
      - `--default-user <user>`: Set the default user for this distro to <user>
      - `--default-uid <uid>`: Set the default user uid for this distro to <uid>
      - `--append-path <on|off>`: Switch of Append Windows PATH to $PATH
      - `--mount-drive <on|off>`: Switch of Mount drives

    get [setting]
      - `--default-uid`: Get the default user uid in this distro
      - `--append-path`: Get on/off status of Append Windows PATH to $PATH
      - `--mount-drive`: Get on/off status of Mount drives
      - `--lxuid`: Get LxUID key for this distro

    backup
      - Output backup.tar.gz to the current directory using tar command.
      
    clean
      - Uninstall the distro.

    help
      - Print this usage message.
```


#### Just Run exe
```cmd
>{InstanceName}.exe
[root@PC-NAME user]#
```

#### Run with command line
```cmd
>{InstanceName}.exe run uname -r
4.4.0-43-Microsoft

```

#### Change Default User(id command required)
```cmd
>{InstanceName}.exe config --default-user user

>{InstanceName}.exe
[user@PC-NAME dir]$
```

#### How to uninstall instance
```cmd
>{InstanceName}.exe clean

```

## 🛠How-to-Build
### Windows

#### Visual Studio or Build Tools 2017+

Use `Developer Command Prompt for Visual Studio` or run these in the Windows Command Prompt
```cmd
:: locate VS base installation path using vswhere
SET vswherePath=%ProgramFiles(x86)%\Microsoft Visual Studio\Installer\vswhere.exe
FOR /F "tokens=*" %i IN ('
      "%vswherePath%" -latest -prerelease -products *               ^
        -requires Microsoft.VisualStudio.Component.VC.Tools.x86.x64 ^
        -property installationPath'
      ) DO SET vsBase=%i

:: initialize x64 build environment
CALL "%vsBase%\vc\Auxiliary\Build\vcvarsall.bat" x64
```

To compile Launcher.exe
```cmd
cl /nologo /O2 /W4 /WX /Ob2 /Oi /Oy /Gs- /GF /Gy /Tc main.c /Fe:Launcher.exe Advapi32.lib Shell32.lib
```

Optionally, to add an icon to the exe, create and link a resource with
```cmd
SET YourDistroName=Fedora

:: create resources
rc /nologo res\%YourDistroName%\res.rc

:: compile to %YourDistroName%.exe
cl /nologo /O2 /W4 /WX /Ob2 /Oi /Oy /Gs- /GF /Gy /Tc main.c /Fe:%YourDistroName%.exe ^
  Advapi32.lib Shell32.lib res\%YourDistroName%\res.res
```

#### MinGW
Install x86_64 version of MSYS2(https://www.msys2.org).

Run these commands in msys shell
```bash
$ pacman -S mingw-w64-x86_64-toolchain # install tool chain
$ gcc -std=c99 --static main.cpp -o Launcher.exe # compile main.c
```

Optionally, to add an icon to the exe, create and link a resource with
```bash
YourDistroName=Fedora
$ windres res/$YourDistroName/res.rc res.o # compile resource
$ gcc -std=c99 --static main.cpp -o Launcher.exe res.o # compile main.cpp
```

### Linux (cross compile)
Install mingw-w64 toolchain include gcc-mingw-w64-x86-64.

Run this command in shell
```bash
 $ x86_64-w64-mingw32-gcc -std=c99 --static main.c -o Launcher.exe # compile main.c
```
## 📄License
[MIT](https://github.com/yuk7/wsldl/blob/master/LICENSES.md)

Copyright (c) 2017-2019 yuk7
