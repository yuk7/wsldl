# Developers Document

## 🛠How-to-Build

## Sidenote: To compile ARM binaries, Please install `Go` 1.17 by running this:
```bash
go get golang.org/dl/go1.17rc1
#Additional step if you want to cross-compile from Linux/MacOS
export PATH=$PATH:~/go/bin
```

### Windows

#### Compile using `Go`

##### Automatic `build.bat` usage

```dos
Usage:
    <no args>/single:
        Build default wsldl.exe
    all:
        Build everything, including exe with icons and default wsldl.exe
    resources:
        Build only .syso files
    icons:
        Build exe with icons(must be executed after running resources)
    clean:
        Clean(remove) .syso files after building exe with icons
    singlewor:
        Build exe without any manifest/resources
```
To cross-compile ARM64 from AMD64 or vice versa, run:
```cmd
set GOARCH=ArchitectureName
```
Architecture names can be:\
`arm64`\
`amd64`\
`arm`\
`386`




### Linux (cross compile)

Run this command in shell
```bash
$ cd src
$ env GOOS=windows GOARCH=amd64 go build -ldflags "-w -s"
```

Optionally, to add an icon to the exe, create and link a resource with
```bash
$ go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
$ export PATH=$PATH:~/go/bin
$ goversioninfo -icon res/DistroName/icon.ico -o src/DistroName.syso
$ env GOOS=windows GOARCH=ArchitectureName go build -ldflags "-w -s" -o DistroName.exe
```
