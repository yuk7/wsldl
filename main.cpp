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
#include "wsld.h"

#define ARRAY_LENGTH(a) (sizeof(a)/sizeof(a[0]))

int InstallDist(wchar_t *TargetName,wchar_t *tgzname);
void show_usage();

int main()
{
    int res = 0;
    HRESULT hr = E_FAIL;
    DWORD exitCode = 0;
    wchar_t **wargv;
    int wargc;
    wargv = CommandLineToArgvW(GetCommandLineW(),&wargc);

    //Get file name of exe
    wchar_t efpath[MAX_PATH];
    if(GetModuleFileNameW(NULL,efpath,ARRAY_LENGTH(efpath)-1) == 0)
        return 1;
    wchar_t TargetName[MAX_PATH];
    _wsplitpath(efpath,NULL,NULL,TargetName,NULL);

    WslApiInit();

    if(!WslIsDistributionRegistered(TargetName))
    {
        wchar_t tgzname[MAX_PATH] = L"rootfs.tar.gz";
        if(wargc >2)
        {
            if(wcscmp(wargv[1],L"tgz")==0)
            {
                wcscpy_s(tgzname,ARRAY_LENGTH(tgzname),wargv[2]);
            }
            else
            {
                fwprintf(stderr,L"ERROR:[%s] is not installed.\nRun with no arguments to install\n",TargetName);
                return 1;
            }
        }
        else if(wargc >1)
        {
            fwprintf(stderr,L"ERROR:[%s] is not installed.\nRun with no arguments to install\n",TargetName);
            return 1;
        }
        hr = InstallDist(TargetName,tgzname);
        return hr;
    }
    else
    {
        if(wargc == 1)
        {
            hr = WslLaunchInteractive(TargetName,L"", false, &exitCode);
        }
        else if(wcscmp(wargv[1],L"run") == 0)
        {
            wchar_t rArgs[500] = L"";
            int i;
            for (i=2;i<wargc;i++)
            {
                wcscat_s(rArgs,ARRAY_LENGTH(rArgs),L" ");
                wcscat_s(rArgs,ARRAY_LENGTH(rArgs),wargv[i]);
            }
            hr = WslLaunchInteractive(TargetName,rArgs, true, &exitCode);
        }
        else if(wcscmp(wargv[1],L"help")==0)
        {
            show_usage();
            return hr;
        }
        else
        {
            fwprintf(stderr,L"ERROR:Invalid Argument.\n\n");
            show_usage();
            return hr;
        }


        if(SUCCEEDED(hr))
        {
            return exitCode;
        }
        else if(hr==E_INVALIDARG)
        {
            fwprintf(stderr,L"ERROR:Invalid Argument.\n\n");
            show_usage();
            return hr;
        }
        else
        {
            fwprintf(stderr,L"ERROR\nHRESULT:0x%x\n",hr);
            wprintf(L"Press any key to continue...");
            getchar();
            return hr;
        }
    }
}

int InstallDist(wchar_t *TargetName,wchar_t *tgzname)
{
    wprintf(L"Installing...\n");
    HRESULT hr = WslRegisterDistribution(TargetName,tgzname);
    if(SUCCEEDED(hr))
    {
        fwprintf(stderr,L"ERROR:Installation Failed!\nHRESULT:0x%x\n",hr);
        wprintf(L"Press any key to continue...");
        getchar();
        return hr;
    }
    wprintf(L"Installation Complete!\n");
    wprintf(L"Press any key to continue...");
    getchar();
    return 0;
}

void show_usage()
{
    wprintf(L"Useage :\n");
    wprintf(L"    <no args>\n");
    wprintf(L"      - Launches the distro's default behavior. By default, this launches your default shell.\n\n");
    wprintf(L"    run <command line>\n");
    wprintf(L"      - Run the given command line in that distro. Inherit current directory.\n\n");
    wprintf(L"    help\n");
    wprintf(L"      - Print this usage message.\n\n");
    
}
