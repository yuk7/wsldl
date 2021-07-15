# wsldl
Advanced WSL Distribution Launcher / Installer


![screenshot](https://raw.githubusercontent.com/wiki/yuk7/wsldl/img/Arch_Alpine_Cent.png)

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/yuk7/wsldl/Continuous%20Integration?logo=GitHub&style=flat-square)](https://github.com/yuk7/wsldl/actions?query=workflow%3A%22Continuous-Integration%22)
[![Github All Releases](https://img.shields.io/github/downloads/yuk7/wsldl/total.svg?style=flat-square)](https://github.com/yuk7/wsldl/releases/latest)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
![License](https://img.shields.io/github/license/yuk7/wsldl.svg?style=flat-square)


### [Detailed documentation is here](https://git.io/wsldl-doc)

## 💻Requirements
* Windows 10 1709 Fall Creators Update or later(x64/arm64).
* Windows Subsystem for Linux feature is enabled.

## 📦Install with Prebuilt Packages
[**You can see List on docs**](https://wsldl-pg.github.io/docs/Using-wsldl/#distros)

**Note:**
Exe filename is using to the instance name to register.
If you rename it, you can register with a different name.


## 🔧Install with any rootfs
#### 1. [Download wsldl.exe](https://github.com/yuk7/wsldl/releases/latest)
(wsldl.exe is x86_64, wsldl_arm64.exe is ARM64 build)
#### 2. Rename it for distribution name to register.
(Ex:Rename to Arch.exe if you want to use "Arch" for the Instance name)
#### 3. Put your rootfs.tar(.gz) in same directory as exe (Installation directory)
#### 4. Run exe to install. This process may take a few minutes.

## 🔗Use as a Launcher for already installed distribution
#### 1. [Download wsldl.exe](https://github.com/yuk7/wsldl/releases/latest)
(wsldl.exe is x86_64, wsldl_arm64.exe is ARM64 build)
#### 2. Rename it for registerd instance name.
Please check the registered instance name of the distribution with `wslconfig /l` command.
(Ex: If the instance name is "Ubuntu-20.04", rename `wsldl.exe` to `Ubuntu-20.04.exe`)
#### 4. Run exe to Launch instance or configuration.
For details, please see the help. (`{InstanceName}.exe help`)

Note: You can distribute your distribution including wsldl exe.

## 📝How-to-Use(for Installed Instance)
#### exe Usage
```
Usage :
    <no args>
      - Open a new shell with your default settings.

    run <command line>
      - Run the given command line in that instance. Inherit current directory.

    runp <command line (includes windows path)>
      - Run the given command line in that instance after converting its path.

    config [setting [value]]
      - `--default-user <user>`: Set the default user of this instance to <user>.
      - `--default-uid <uid>`: Set the default user uid of this instance to <uid>.
      - `--append-path <true|false>`: Switch of Append Windows PATH to $PATH
      - `--mount-drive <true|false>`: Switch of Mount drives
      - `--default-term <default|wt|flute>`: Set default type of terminal window.

    get [setting]
      - `--default-uid`: Get the default user uid in this instance.
      - `--append-path`: Get true/false status of Append Windows PATH to $PATH.
      - `--mount-drive`: Get true/false status of Mount drives.
      - `--wsl-version`: Get the version os the WSL (1/2) of this instance.
      - `--default-term`: Get Default Terminal type of this instance launcher.
      - `--lxguid`: Get WSL GUID key for this instance.

    backup [contents]
      - `--tar`: Output backup.tar to the current directory.
      - `--reg`: Output settings registry file to the current directory.

    clean
      - Uninstall that instance.

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
Please see [DEVELOPERS.md](DEVELOPERS.md)

## 📄License
[MIT](LICENSES.md)

Copyright (c) 2017-2021 [yuk7](https://github.com/yuk7)
