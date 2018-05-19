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

unsigned long QueryUser(wchar_t *TargetName,wchar_t *username);
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
        unsigned long distributionVersion;
        unsigned long defaultUID;
        int distributionFlags;
        LPSTR defaultEnv;
        unsigned long defaultEnvCnt;
        res = WslGetDistributionConfiguration(TargetName,&distributionVersion,&defaultUID,&distributionFlags,&defaultEnv,&defaultEnvCnt);

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
        else if(wcscmp(wargv[1],L"config") == 0)
        {
            if(wargc == 4)
            {
                if(wcscmp(wargv[2],L"--default-user") == 0)
                {
                    unsigned long uid;
                    uid = QueryUser(TargetName,wargv[3]);
                    if(uid != E_FAIL)
                    {
                        hr = WslConfigureDistribution(TargetName,uid,distributionFlags);
                    }
                }
                else if(wcscmp(wargv[2],L"--default-uid") == 0)
                {
                    unsigned long uid;
                    if(swscanf(wargv[3],L"%d",&uid)==1)
                    {
                        hr = WslConfigureDistribution(TargetName,uid,distributionFlags);
                    }
                    else
                    {
                        hr = E_INVALIDARG;
                    }
                }
                else if(wcscmp(wargv[2],L"--append-path") == 0)
                {
                    if(wcscmp(wargv[3],L"on") == 0)
                        distributionFlags |= 0x2;
                    else if(wcscmp(wargv[3],L"off") == 0)
                        distributionFlags &= ~0x2;
                    else
                    {
                        hr = E_INVALIDARG;
                    }
                    if(hr != E_INVALIDARG)
                    {
                        res = WslConfigureDistribution(TargetName,defaultUID,distributionFlags);
                    }
                }
                else if(wcscmp(wargv[2],L"--mount-drive") == 0)
                {
                    if(wcscmp(wargv[3],L"on") == 0)
                        distributionFlags |= 0x4;
                    else if(wcscmp(wargv[3],L"off") == 0)
                        distributionFlags &= ~0x4;
                    else
                    {
                        hr = E_INVALIDARG;
                    }
                    if(hr != E_INVALIDARG)
                    {
                        res = WslConfigureDistribution(TargetName,defaultUID,distributionFlags);
                    }
                }
                else
                {
                    hr = E_INVALIDARG;
                }
            }
            else
            {
                hr = E_INVALIDARG;
            }
        }
        else if(wcscmp(wargv[1],L"get") == 0)
        {
            if(wargc == 3)
            {
                if(wcscmp(wargv[2],L"--default-uid") == 0)
                {
                    wprintf(L"%d",defaultUID);
                    hr = S_OK;
                }
                else if(wcscmp(wargv[2],L"--append-path") == 0)
                {
                    if(distributionFlags & 0x2)
                        wprintf(L"on");
                    else
                        wprintf(L"off");
                    hr = S_OK;
                }
                else if(wcscmp(wargv[2],L"--mount-drive") == 0)
                {
                    if(distributionFlags & 0x4)
                        wprintf(L"on");
                    else
                        wprintf(L"off");
                hr = S_OK;
                }
                else if(wcscmp(wargv[2],L"--lxuid") == 0)
                {
                    wchar_t LxUID[50] = L"";
                    if(WslGetLxUID(TargetName,LxUID) == NULL)
                    {
                        hr = E_FAIL;
                    }
                    wprintf(L"%s",LxUID);
                    hr = S_OK;
                }
                else
                {
                    hr = E_INVALIDARG;
                }
            }
            else
            {
                hr = E_INVALIDARG;
            }
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

unsigned long QueryUser(wchar_t *TargetName,wchar_t *username)
{
    HANDLE hProcess;
    HANDLE hOutTmp,hOut;
    HANDLE hInTmp,hIn;
    SECURITY_ATTRIBUTES sa;
    sa.nLength = sizeof(sa);
    sa.bInheritHandle = TRUE;
    sa.lpSecurityDescriptor = NULL;
    unsigned long uid;
    wchar_t idcmd[30] = L"id -u ";
    wcscat_s(idcmd,ARRAY_LENGTH(idcmd),username);
    
    CreatePipe(&hOut, &hOutTmp, &sa, 0);
    CreatePipe(&hIn, &hInTmp, &sa, 0);
    if(WslLaunch(TargetName,idcmd,0,hInTmp,hOutTmp,hOutTmp,&hProcess))
    {
        fwprintf(stderr,L"ERROR:Failed to Excute id command.\n");
        return E_FAIL;
    }
    CloseHandle(hInTmp);
    CloseHandle(hOutTmp);

    char buf[300];
    DWORD len = 0;
    if(!ReadFile(hOut, &buf, sizeof(buf), &len, NULL))
    {
        fwprintf(stderr,L"ERROR:Failed to read result.\n");
        return E_FAIL;
    }
    
    CloseHandle(hInTmp);
    CloseHandle(hOutTmp);
    CloseHandle(hProcess);

    //read output
    if(sscanf(buf,"%d",&uid)==1)
    {
        return uid;
    }
    return E_FAIL;
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
    wprintf(L"    help\n");
    wprintf(L"      - Print this usage message.\n\n");
    
}
