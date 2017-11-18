# ArchWSL
Install ArchLinux as a WSL Instance


![screenshot](https://raw.githubusercontent.com/wiki/yuk7/ArchWSL/img/Arch_and_Ubuntu.png)

## Install (17111800 version)
#### Download testing installer and rootfs.tar.gz
[Arch.exe](https://github.com/yuk7/ArchWSL/releases/download/17111800/Arch.exe) (Release:17111800/md5:c419f21c77c8987923d48ba7ff1787f7)

[rootfs.tar.gz](https://github.com/yuk7/ArchWSL/releases/download/17102300/rootfs.tar.gz) (Release:17102300/md5:f0660ee8b236413429de8d05ea785d3b)


#### First Run Arch.exe to Extract rootfs and Register to WSL
Excutable filename is using to distribution name to register.

If you rename it you can register with a diffrent name.

```dos
>Arch.exe
~
Installation Complete!
```
This process may take a few minutes.



#### Check Registerd Distribution
```dos
>wslconfig /l
~
Arch
```

## How-to-Use
#### Arch.exe Usage
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


#### Just Run Arch.exe
```dos
>Arch.exe
[root@PC-NAME user]#
```

#### Run with command line
```dos
>Arch.exe run uname -r
4.4.0-43-Microsoft

```

#### Change Default Distribution to Arch and Run it
```dos
>wslconfig /s Arch
>wsl
[root@PC-NAME dir]#
```

#### Change Default User
```dos
>Arch.exe config --default-user user

>Arch.exe
[user@PC-NAME dir]$
```


#### How to uninstall instance
```dos
>wslconfig /u Arch

```
