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

char *GetLxUID(char *DistributionName,char *LxUID);


int main()
{
    HMODULE hmod;
    ISDISTRIBUTIONREBISTERED IsDistributionRegistered;
    REGISTERDISTRIBUTION RegisterDistribution;

    hmod = LoadLibrary(TEXT("wslapi.dll"));
    if (hmod == NULL) {
        printf("ERROR:wslapi.dll load failed\n");
        return 1;
    }

    IsDistributionRegistered = (ISDISTRIBUTIONREBISTERED)GetProcAddress(hmod, "WslIsDistributionRegistered");
    if (IsDistributionRegistered == NULL) {
        FreeLibrary(hmod);
        printf("ERROR:GetProcAddress failed\n");
        return 1;
    }
    RegisterDistribution = (REGISTERDISTRIBUTION)GetProcAddress(hmod, "WslRegisterDistribution");
    if (RegisterDistribution == NULL) {
        FreeLibrary(hmod);
        printf("ERROR:GetProcAddress failed\n");
        return 1;
    }

    if(IsDistributionRegistered(L"Arch"))
    {
        char LxUID[50] = "";
        if(GetLxUID(L"Arch",LxUID) != NULL)
        {
            char wcmd[70] = "wsl.exe ";
            strcat(wcmd,LxUID);
            int res = system(wcmd);//Excute wsl with LxUID
            return res;
        }
        else
        {
            printf("ERROR:GetLxUID failed!");
            return 1;
        }
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

char *GetLxUID(char *DistributionName,char *LxUID)
{
    char RKey[]="Software\\Microsoft\\Windows\\CurrentVersion\\Lxss";
    HKEY hKey;
    LONG rres;
    if(RegOpenKeyEx(HKEY_CURRENT_USER,RKey, 0, KEY_READ, &hKey) == ERROR_SUCCESS)
    {
	    for(int i=0;;i++)
	    {
            char subKey[200];
	        char subKeyF[200];
	        strcpy(subKeyF,RKey);
	        char regDistName[100];
	        DWORD subKeySz = 100;
	        DWORD dwType;
	        DWORD dwSize;
	        FILETIME ftLastWriteTime;

	        rres = RegEnumKeyEx(hKey, i, subKey, &subKeySz, NULL, NULL, NULL, &ftLastWriteTime);
	        if (rres == ERROR_NO_MORE_ITEMS)
                break;
	        else if(rres != ERROR_SUCCESS)
	        {
	            //ERROR
	            LxUID = NULL;
	            return LxUID;
	        }

	        HKEY hKeyS;
            strcat(subKeyF,"\\");
            strcat(subKeyF,subKey);
	        RegOpenKeyEx(HKEY_CURRENT_USER,subKeyF, 0, KEY_READ, &hKeyS);
	        RegQueryValueEx(hKeyS, "DistributionName", NULL, &dwType, &regDistName,&dwSize);
	        RegQueryValueEx(hKeyS, "DistributionName", NULL, &dwType, &regDistName,&dwSize);
	        if((subKeySz == 38)&&(strcmp(regDistName,DistributionName)==0))
	        {
                //SUCCESS:Distribution found!
                //return LxUID
	            RegCloseKey(hKey);
	            RegCloseKey(hKeyS);
	            strcpy(LxUID,subKey);
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