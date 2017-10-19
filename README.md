# ArchWSL
ArchLinux on WSL (Windows 10 FCU or later)


## 2017101900 Testing Installation
#### Download testing installer and rootfs.tar.gz
[arch.exe](https://github.com/yuk7/ArchWSL/releases/download/17101900/arch.exe)

[rootfs.tar.gz](https://github.com/yuk7/ArchWSL/releases/download/17101700/rootfs.tar.gz)


#### First Run arch.exe to Extract rootfs and Register to WSL
```dos
>arch.exe
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
#### Just Run arch.exe
```dos
>arch.exe
[root@PC-NAME user]#
```

#### Change Default Distribution to Arch and Run it
```dos
>wslconfig /s Arch
>wsl.exe
[root@PC-NAME user]#
```
