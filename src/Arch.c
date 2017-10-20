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
    //Get file name of exe
    char efpath[300];
    if(GetModuleFileName(NULL,efpath,300) == 0)
        return 1;
    char efName[50];
    _splitpath(efpath,NULL,NULL,efName,NULL);

    wchar_t TargetName[30];
    mbstowcs_s(NULL,TargetName,30,efName,_TRUNCATE);


    HMODULE hmod;
    ISDISTRIBUTIONREBISTERED IsDistributionRegistered;
    REGISTERDISTRIBUTION RegisterDistribution;
    CONFIGUREDISTRIBUTION ConfigureDistribution;

    hmod = LoadLibrary(TEXT("wslapi.dll"));
    if (hmod == NULL) {
        printf("ERROR:wslapi.dll load failed\n");
        return 1;
    }

    IsDistributionRegistered = (ISDISTRIBUTIONREBISTERED)GetProcAddress(hmod, "WslIsDistributionRegistered");
    RegisterDistribution = (REGISTERDISTRIBUTION)GetProcAddress(hmod, "WslRegisterDistribution");
    ConfigureDistribution = (CONFIGUREDISTRIBUTION)GetProcAddress(hmod, "WslConfigureDistribution");
    if (IsDistributionRegistered == NULL | RegisterDistribution == NULL | ConfigureDistribution == NULL) {
        FreeLibrary(hmod);
        printf("ERROR:GetProcAddress failed\n");
        return 1;
    }


    if(IsDistributionRegistered(TargetName))
    {
        if(argc >1)
        {
            if(strcmp(argv[1],"run") == 0)
            {
            }
            else if((strcmp(argv[1],"config") == 0)&&argc>3)
            {
                if(strcmp(argv[2],"--default-uid") == 0)
                {
                    long uid;
                    if(sscanf(argv[3],"%d",&uid)==1)
                    {
                        int a = ConfigureDistribution(TargetName,uid,0x7);
                        if(a != 0)
                        {
                            printf("ERROR:Configure Failed! 0x%x",a);
                            return 1;
                        }
                        return 0;
                    }
                    else
                    {
                        printf("ERROR:Invalid Argument.\nInput UID");

                    }
                    return 1;
                }
                else
                {
                    printf("ERROR:Invalid Arguments");
                    return 1;
                }
            }
            else
            {
                printf("ERROR:Invalid Arguments\n\n");
                printf("Useage :\n");
                printf("    <no args>\n");
                printf("      - Launches the distro's default behavior. By default, this launches your default shell.\n\n");
                printf("    run <command line>\n");
                printf("      - Run the given command line in that distro.\n\n");
                printf("    config [setting [value]]\n");
                printf("      - `--default-uid <uid>`: Set the default user uid for this distro to <uid>\n\n");

                return 1;
            }
        }

        char rArgs[100] = "";
        for (int i=2;i<argc;i++)
        {
            strcat(rArgs," ");
            strcat(rArgs,argv[i]);
        }
        wchar_t wRcmd[100] = L"";
        mbstowcs_s(NULL,wRcmd,100,rArgs,_TRUNCATE);
        wchar_t LxUID[50] = L"";
        if(GetLxUID(TargetName,LxUID) != NULL)
        {
            wchar_t wcmd[120] = L"wsl.exe ";
            wcscat(wcmd,LxUID);
            wcscat(wcmd,wRcmd);
            int res = _wsystem(wcmd);//Excute wsl with LxUID
            return res;
        }
        else
        {
            printf("ERROR:GetLxUID failed!");
            return 1;
        }
    }
    else
    {
        if(argc >1)
        {
            wprintf(L"ERROR:[%s] is not installed.\nRun with no arguments to install",TargetName);
            return 1;
        }
        printf("Installing...\n\n");
        int a = RegisterDistribution(TargetName,L"rootfs.tar.gz");
        if(a != 0)
        {
            printf("ERROR:Installation Failed! 0x%x",a);
            return 1;
        }
        printf("Installation Complete!",a);
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
            if((subKeySz == 38)&&(strcmp(regDistName,DistributionName)==0))
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