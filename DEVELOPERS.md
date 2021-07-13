# Developers Document

## 🛠How-to-Build

### Windows

#### Compile using `Go`(Manual build)
Create wsldl.exe(Manual)
```cmd
cd src
go build
```


Optionally, to add an icon to exe run:(Manual)
```cmd
cd src
go get github.com/akavel/rsrc
rsrc -ico ../res/%YourDistroName%/icon.ico -o wsldl.syso
go build
```

### Automatic `build.bat` usage
```
Usage:
    <no args>:
        build single wsldl.exe
    all:
        build everything, including exe with icons and default wsldl.exe
    resources:
        build only .syso files
    icons:
        build exe with icons(must be executed after running resources)
    clean:
        clean(remove) .syso files after building exe with icons
```



**Note: Creating wsldl.exe for ARM is currently not supported since Go does not support cross-compiling for `windows/arm64`.**

### Linux (cross compile)

Run this command in shell
```bash
$ cd src
$ env GOOS=windows GOARCH=amd64 go build
```

Optionally, to add an icon to the exe, create and link a resource with
```bash
$ cd src
$ go get github.com/akavel/rsrc
$ rsrc -ico ../res/%YourDistroName%/icon.ico -o wsldl.syso
$ env GOOS=windows GOARCH=amd64 go build
```
