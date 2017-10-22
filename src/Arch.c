/*
 * Copyright (c) 2017 yuk7
 * Author: yuk7 <yukx00@gmail.com>
 *
 * Released under the MIT license
 * http://opensource.org/licenses/mit-license.php
 */


#include <stdio.h>
#include <windows.h>

typedef int (WINAPI *ISDISTRIBUTIONREBISTERED)(PCWSTR);
typedef int (WINAPI *REGISTERDISTRIBUTION)(PCWSTR,PCWSTR);
typedef int (WINAPI *CONFIGUREDISTRIBUTION)(PCWSTR,ULONG,INT);

wchar_t *GetLxUID(wchar_t *DistributionName,wchar_t *LxUID);


int main(int argc,char *argv[])
{
    wchar_t **wargv;
    int wargc;
    wargv = CommandLineToArgvW(GetCommandLineW(),&wargc);

    //Get file name of exe
    wchar_t efpath[300];
    if(GetModuleFileNameW(NULL,efpath,300) == 0)
        return 1;
    wchar_t TargetName[50];
    _wsplitpath(efpath,NULL,NULL,TargetName,NULL);


    HMODULE hmod;
    ISDISTRIBUTIONREBISTERED IsDistributionRegistered;
    REGISTERDISTRIBUTION RegisterDistribution;
    CONFIGUREDISTRIBUTION ConfigureDistribution;

    hmod = LoadLibraryW(L"wslapi.dll");
    if (hmod == NULL) {
        fwprintf(stderr,L"ERROR:wslapi.dll load failed\n");
        return 1;
    }

    IsDistributionRegistered = (ISDISTRIBUTIONREBISTERED)GetProcAddress(hmod, "WslIsDistributionRegistered");
    RegisterDistribution = (REGISTERDISTRIBUTION)GetProcAddress(hmod, "WslRegisterDistribution");
    ConfigureDistribution = (CONFIGUREDISTRIBUTION)GetProcAddress(hmod, "WslConfigureDistribution");
    if (IsDistributionRegistered == NULL | RegisterDistribution == NULL | ConfigureDistribution == NULL) {
        FreeLibrary(hmod);
        fwprintf(stderr,L"ERROR:GetProcAddress failed\n");
        return 1;
    }


    if(IsDistributionRegistered(TargetName))
    {
        if(wargc >1)
        {
            if(wcscmp(wargv[1],L"run") == 0)
            {
            }
            else if((wcscmp(wargv[1],L"config") == 0)&&wargc>3)
            {
                if(wcscmp(wargv[2],L"--default-uid") == 0)
                {
                    long uid;
                    if(swscanf(wargv[3],L"%d",&uid)==1)
                    {
                        int a = ConfigureDistribution(TargetName,uid,0x7);
                        if(a != 0)
                        {
                            fwprintf(stderr,L"ERROR:Configure Failed! 0x%x",a);
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
                else
                {
                    fwprintf(stderr,L"ERROR:Invalid Arguments");
                    return 1;
                }
            }
            else
            {
                fwprintf(stderr,L"ERROR:Invalid Arguments\n\n");
                wprintf(L"Useage :\n");
                wprintf(L"    <no args>\n");
                wprintf(L"      - Launches the distro's default behavior. By default, this launches your default shell.\n\n");
                wprintf(L"    run <command line>\n");
                wprintf(L"      - Run the given command line in that distro.\n\n");
                wprintf(L"    config [setting [value]]\n");
                wprintf(L"      - `--default-uid <uid>`: Set the default user uid for this distro to <uid>\n\n");

                return 1;
            }
        }

        wchar_t rArgs[100] = L"";
        for (int i=2;i<wargc;i++)
        {
            wcscat(rArgs,L" ");
            wcscat(rArgs,wargv[i]);
        }
        wchar_t LxUID[50] = L"";
        if(GetLxUID(TargetName,LxUID) != NULL)
        {
            wchar_t wcmd[120] = L"wsl.exe ";
            wcscat(wcmd,LxUID);
            wcscat(wcmd,rArgs);
            int res = _wsystem(wcmd);//Excute wsl with LxUID
            return res;
        }
        else
        {
            fwprintf(stderr,L"ERROR:GetLxUID failed!");
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
        int a = RegisterDistribution(TargetName,L"rootfs.tar.gz");
        if(a != 0)
        {
            fwprintf(stderr,L"ERROR:Installation Failed! 0x%x",a);
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
            wcscpy(subKeyF,RKey);
            wchar_t regDistName[100];
            DWORD subKeySz = 100;
            DWORD dwType;
            DWORD dwSize;
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
            wcscat(subKeyF,L"\\");
            wcscat(subKeyF,subKey);
            RegOpenKeyExW(HKEY_CURRENT_USER,subKeyF, 0, KEY_READ, &hKeyS);
            RegQueryValueExW(hKeyS, L"DistributionName", NULL, &dwType, &regDistName,&dwSize);
            RegQueryValueExW(hKeyS, L"DistributionName", NULL, &dwType, &regDistName,&dwSize);
            if((subKeySz == 38)&&(wcscmp(regDistName,DistributionName)==0))
            {
                //SUCCESS:Distribution found!
                //return LxUID
                RegCloseKey(hKey);
                RegCloseKey(hKeyS);
                wcscpy(LxUID,subKey);
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