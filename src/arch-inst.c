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

wchar_t *GetLxUID(wchar_t *DistributionName,wchar_t *LxUID);


int main()
{
    //Set Target Name
    wchar_t TargetName[] = L"Arch";


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

    if(IsDistributionRegistered(TargetName))
    {
        wchar_t LxUID[50] = L"";
        if(GetLxUID(TargetName,LxUID) != NULL)
        {
            wchar_t wcmd[70] = L"wsl.exe ";
            wcscat(wcmd,LxUID);
            int res = _wsystem(wcmd);//Excute wsl with LxUID
            return res;
        }
        else
        {
            printf("ERROR:GetLxUID failed!");
            return 1;
        }
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