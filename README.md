# ArchWSL
ArchLinux on WSL (Windows 10 FCU or later)


## 2017101700 Testingi Installation
#### Download testing installer and rootfs.tar.gz
[arch-inst.exe](https://github.com/yuk7/ArchWSL/releases/download/17101700/arch-inst.exe)
[rootfs.tar.gz](https://github.com/yuk7/ArchWSL/releases/download/17101700/rootfs.tar.gz)

#### Run arch-inst.exe to Extract rootfs and Register WSL
```dos
>arch-inst.exe
Installation Complete!
```

#### Check Registerd Distribution
```dos
>wslconfig /l
~
Arch
```


## How-to-Run
#### Change Default Distribution to Arch
```
>wslconfig /s Arch
```
#### Run Default Distribution 
```dos
>wsl.exe
[root@PC-NAME user]#
```
