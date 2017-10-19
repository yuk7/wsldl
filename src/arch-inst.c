/*
 * Copyright (c) 2017 yuk7
 * Author: yuk7 <yukx00@gmail.com>
 *
 * Released under the MIT license
 * http://opensource.org/licenses/mit-license.php
 */

#include <stdio.h>
#include <windows.h>

typedef int (WINAPI *REGISTERDISTRIBUTION)(PCWSTR,PCWSTR);

int main()
{
    HMODULE hmod;
    REGISTERDISTRIBUTION RegisterDistribution;

    hmod = LoadLibrary(TEXT("wslapi.dll"));
    if (hmod == NULL) {
        printf("ERROR:wslapi.dll load failed\n");
        return 1;
    }


    RegisterDistribution = (REGISTERDISTRIBUTION)GetProcAddress(hmod, "WslRegisterDistribution");
    if (RegisterDistribution == NULL) {
        FreeLibrary(hmod);
        printf("ERROR:GetProcAddress failed\n");
        return 1;
    }

    printf("Installing...\n\n");
    int a = RegisterDistribution(L"Arch",L"rootfs.tar.gz");
    if(a != 0)
    {
        printf("ERROR:Installation Failed! 0x%x",a);
        return 1;
    }
    printf("Installation Complete!",a);
    return 0;
}