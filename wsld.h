/*
 * Copyright (c) 2017-2020 yuk7
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

#define LXSS_BASE_RKEY L"Software\\Microsoft\\Windows\\CurrentVersion\\Lxss"
#define WSLDL_TERM_KEY L"wsldl-term"
#define MAX_DISTRO_NAME_SIZE 50
#define MAX_BASEPATH_SIZE 128
#define UUID_SIZE 38

struct WslInstallation {
    wchar_t uuid[UUID_SIZE+1];
    wchar_t basePath[MAX_BASEPATH_SIZE];
} WslInstallation;

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

        BOOL isWow64 = FALSE;
        if (IsWow64Process(GetCurrentProcess(), &isWow64))
        {
            if (isWow64)
            {
                fwprintf(stderr,L"       wslapi.dll is only known to support 64-bit executables but this executable is not.\n");
            }
        }
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
    if (WslIsDistributionRegistered == NULL || WslRegisterDistribution == NULL || WslUnregisterDistribution == NULL
        || WslConfigureDistribution == NULL || WslGetDistributionConfiguration == NULL || WslLaunchInteractive == NULL || WslLaunch == NULL)
    {
        FreeLibrary(WslHmod);
        fwprintf(stderr,L"ERROR: GetProcAddress() failed to get function address\n");
        wprintf(L"Press any key to exit...");
        getchar();
        exit(EXIT_FAILURE);
    }
return 0;
}

struct WslInstallation WslGetInstallationInfo(wchar_t *DistributionName) {
    struct WslInstallation wslInstallation = {.uuid = {0}, .basePath = {0}};

    HKEY hKey;
    LONG rres;
    if(RegOpenKeyExW(HKEY_CURRENT_USER,LXSS_BASE_RKEY, 0, KEY_READ, &hKey) == ERROR_SUCCESS)
    {
        int i;
        for(i=0;;i++)
        {
            wchar_t subKeyF[200];
            wcscpy_s(subKeyF,(sizeof(subKeyF)/sizeof(subKeyF[0])),LXSS_BASE_RKEY);

            wchar_t subKey[200];
            DWORD subKeySz = 100;
            FILETIME ftLastWriteTime;
            rres = RegEnumKeyExW(hKey, i, subKey, &subKeySz, NULL, NULL, NULL, &ftLastWriteTime);
            if (rres == ERROR_NO_MORE_ITEMS)
                break;
            else if(rres != ERROR_SUCCESS)
            {
                return wslInstallation;
            }

            DWORD dwType;
            HKEY hKeyS;
            wcscat_s(subKeyF,(sizeof(subKeyF)/sizeof(subKeyF[0])),L"\\");
            wcscat_s(subKeyF,(sizeof(subKeyF)/sizeof(subKeyF[0])),subKey);
            RegOpenKeyExW(HKEY_CURRENT_USER,subKeyF, 0, KEY_READ, &hKeyS);

            wchar_t regDistName[MAX_DISTRO_NAME_SIZE*2];
            DWORD dwSize = MAX_DISTRO_NAME_SIZE;
            rres = RegQueryValueExW(hKeyS, L"DistributionName", NULL, &dwType, (LPBYTE)&regDistName,&dwSize);
            if (rres != ERROR_SUCCESS)
            {
                // TODO: this helps for diagnostic, but we should implement a better error handling in the future
                fwprintf(stderr,L"ERROR:[%i] Could not read registry key\n", rres);
            }
            if((subKeySz == UUID_SIZE) && (_wcsicmp(regDistName,DistributionName)==0))
            {
                // SUCCESS: Distribution found
                wcscpy_s(wslInstallation.uuid, UUID_SIZE*2, subKey);
                DWORD pathSize = MAX_BASEPATH_SIZE*2;
                rres = RegQueryValueExW(hKeyS, L"BasePath", NULL, &dwType, (LPBYTE)&wslInstallation.basePath, &pathSize);
                if (rres != ERROR_SUCCESS)
                {
                    fwprintf(stderr,L"ERROR:[%i] Could not read registry key\n", rres);
                }

                RegCloseKey(hKey);
                RegCloseKey(hKeyS);
                return wslInstallation;
            }
            RegCloseKey(hKeyS);
        }
    }
    RegCloseKey(hKey);

    return wslInstallation;
}

 int WslSetInstallationPathInfo(wchar_t *uuid,wchar_t *basePath)
 {
    LONG rres;
    HKEY hKey;
    wchar_t RKey[200];
    wcscpy_s(RKey,(sizeof(RKey)/sizeof(RKey[0])),LXSS_BASE_RKEY);
    wcscat_s(RKey,(sizeof(RKey)/sizeof(RKey[0])),L"\\");
    wcscat_s(RKey,(sizeof(RKey)/sizeof(RKey[0])),uuid);

    RegOpenKeyExW(HKEY_CURRENT_USER,RKey, 0, KEY_SET_VALUE, &hKey);
    rres = RegSetValueExW(hKey,L"BasePath",0,REG_SZ,(const BYTE*)basePath,MAX_BASEPATH_SIZE);
    if (rres != ERROR_SUCCESS)
    {
        fwprintf(stderr,L"ERROR:[%i] Write registory key failed\n", rres);
        return 1;
    }
    return 0;
 }

int WsldlSetTerminalInfo(wchar_t *uuid ,long id)
{
    LONG rres;
    HKEY hKey;
    wchar_t RKey[200];
    wcscpy_s(RKey,(sizeof(RKey)/sizeof(RKey[0])),LXSS_BASE_RKEY);
    wcscat_s(RKey,(sizeof(RKey)/sizeof(RKey[0])),L"\\");
    wcscat_s(RKey,(sizeof(RKey)/sizeof(RKey[0])),uuid);

    RegOpenKeyExW(HKEY_CURRENT_USER,RKey, 0, KEY_SET_VALUE, &hKey);
    rres = RegSetValueExW(hKey, WSLDL_TERM_KEY, 0, REG_DWORD, (const BYTE*)&id, (int)sizeof(id));
    return rres;
}

 unsigned long WslExec(wchar_t *DistroName, wchar_t *command, char *result, long unsigned int *len)
 {
     HANDLE hProcess;
     HANDLE hOutTmp,hOut;
     HANDLE hInTmp,hIn;
     SECURITY_ATTRIBUTES sa;
     sa.nLength = sizeof(sa);
     sa.bInheritHandle = TRUE;
     sa.lpSecurityDescriptor = NULL;
    
    CreatePipe(&hOut, &hOutTmp, &sa, 0);
    CreatePipe(&hIn, &hInTmp, &sa, 0);
    if(WslLaunch(DistroName, command, 1, hInTmp, hOutTmp, hOutTmp, &hProcess))
    {
        return 100000;
    }
    CloseHandle(hInTmp);
    CloseHandle(hOutTmp);

    WaitForSingleObject(hProcess, INFINITE);
    unsigned long exitcode;
    GetExitCodeProcess(hProcess, &exitcode);

    if(!ReadFile(hOut, result, *len, len, NULL))
    {
        return 200000;
    }

    if(result[(strrchr(result, '\0') - result) - 1] == '\n')
    {
        result[(strrchr(result, '\0') - result) - 1] = '\0';
    }
    
    CloseHandle(hIn);
    CloseHandle(hOut);
    CloseHandle(hProcess);

    return exitcode;
 }

#ifdef __cplusplus
}
#endif

#endif
