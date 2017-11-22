# WSL-DistroLauncher
Gernal Perpose WSL Distribution Installer&Launcher

![screenshot](https://raw.githubusercontent.com/wiki/yuk7/WSL-DistroLauncher/img/Arch_Alpine_Ubuntu.png)

## Install with prebuilt rootfs
#### 1. Download installer zip
[Alpine Linux 3.6.2](https://github.com/yuk7/WSL-DistroLauncher/releases/download/17112101/Alpine.zip)(md5:e7476fd779e27c75aa61265bcc10df5)

[Arch Linux](https://github.com/yuk7/WSL-DistroLauncher/releases/download/17112100/Arch.zip)(md5:50e300a63395c8f11f2f13663c009c85)

#### 2. Extract all files in zip file to same directory

#### 3.Run exe to Extract rootfs and Register to WSL
Exe filename is using to the instance name to register.
If you rename it you can register with a diffrent name.


## Install with any rootfs
#### 1. Download [Launcher.exe](https://github.com/yuk7/WSL-DistroLauncher/releases/download/17112100/Launcher.exe) and rename it for distribution name to register.
(Ex:Rename to Arch.exe if you want to "Arch" for the Instance name)
#### 2. Put your rootfs.tar.gz in same directory as exe (Installation directory)
#### 3. Run exe to install. This process may take a few minutes.


## How-to-Use(for Installed Instance)
#### exe Usage
```dos
Useage :
    <no args>
      - Launches the distro's default behavior. By default, this launches your default shell.

    run <command line>
      - Run the given command line in that distro.

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
```


#### Just Run exe
```dos
>{InstanceName}.exe
[root@PC-NAME user]#
```

#### Run with command line
```dos
>{InstanceName}.exe run uname -r
4.4.0-43-Microsoft

```

#### Change Default User(id command required)
```dos
>{InstanceName}.exe config --default-user user

>{InstanceName}.exe
[user@PC-NAME dir]$
```


#### How to uninstall instance
```dos
>wslconfig /u {InstanceName}

```
