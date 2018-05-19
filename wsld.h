/*
 * Copyright (c) 2017 yuk7
 * Author: yuk7 <yukx00@gmail.com>
 *
 * Released under the MIT license
 * http://opensource.org/licenses/mit-license.php
 */

 #ifndef WSLD_H_
 #define WSLD_H_

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <wchar.h>
#include <windows.h>

#ifdef __cplusplus
extern "C" {
#endif



typedef HRESULT (WINAPI *WSLISDISTRIBUTIONREBISTERED)(PCWSTR);
typedef HRESULT (WINAPI *WSLREGISTERDISTRIBUTION)(PCWSTR,PCWSTR);
typedef HRESULT (WINAPI *WSLUNREGISTERDISTRIBUTION)(PCWSTR);
typedef HRESULT (WINAPI *WSLCONFIGUREDISTRIBUTION)(PCWSTR,ULONG,INT);
typedef HRESULT (WINAPI *WSLGETDISTRIBUTIONCONFIGURATION)(PCWSTR,ULONG*,ULONG*,INT*,PSTR*,ULONG*);
typedef HRESULT (WINAPI *WSLLAUNCHINTERACTIVE)(PCWSTR,PCWSTR,BOOL,DWORD*);
typedef HRESULT (WINAPI *WSLLAUNCH)(PCWSTR,PCWSTR,BOOL,HANDLE,HANDLE,HANDLE,HANDLE*);

HMODULE WslHmod;
WSLISDISTRIBUTIONREBISTERED WslIsDistributionRegistered;
WSLREGISTERDISTRIBUTION WslRegisterDistribution;
WSLUNREGISTERDISTRIBUTION WslUnregisterDistribution;
WSLCONFIGUREDISTRIBUTION WslConfigureDistribution;
WSLGETDISTRIBUTIONCONFIGURATION WslGetDistributionConfiguration;
WSLLAUNCHINTERACTIVE WslLaunchInteractive;
WSLLAUNCH WslLaunch;

void WslApiFree()
{
    FreeLibrary(WslHmod);
}

int WslApiInit()
{
    WslHmod = LoadLibraryExW(L"wslapi.dll", NULL, LOAD_LIBRARY_SEARCH_SYSTEM32);
    if (WslHmod == NULL)
    {
        fwprintf(stderr,L"ERROR: LoadLibraryEx() failed to load wslapi.dll\n");
        wprintf(L"Press any key to exit...");
        getchar();
        exit(EXIT_FAILURE);
    }

    WslIsDistributionRegistered = (WSLISDISTRIBUTIONREBISTERED)GetProcAddress(WslHmod, "WslIsDistributionRegistered");
    WslRegisterDistribution = (WSLREGISTERDISTRIBUTION)GetProcAddress(WslHmod, "WslRegisterDistribution");
    WslUnregisterDistribution = (WSLUNREGISTERDISTRIBUTION)GetProcAddress(WslHmod, "WslUnregisterDistribution");
    WslConfigureDistribution = (WSLCONFIGUREDISTRIBUTION)GetProcAddress(WslHmod, "WslConfigureDistribution");
    WslGetDistributionConfiguration = (WSLGETDISTRIBUTIONCONFIGURATION)GetProcAddress(WslHmod, "WslGetDistributionConfiguration");
    WslLaunchInteractive = (WSLLAUNCHINTERACTIVE)GetProcAddress(WslHmod, "WslLaunchInteractive");
    WslLaunch = (WSLLAUNCH)GetProcAddress(WslHmod, "WslLaunch");
    if (WslIsDistributionRegistered == NULL | WslRegisterDistribution == NULL | WslUnregisterDistribution == NULL
    | WslConfigureDistribution == NULL | WslGetDistributionConfiguration == NULL | WslLaunchInteractive == NULL | WslLaunch == NULL)
    {
        FreeLibrary(WslHmod);
        fwprintf(stderr,L"ERROR: GetProcAddress() failed to get function address\n");
        wprintf(L"Press any key to exit...");
        getchar();
        exit(EXIT_FAILURE);
    }
return 0;
}

wchar_t *WslGetLxUID(wchar_t *DistributionName,wchar_t *LxUID)
{
    wchar_t RKey[]=L"Software\\Microsoft\\Windows\\CurrentVersion\\Lxss";
    HKEY hKey;
    LONG rres;
    if(RegOpenKeyExW(HKEY_CURRENT_USER,RKey, 0, KEY_READ, &hKey) == ERROR_SUCCESS)
    {
        int i;
        for(i=0;;i++)
        {
            wchar_t subKey[200];
            wchar_t subKeyF[200];
            wcscpy_s(subKeyF,(sizeof(subKeyF)/sizeof(subKeyF[0])),RKey);
            wchar_t regDistName[100];
            DWORD subKeySz = 100;
            DWORD dwType;
            DWORD dwSize = 50;
            FILETIME ftLastWriteTime;

            rres = RegEnumKeyExW(hKey, i, subKey, &subKeySz, NULL, NULL, NULL, &ftLastWriteTime);
            if (rres == ERROR_NO_MORE_ITEMS)
                break;
            else if(rres != ERROR_SUCCESS)
            {
                //ERROR
                LxUID = NULL;
                return LxUID;
            }

            HKEY hKeyS;
            wcscat_s(subKeyF,(sizeof(subKeyF)/sizeof(subKeyF[0])),L"\\");
            wcscat_s(subKeyF,(sizeof(subKeyF)/sizeof(subKeyF[0])),subKey);
            RegOpenKeyExW(HKEY_CURRENT_USER,subKeyF, 0, KEY_READ, &hKeyS);
            RegQueryValueExW(hKeyS, L"DistributionName", NULL, &dwType, (LPBYTE)&regDistName,&dwSize);
            if((subKeySz == 38)&&(wcscmp(regDistName,DistributionName)==0))
            {
                //SUCCESS:Distribution found!
                //return LxUID
                RegCloseKey(hKey);
                RegCloseKey(hKeyS);
                wcscpy_s(LxUID,40,subKey);
                return LxUID;
            }
            RegCloseKey(hKeyS);
            }
        }
        else
        {
        //ERROR
        LxUID = NULL;
        return LxUID;
        }
    RegCloseKey(hKey);
    //ERROR:Distribution Not Found
    LxUID = NULL;
    return LxUID;
}

#ifdef __cplusplus
}
#endif

#endif
