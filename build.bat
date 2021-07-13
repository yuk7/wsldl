:: Copyright (c) 2021 yuk7 <yukx00@gmail.com>
:: Released under the MIT license
:: http://opensource.org/licenses/mit-license.php

@echo off
cd /d %~dp0

set PATH="%GOPATH%\bin";%PATH%
set PATH="%USERPROFILE%\go\bin";%PATH%

if "%~1"=="all" (
    echo Building everything
    call :resources
    call :icons
    call :single
    exit /b
)
if "%~1"=="resources" (
    echo Building resources
    call :resources
    exit /b
)
if "%~1"=="icons" (
    echo Building icon binaries
    call :icons
    exit /b
)
if "%~1"=="single" (
    echo Building binary...
    call :singlewor
    exit /b
)
if "%~1"=="singlewor" (
    echo Building binary without resource...
    call :singlewor
    exit /b
)
if "%~1"=="clean" (
    echo Removal of .syso files
    call :clean
    exit /b
)
call :single
exit /b



:resources
cd /d %~dp0
mkdir tools >NUL 2>&1
echo Installing goversioninfo...
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo

echo Compiling all resources...
FOR /D /r %%D in ("res/*") DO (
    goversioninfo -icon res\%%~nxD\icon.ico -o res\%%~nxD\resource.syso src\versioninfo.json
)
exit /b

:icons
cd /d %~dp0
echo Building wsldl with icons...
mkdir out\icons >NUL 2>&1
FOR /D /r %%D in ("res/*") DO (
    copy /y res\%%~nxD\resource.syso src\resource.syso
    cd src
    echo go build %GO_BUILD_OPTS% -o "%~dp0\out\icons\%%~nxD.exe"
    go build %GO_BUILD_OPTS% -o "%~dp0\out\icons\%%~nxD.exe"
    cd ..
    del /f src\resource.syso
)
exit /b

:single
cd /d %~dp0
mkdir out >NUL 2>&1
echo Installing goversioninfo...
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo

echo Compiling resource object...
goversioninfo -o src\resource.syso src\versioninfo.json
:singlewor
cd src
echo Building default wsldl.exe...
echo go build %GO_BUILD_OPTS% -o "%~dp0\out\wsldl.exe"
go build %GO_BUILD_OPTS% -o "%~dp0\out\wsldl.exe"
cd ..
:end
exit /b

:clean
FOR /D /r %%D in ("res/*") DO (
    cd /d %~dp0
    del res\%%~nxD\resource.syso
)
del src\resource.syso
rmdir /s /q out tools
exit /b
