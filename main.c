/*
 * Copyright (c) 2017 yuk7
 * Author: yuk7 <yukx00@gmail.com>
 *
 * Released under the MIT license
 * http://opensource.org/licenses/mit-license.php
 */


#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <wchar.h>
#include <windows.h>

#define ARRAY_LENGTH(a) (sizeof(a)/sizeof(a[0]))

typedef int (WINAPI *ISDISTRIBUTIONREBISTERED)(PCWSTR);
typedef int (WINAPI *REGISTERDISTRIBUTION)(PCWSTR,PCWSTR);
typedef int (WINAPI *CONFIGUREDISTRIBUTION)(PCWSTR,ULONG,INT);
typedef int (WINAPI *GETDISTRIBUTIONCONFIGURATION)(PCWSTR,ULONG*,ULONG*,INT*,PSTR*,ULONG*);
typedef int (WINAPI *LAUNCHINTERACTIVE)(PCWSTR,PCWSTR,INT,DWORD*);

wchar_t *GetLxUID(wchar_t *DistributionName,wchar_t *LxUID);


int main(int argc,char *argv[])
{
    int res = 0;
    wchar_t **wargv;
    int wargc;
    wargv = CommandLineToArgvW(GetCommandLineW(),&wargc);

    //Get file name of exe
    wchar_t efpath[300];
    if(GetModuleFileNameW(NULL,efpath,ARRAY_LENGTH(efpath)-1) == 0)
        return 1;
    wchar_t TargetName[300];
    _wsplitpath(efpath,NULL,NULL,TargetName,NULL);


    HMODULE hmod;
    ISDISTRIBUTIONREBISTERED IsDistributionRegistered;
    REGISTERDISTRIBUTION RegisterDistribution;
    CONFIGUREDISTRIBUTION ConfigureDistribution;
    GETDISTRIBUTIONCONFIGURATION GetDistributionConfiguration;
    LAUNCHINTERACTIVE LaunchInteractive;

    hmod = LoadLibraryW(L"wslapi.dll");
    if (hmod == NULL) {
        fwprintf(stderr,L"ERROR:wslapi.dll load failed\n");
        return 1;
    }

    IsDistributionRegistered = (ISDISTRIBUTIONREBISTERED)GetProcAddress(hmod, "WslIsDistributionRegistered");
    RegisterDistribution = (REGISTERDISTRIBUTION)GetProcAddress(hmod, "WslRegisterDistribution");
    ConfigureDistribution = (CONFIGUREDISTRIBUTION)GetProcAddress(hmod, "WslConfigureDistribution");
    GetDistributionConfiguration = (GETDISTRIBUTIONCONFIGURATION)GetProcAddress(hmod, "WslGetDistributionConfiguration");
    LaunchInteractive = (LAUNCHINTERACTIVE)GetProcAddress(hmod, "WslLaunchInteractive");
    if (IsDistributionRegistered == NULL | RegisterDistribution == NULL | ConfigureDistribution == NULL 
        | GetDistributionConfiguration == NULL | LaunchInteractive ==NULL) {
        FreeLibrary(hmod);
        fwprintf(stderr,L"ERROR:GetProcAddress failed\n");
        return 1;
    }


    if(IsDistributionRegistered(TargetName))
    {
        unsigned long distributionVersion;
        unsigned long defaultUID;
        int distributionFlags;
        LPSTR defaultEnv;
        unsigned long defaultEnvCnt;
        res = GetDistributionConfiguration(TargetName,&distributionVersion,&defaultUID,&distributionFlags,&defaultEnv,&defaultEnvCnt);
        if(res!=0)
        {
            fwprintf(stderr,L"ERROR:Get Configuration failed! 0x%x",res);
        }

        wchar_t LxUID[50] = L"";
        if(GetLxUID(TargetName,LxUID) == NULL)
        {
            fwprintf(stderr,L"ERROR:GetLxUID failed!");
            return 1;
        }



        if(wargc >1)
        {
            if(wcscmp(wargv[1],L"run") == 0)
            {
            }
            else if((wcscmp(wargv[1],L"config") == 0)&&wargc>3)
            {
                if(wcscmp(wargv[2],L"--default-user") == 0)
                {
                    (void) ConfigureDistribution(TargetName,0,distributionFlags); //set default uid to 0(root)
                    FILE *fp;
                    unsigned long uid;
                    wchar_t wcmd[300] = L"wsl.exe ";
                    wcscat_s(wcmd,ARRAY_LENGTH(wcmd),LxUID);
                    wcscat_s(wcmd,ARRAY_LENGTH(wcmd),L"id -u ");
                    wcscat_s(wcmd,ARRAY_LENGTH(wcmd),wargv[3]);
                    if((fp=_wpopen(wcmd,L"r")) ==NULL) {
                        fwprintf(stderr,L"ERROR:Command Excute Failed!");
                        return 1;
                    }
                    wchar_t buf[256];
                    if (!feof(fp))
                        fgetws(buf, sizeof(buf), fp);

                    (void) pclose(fp);
                    if(swscanf(buf,L"%d",&uid)==1)
                    {
                        res = ConfigureDistribution(TargetName,uid,distributionFlags);
                        if(res != 0)
                        {
                            (void) ConfigureDistribution(TargetName,defaultUID,distributionFlags); //revert uid
                            fwprintf(stderr,L"ERROR:Configure Failed! 0x%x",res);
                            return 1;
                        }
                        return 0;
                    }
                    else
                    {
                        (void) ConfigureDistribution(TargetName,defaultUID,distributionFlags); //revert uid
                        wprintf(L"\n");
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nFailed to detect user.");
                    }
                    return 1;
                }
                else if(wcscmp(wargv[2],L"--default-uid") == 0)
                {
                    unsigned long uid;
                    if(swscanf(wargv[3],L"%d",&uid)==1)
                    {
                        res = ConfigureDistribution(TargetName,uid,distributionFlags);
                        if(res != 0)
                        {
                            fwprintf(stderr,L"ERROR:Configure Failed! 0x%x",res);
                            return 1;
                        }
                        return 0;
                    }
                    else
                    {
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nInput UID");
                    }
                    return 1;
                }
                else if(wcscmp(wargv[2],L"--append-path") == 0)
                {
                    if(wcscmp(wargv[3],L"on") == 0)
                        distributionFlags |= 0x2;
                    else if(wcscmp(wargv[3],L"off") == 0)
                        distributionFlags &= ~0x2;
                    else
                    {
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nInput on/off");
                        return 1;
                    }
                    res = ConfigureDistribution(TargetName,defaultUID,distributionFlags);
                    if(res != 0)
                    {
                        fwprintf(stderr,L"ERROR:Configure Failed! 0x%x",res);
                        return 1;
                    }
                    return 0;
                }
                else if(wcscmp(wargv[2],L"--mount-drive") == 0)
                {
                    if(wcscmp(wargv[3],L"on") == 0)
                        distributionFlags |= 0x4;
                    else if(wcscmp(wargv[3],L"off") == 0)
                        distributionFlags &= ~0x4;
                    else
                    {
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nInput on/off");
                        return 1;
                    }
                    res = ConfigureDistribution(TargetName,defaultUID,distributionFlags);
                    if(res != 0)
                    {
                        fwprintf(stderr,L"ERROR:Configure Failed! 0x%x",res);
                        return 1;
                    }
                    return 0;
                }
                else
                {
                    fwprintf(stderr,L"ERROR:Invalid Arguments");
                    return 1;
                }
            }
            else if((wcscmp(wargv[1],L"get") == 0)&&wargc>2)
            {
                if(wcscmp(wargv[2],L"--default-uid") == 0)
                {
                    wprintf(L"%d",defaultUID);
                    return 0;
                }
                if(wcscmp(wargv[2],L"--append-path") == 0)
                {
                    if(distributionFlags & 0x2)
                        wprintf(L"on");
                    else
                        wprintf(L"off");
                    return 0;
                }
                if(wcscmp(wargv[2],L"--mount-drive") == 0)
                {
                    if(distributionFlags & 0x4)
                    wprintf(L"on");
                else
                    wprintf(L"off");
                return 0;
                }
                if(wcscmp(wargv[2],L"--lxuid") == 0)
                {
                    wprintf(L"%s",LxUID);
                    return 0;
                }
                fwprintf(stderr,L"ERROR:Invalid Arguments");
                return 1;
            }
            else
            {
                fwprintf(stderr,L"ERROR:Invalid Arguments.");
                wprintf(L"\n\n");
                wprintf(L"Useage :\n");
                wprintf(L"    <no args>\n");
                wprintf(L"      - Launches the distro's default behavior. By default, this launches your default shell.\n\n");
                wprintf(L"    run <command line>\n");
                wprintf(L"      - Run the given command line in that distro.\n\n");
                wprintf(L"    config [setting [value]]\n");
                wprintf(L"      - `--default-user <user>`: Set the default user for this distro to <user>\n");
                wprintf(L"      - `--default-uid <uid>`: Set the default user uid for this distro to <uid>\n");
                wprintf(L"      - `--append-path <on|off>`: Switch of Append Windows PATH to $PATH\n");
                wprintf(L"      - `--mount-drive <on|off>`: Switch of Mount drives\n\n");
                wprintf(L"    get [setting]\n");
                wprintf(L"      - `--default-uid`: Get the default user uid in this distro\n");
                wprintf(L"      - `--append-path`: Get on/off status of Append Windows PATH to $PATH\n");
                wprintf(L"      - `--mount-drive`: Get on/off status of Mount drives\n");
                wprintf(L"      - `--lxuid`: Get LxUID key for this distro\n\n");

                return 1;
            }
        }

        wchar_t rArgs[300] = L"";
        for (int i=2;i<wargc;i++)
        {
            wcscat_s(rArgs,ARRAY_LENGTH(rArgs),L" ");
            wcscat_s(rArgs,ARRAY_LENGTH(rArgs),wargv[i]);
        }

        unsigned long exitcode;
        res = LaunchInteractive(TargetName,rArgs,1,&exitcode);
        if(res==0)
            return exitcode;
        else
        {
            fwprintf(stderr,L"ERROR:Launch Interactive mode Failed! 0x%x",res);
            return 1;
        }
    }
    else
    {
        if(wargc >1)
        {
            fwprintf(stderr,L"ERROR:[%s] is not installed.\nRun with no arguments to install",TargetName);
            return 1;
        }
        wprintf(L"Installing...\n");
        res = RegisterDistribution(TargetName,L"rootfs.tar.gz");
        if(res != 0)
        {
            fwprintf(stderr,L"ERROR:Installation Failed! 0x%x",res);
            return 1;
        }
        wprintf(L"Installation Complete!");
        return 0;
    }
    return 0;
}

wchar_t *GetLxUID(wchar_t *DistributionName,wchar_t *LxUID)
{
    wchar_t RKey[]=L"Software\\Microsoft\\Windows\\CurrentVersion\\Lxss";
    HKEY hKey;
    LONG rres;
    if(RegOpenKeyExW(HKEY_CURRENT_USER,RKey, 0, KEY_READ, &hKey) == ERROR_SUCCESS)
    {
        for(int i=0;;i++)
        {
            wchar_t subKey[200];
            wchar_t subKeyF[200];
            wcscpy_s(subKeyF,ARRAY_LENGTH(subKeyF),RKey);
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
            wcscat_s(subKeyF,ARRAY_LENGTH(subKeyF),L"\\");
            wcscat_s(subKeyF,ARRAY_LENGTH(subKeyF),subKey);
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