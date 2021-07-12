:: Copyright (c) 2021 yuk7 <yukx00@gmail.com>
:: Released under the MIT license
:: http://opensource.org/licenses/mit-license.php

@echo off
cd /d %~dp0

if "%~1"=="all" (
    echo Build All
    call :resources
    call :icons
    call :single
    exit /b
)
if "%~1"=="resources" (
    echo Building Resources
    call :resources
    exit /b
)
if "%~1"=="icons" (
    echo Building binary with icons
    call :icons
    exit /b
)
if "%~1"=="clean" (
    echo Clean Files
    call :clean
    exit /b
)
call :single
exit /b



:resources
cd /d %~dp0
mkdir tools >NUL 2>&1
echo Downloading rsrc...
curl -sSfL https://github.com/akavel/rsrc/releases/download/v0.10.2/rsrc_windows_amd64.exe -o tools\rsrc.exe
echo Compiling All Resources...
FOR /D /r %%D in ("res/*") DO (
    tools\rsrc.exe -ico res\%%~nxD\icon.ico -o res\%%~nxD\res.syso
)
exit /b

:icons
cd /d %~dp0
echo Building with icons...
mkdir out\icons >NUL 2>&1
FOR /D /r %%D in ("res/*") DO (
    copy /y res\%%~nxD\res.syso src\res.syso
    cd src
    echo go build %GO_BUILD_OPTS% -o "%~dp0\out\icons\%%~nxD.exe"
    go build %GO_BUILD_OPTS% -o "%~dp0\out\icons\%%~nxD.exe"
    cd ..
    del /f src\res.syso
)
exit /b

:single
cd /d %~dp0
mkdir out >NUL 2>&1
cd src
echo running go build...
echo go build %GO_BUILD_OPTS% -o "%~dp0\out\wsldl.exe"
go build %GO_BUILD_OPTS% -o "%~dp0\out\wsldl.exe"
cd ..
:end
exit /b

:clean
FOR /D /r %%D in ("res/*") DO (
    cd /d %~dp0
    del res\%%~nxD\res.syso
)
rmdir /s /q out tools
exit /b