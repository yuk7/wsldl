# ArchWSL
ArchLinux on WSL (Windows 10 FCU or later)


## 2017101901 Testing Installation
#### Download testing installer and rootfs.tar.gz
[Arch.exe](https://github.com/yuk7/ArchWSL/releases/download/17101901/Arch.exe) (Release:17101901/md5:39c4da30287d5ac408b8ccf508f7cfc6)

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


## How-to-Run
#### Just Run Arch.exe
```dos
>Arch.exe
[root@PC-NAME user]#
```

#### Change Default Distribution to Arch and Run it
```dos
>wslconfig /s Arch
>wsl
[root@PC-NAME user]#
```
