# Developers Document

## 🛠How-to-Build

### Windows

#### Visual Studio or Build Tools 2017+

Use `x64 Native Tools Command Prompt for VS 2017` or run this in the Windows Command Prompt (replace `2017` with `2019` for VS 2019)
```cmd
:: initialize x64 build environment
CALL "%ProgramFiles(x86)%\\Microsoft Visual Studio\\2017\\Community\\VC\\Auxiliary\\Build\\vcvarsall.bat" x64
```

Generate wslapi.lib
```cmd
lib /nologo /def:wslapi.def /machine:X64
```

To compile Launcher.exe
```cmd
cl /nologo /O2 /W4 /WX /Ob2 /Oi /Oy /Gs- /GF /Gy /Tc main.c /Fe:Launcher.exe Advapi32.lib Shell32.lib shlwapi.lib wslapi.lib
```

Optionally, to add an icon to the exe, create and link a resource with
```cmd
SET YourDistroName=Fedora

:: create resources
rc /nologo res\%YourDistroName%\res.rc

:: compile to %YourDistroName%.exe
cl /nologo /O2 /W4 /WX /Ob2 /Oi /Oy /Gs- /GF /Gy /Tc main.c /Fe:%YourDistroName%.exe ^
  Advapi32.lib Shell32.lib shlwapi.lib wslapi.lib res\%YourDistroName%\res.res
```

### MinGW
Install x86_64 version of MSYS2(https://www.msys2.org).

Run these commands in msys shell
```bash
$ pacman -S mingw-w64-x86_64-toolchain # install tool chain
$ gcc -std=c99 --static main.c -lshlwapi -lwslapi -o Launcher.exe # compile main.c
```

Optionally, to add an icon to the exe, create and link a resource with
```bash
$ YourDistroName=Fedora
$ windres res/$YourDistroName/res.rc res.o # compile resource
$ gcc -std=c99 --static main.c -lshlwapi -lwslapi -o ${YourDistroName}.exe res.o # compile main.c
```

### Linux (cross compile)
Install mingw-w64 toolchain include gcc-mingw-w64-x86-64.

Run this command in shell
```bash
 $ x86_64-w64-mingw32-gcc -std=c99 --static main.c -lshlwapi -lwslapi -o Launcher.exe # compile main.c
```
