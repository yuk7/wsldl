:: Copyright (c) 2021 yuk7 <yukx00@gmail.com>
:: Released under the MIT license
:: http://opensource.org/licenses/mit-license.php

@echo off
cd /d %~dp0

set PATH="%GOPATH%\bin";%PATH%
set PATH="%USERPROFILE%\go\bin";%PATH%

if not defined GOBIN (
    set GOBIN=go
)

if not defined GOARCH (
    echo GOARCH is not defined, detecting cpu architecture...
    if "%PROCESSOR_ARCHITECTURE%"=="x86" (
        echo 32bit 386 processor detected
        set GOARCH=386
    )
    if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
        echo 64bit amd64 processor detected
        set GOARCH=amd64
    )
    if "%PROCESSOR_ARCHITECTURE%"=="ARM" (
        echo 32bit arm processor detected
        set GOARCH=arm
    )
    if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
        echo 64bit arm64 processor detected
        set GOARCH=arm64
    )
)

if "%GOARCH%"=="386" (
    echo GOARCH is 386
    set GOVERSIONINFO_OPTS=
)
if "%GOARCH%"=="amd64" (
    echo GOARCH is amd64
    set GOVERSIONINFO_OPTS=-64
)
if "%GOARCH%"=="arm" (
    echo GOARCH is arm
    set GOVERSIONINFO_OPTS=-arm
)
if "%GOARCH%"=="arm64" (
    echo GOARCH is arm64
    set GOVERSIONINFO_OPTS=-arm -64
)


if "%~1"=="all" (
    echo Building everything
    call :dlgoversioninfo
    call :resources
    call :icons
    call :single
    exit /b
)
if "%~1"=="resources" (
    echo Building resources
    call :dlgoversioninfo
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
    call :dlgoversioninfo
    call :single
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
call :dlgoversioninfo
call :single
exit /b



:resources
cd /d %~dp0
echo Compiling all resources...
FOR /D /r %%D in ("res/*") DO (
    tools\goversioninfo %GOVERSIONINFO_OPTS% -icon res\%%~nxD\icon.ico -o res\%%~nxD\resource.syso src\versioninfo.json
)
exit /b

:icons
cd /d %~dp0
echo Building wsldl with icons...
mkdir out\icons >NUL 2>&1
FOR /D /r %%D in ("res/*") DO (
    copy /y res\%%~nxD\resource.syso src\resource.syso
    cd src
    echo %GOBIN% build %GO_BUILD_OPTS% -o "%~dp0\out\icons\%%~nxD.exe"
    %GOBIN% build %GO_BUILD_OPTS% -o "%~dp0\out\icons\%%~nxD.exe"
    cd ..
    del /f src\resource.syso
)
exit /b

:single
cd /d %~dp0
echo Compiling resource object...
tools\goversioninfo %GOVERSIONINFO_OPTS% -o src\resource.syso src\versioninfo.json
:singlewor
cd /d %~dp0
mkdir out >NUL 2>&1
cd src
echo Building default wsldl.exe...
echo %GOBIN% build %GO_BUILD_OPTS% -o "%~dp0\out\wsldl.exe"
%GOBIN% build %GO_BUILD_OPTS% -o "%~dp0\out\wsldl.exe"
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

:dlgoversioninfo
cd /d %~dp0
mkdir tools >NUL 2>&1
if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    echo Downaloding goversioninfo 386...
    curl -sSfL https://github.com/yuk7/goversioninfo/releases/download/v1.2.0-arm/goversioninfo_386.exe -o tools\goversioninfo.exe
)
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    echo Downaloding goversioninfo amd64...
    curl -sSfL https://github.com/yuk7/goversioninfo/releases/download/v1.2.0-arm/goversioninfo_amd64.exe -o tools\goversioninfo.exe
)
if "%PROCESSOR_ARCHITECTURE%"=="ARM" (
    curl -sSfL https://github.com/yuk7/goversioninfo/releases/download/v1.2.0-arm/goversioninfo_arm.exe -o tools\goversioninfo.exe
)
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    curl -sSfL https://github.com/yuk7/goversioninfo/releases/download/v1.2.0-arm/goversioninfo_arm64.exe -o tools\goversioninfo.exe
)

exit /b