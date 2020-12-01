# wsldl
Advanced WSL Distribution Launcher / Installer


![screenshot](https://raw.githubusercontent.com/wiki/yuk7/wsldl/img/Arch_Alpine_Cent.png)

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/yuk7/wsldl/Mingw-w64%20Cross%20CI?logo=GitHub&style=flat-square)](https://github.com/yuk7/wsldl/actions?query=workflow%3A%22Mingw-w64+Cross+CI%22)
[![AppVeyor](https://img.shields.io/appveyor/ci/yuk7/wsldl.svg?logo=AppVeyor&style=flat-square)](https://ci.appveyor.com/project/yuk7/wsldl)
[![Github All Releases](https://img.shields.io/github/downloads/yuk7/wsldl/total.svg?style=flat-square)](https://github.com/yuk7/wsldl/releases/latest)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
![License](https://img.shields.io/github/license/yuk7/wsldl.svg?style=flat-square)


### [Detailed documentation is here](https://git.io/wsldl-doc)

## 💻Requirements
* Windows 10 1709 Fall Creators Update 64bit or later.
* Windows Subsystem for Linux feature is enabled.

## 📦Install with Prebuilt Packages
[**You can see List on docs**](https://wsldl-pg.github.io/docs/Using-wsldl/#distros)

**Note:**
Exe filename is using to the instance name to register.
If you rename it, you can register with a different name.


## 🔧Install with any rootfs
#### 1. [Download Launcher.exe](https://github.com/yuk7/wsldl/releases/latest)
#### 2. Rename it for distribution name to register.
(Ex:Rename to Arch.exe if you want to use "Arch" for the Instance name)
#### 3. Put your rootfs.tar(.gz) in same directory as exe (Installation directory)
#### 4. Run exe to install. This process may take a few minutes.

## 🔗Use as a Launcher for already installed distribution
#### 1. [Download Launcher.exe](https://github.com/yuk7/wsldl/releases/latest)
#### 2. Rename it for registerd instance name.
Please check the registered instance name of the distribution with `wslconfig /l` command.
(Ex: If the instance name is "Ubuntu-20.04", rename `Launcher.exe` to `Ubuntu-20.04.exe`)
#### 4. Run exe to Launch instance or configuration.
For details, please see the help. (`{InstanceName}.exe --help`)

Note: You can distribute your distribution including wsldl exe.

## 📝How-to-Use(for Installed Instance)
#### exe Usage
```
Usage :
    <no args>
      - Open a new shell with your default settings.

    run <command line>
      - Run the given command line in that distro. Inherit current directory.

    runp <command line (includes windows path)>
      - Run the path translated command line in that distro.

    config [setting [value]]
      - `--default-user <user>`: Set the default user for this distro to <user>
      - `--default-uid <uid>`: Set the default user uid for this distro to <uid>
      - `--append-path <on|off>`: Switch of Append Windows PATH to $PATH
      - `--mount-drive <on|off>`: Switch of Mount drives
      - `--default-term <default|wt|flute>`: Set default terminal window

    get [setting]
      - `--default-uid`: Get the default user uid in this distro
      - `--append-path`: Get on/off status of Append Windows PATH to $PATH
      - `--mount-drive`: Get on/off status of Mount drives
      - `--wsl-version`: Get WSL Version 1/2 for this distro
      - `--default-term`: Get Default Terminal for this distro launcher
      - `--lxguid`: Get WSL GUID key for this distro

    backup [contents]
      - `--tgz`: Output backup.tar.gz to the current directory using tar command
      - `--reg`: Output settings registry file to the current directory

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

#### Run with command line with path translation
```cmd
>{InstanceName}.exe runp echo C:\Windows\System32\cmd.exe
/mnt/c/Windows/System32/cmd.exe
```

#### Change Default User(id command required)
```cmd
>{InstanceName}.exe config --default-user user

>{InstanceName}.exe
[user@PC-NAME dir]$
```

#### Set "Windows Terminal" as default terminal
```cmd
>{InstanceName}.exe config --default-term wt
```

#### How to uninstall instance
```cmd
>{InstanceName}.exe clean

```

## 🛠How-to-Build
Please see [DEVELOPERS.md](https://github.com/yuk7/wsldl/blob/main/DEVELOPERS.md)

## 📄License
[MIT](https://github.com/yuk7/wsldl/blob/main/LICENSES.md)

Copyright (c) 2017-2020 yuk7
