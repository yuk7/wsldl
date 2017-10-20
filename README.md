# ArchWSL
ArchLinux on WSL (Windows 10 FCU or later)


## 2017102000 Testing Installation
#### Download testing installer and rootfs.tar.gz
[Arch.exe](https://github.com/yuk7/ArchWSL/releases/download/17102000/Arch.exe) (Release:17102000/md5:e281f65b65aae3cc976cbe65eb0cc287)

[rootfs.tar.gz](https://github.com/yuk7/ArchWSL/releases/download/17101901/rootfs.tar.gz) (Release:17101901/md5:0080e1df5b1de2b567b288b7a1bd3f5e)


#### First Run Arch.exe to Extract rootfs and Register to WSL
Excutable filename is using to distribution name to register.

If you rename it you can register with a diffrent name.

```dos
>Arch.exe
~
Installation Complete!
```
This process may take a few minutes.

be partient:)


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
          -  Run the given command line in that distro.

    config [setting [value]]
      - `--default-uid <uid>`: Set the default user uid for this distro to <uid>
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
>Arch.exe run id user -u
1000

>Arch.exe config --default-uid 1000

>Arch.exe
[user@PC-NAME dir]$
```


#### How to uninstall instance
```dos
>wslconfig /u Arch

```
