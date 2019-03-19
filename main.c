/*
 * Copyright (c) 2017-2019 yuk7
 * Author: yuk7 <yukx00@gmail.com>
 *
 * Released under the MIT license
 * http://opensource.org/licenses/mit-license.php
 */


#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <wchar.h>
#include <windows.h>
#include "wsld.h"
#include "version.h"

#define ARRAY_LENGTH(a) (sizeof(a)/sizeof(a[0]))
#define WARGV_CMP(a,b) ((wargc>a)?wcscmp(wargv[a],b)==0:false)

unsigned long QueryUser(wchar_t *TargetName,wchar_t *username);
bool dirExists(const char* dirName);
int InstallDist(wchar_t *TargetName,wchar_t *tgzname);
HRESULT RemoveDist(wchar_t *TargetName);
void show_usage();
void show_version();

int main()
{
    HRESULT hr = E_FAIL;
    DWORD exitCode = 0;
    wchar_t **wargv;
    int wargc;
    wargv = CommandLineToArgvW(GetCommandLineW(),&wargc);

    if(WARGV_CMP(1,L"version"))
    {
        show_version();
        return 0;
    }

    //Get file name of exe
    wchar_t efpath[MAX_PATH];
    if(GetModuleFileNameW(NULL,efpath,ARRAY_LENGTH(efpath)-1) == 0)
        return 1;
    wchar_t TargetName[MAX_PATH];
    _wsplitpath_s(efpath,NULL,0,NULL,0,TargetName,MAX_PATH,NULL,0);

    WslApiInit();

    if(WARGV_CMP(1,L"isregd"))
    {
        return (!WslIsDistributionRegistered(TargetName));
    }

    if(!WslIsDistributionRegistered(TargetName))
    {
        bool InstSilent = false;
        wchar_t tgzname[MAX_PATH] = L"rootfs.tar.gz";
        if(wargc >1)
        {
            //"tgz" and "silent" will be discontinued in the future.
            if( WARGV_CMP(1,L"install") | WARGV_CMP(1,L"tgz") | WARGV_CMP(1,L"silent") )
            {
                InstSilent = true;
                if( (!WARGV_CMP(2,L"--root")) & (wargc>2) )
                {
                    wcscpy_s(tgzname,ARRAY_LENGTH(tgzname),wargv[2]);
                }
            }
            else
            {
                fwprintf(stderr,L"ERROR:[%s] is not installed.\nRun with no arguments to install\n",TargetName);
                return 1;
            }
        }
        if(InstSilent)
        {
            hr = WslRegisterDistribution(TargetName,tgzname);
        }
        else
        {
            hr = InstallDist(TargetName,tgzname);
        }
        return hr;
    }
    else
    {
        unsigned long distributionVersion;
        unsigned long defaultUID;
        int distributionFlags;
        LPSTR defaultEnv;
        unsigned long defaultEnvCnt;
        WslGetDistributionConfiguration(TargetName,&distributionVersion,&defaultUID,&distributionFlags,&defaultEnv,&defaultEnvCnt);

        if(wargc == 1)
        {
            struct WslInstallation wslInstallation = WslGetInstallationInfo(TargetName);
            char buffer[MAX_BASEPATH_SIZE];
            size_t *retSize = 0;
            wcstombs_s(retSize, buffer, MAX_BASEPATH_SIZE, wslInstallation.basePath, MAX_BASEPATH_SIZE);
            if (!dirExists(buffer))
            {
                fwprintf(stderr,L"Installation directory not found: %s.\nMake sure it exists or reinstall.\n",wslInstallation.basePath);
                hr = E_ABORT;
            }
            else
            {
                SetConsoleTitleW(TargetName);
                hr = WslLaunchInteractive(TargetName,L"", false, &exitCode);
            }
        }
        else if( WARGV_CMP(1,L"run") | WARGV_CMP(1,L"-c") | WARGV_CMP(1,L"/c") )
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
        else if(WARGV_CMP(1,L"config"))
        {
            if((WARGV_CMP(2,L"--default-user"))&(wargc == 4))
            {
                unsigned long uid;
                uid = QueryUser(TargetName,wargv[3]);
                if(uid != E_FAIL)
                {
                    hr = WslConfigureDistribution(TargetName,uid,distributionFlags);
                }
            }
            else if((WARGV_CMP(2,L"--default-uid"))&(wargc == 4))
            {
                unsigned long uid;
                if(swscanf_s(wargv[3],L"%d",&uid)==1)
                {
                    hr = WslConfigureDistribution(TargetName,uid,distributionFlags);
                }
                else
                {
                    hr = E_INVALIDARG;
                }
            }
            else if(WARGV_CMP(2,L"--append-path"))
            {
                if(WARGV_CMP(3,L"on"))
                    distributionFlags |= 0x2;
                else if(WARGV_CMP(3,L"off"))
                    distributionFlags &= ~0x2;
                else
                    hr = E_INVALIDARG;

                if(hr != E_INVALIDARG)
                {
                    hr = WslConfigureDistribution(TargetName,defaultUID,distributionFlags);
                }
            }
            else if(WARGV_CMP(2,L"--mount-drive"))
            {
                if(WARGV_CMP(3,L"on"))
                    distributionFlags |= 0x4;
                else if(WARGV_CMP(3,L"off"))
                    distributionFlags &= ~0x4;
                else
                    hr = E_INVALIDARG;

                if(hr != E_INVALIDARG)
                {
                    hr = WslConfigureDistribution(TargetName,defaultUID,distributionFlags);
                }
            }
            else
            {
                hr = E_INVALIDARG;
            }
        }
        else if(WARGV_CMP(1,L"get"))
        {
            if(WARGV_CMP(2,L"--default-uid"))
            {
                wprintf(L"%d",defaultUID);
                hr = S_OK;
            }
            else if(WARGV_CMP(2,L"--append-path"))
            {
                if(distributionFlags & 0x2)
                    wprintf(L"on");
                else
                    wprintf(L"off");
                hr = S_OK;
            }
            else if(WARGV_CMP(2,L"--mount-drive"))
            {
                if(distributionFlags & 0x4)
                    wprintf(L"on");
                else
                    wprintf(L"off");
                hr = S_OK;
            }
            else if(WARGV_CMP(2,L"--lxuid"))
            {
                struct WslInstallation wsl = WslGetInstallationInfo(TargetName);
                if(wsl.uuid == NULL)
                {
                    hr = E_FAIL;
                }
                wprintf(L"%.*s",UUID_SIZE,wsl.uuid);
                hr = S_OK;
            }
            else
            {
                hr = E_INVALIDARG;
            }
        }
        else if(WARGV_CMP(1,L"backup"))
        {
            if(distributionFlags & 0x4)
            {
                WslConfigureDistribution(TargetName,0,distributionFlags);
                wprintf(L"Running backup command.\n");
                wprintf(L"If a password is requested, please enter the root password.\n\n");
                hr = WslLaunchInteractive(TargetName,L"su root -c \'tar -zcpf backup.tar.gz --exclude \"mnt/*\" --exclude \"dev/*\" --exclude \"proc/*\" --exclude \"sys/*\" --exclude \"run/*\" /\'", true, &exitCode);
                WslConfigureDistribution(TargetName,defaultUID,distributionFlags);
            }
            else
            {
                fwprintf(stderr,L"ERROR:Mount drive feature is not enabled.\n");
                fwprintf(stderr,L"Please enable it and retry.\n");
                hr = E_FAIL;
            }
            
        }
        else if(WARGV_CMP(1,L"clean"))
        {
            if(WARGV_CMP(2,L"-y"))
            {
                hr = WslUnregisterDistribution(TargetName);
                return hr;
            }
            else
            {
                hr = RemoveDist(TargetName);
            }
        }
        else if( WARGV_CMP(1,L"help") | WARGV_CMP(1,L"-h") | WARGV_CMP(1,L"/h") )
        {
            show_usage();
            hr = S_OK;
        }
        else
        {
            hr = E_INVALIDARG;
        }


        if(SUCCEEDED(hr))
        {
            return exitCode;
        }
        // already controlled error cases
        else if(hr==E_ABORT)
        {
            wprintf(L"Press any key to continue...");
            getchar();
            return hr;
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

bool dirExists(const char* dirName)
{
    DWORD ftyp = GetFileAttributesA(dirName);
    return (ftyp != INVALID_FILE_ATTRIBUTES) && (ftyp & FILE_ATTRIBUTE_DIRECTORY);
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
        fwprintf(stderr,L"ERROR:Failed to execute id command.\n");
        return (unsigned long)E_FAIL;
    }
    CloseHandle(hInTmp);
    CloseHandle(hOutTmp);

    char buf[300];
    DWORD len = 0;
    if(!ReadFile(hOut, &buf, sizeof(buf), &len, NULL))
    {
        fwprintf(stderr,L"ERROR:Failed to read result.\n");
        return (unsigned long)E_FAIL;
    }
    
    CloseHandle(hInTmp);
    CloseHandle(hOutTmp);
    CloseHandle(hProcess);

    //read output
    if(sscanf_s(buf,"%d",&uid)==1)
    {
        return uid;
    }
    return (unsigned long)E_FAIL;
}

int InstallDist(wchar_t *TargetName,wchar_t *tgzname)
{
    wprintf(L"Installing...\n");
    HRESULT hr = WslRegisterDistribution(TargetName,tgzname);
    if(FAILED(hr))
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

HRESULT RemoveDist(wchar_t *TargetName)
{
    char yn;
    wprintf(L"This will remove this distro (%s) from the filesystem.\n",TargetName);
    wprintf(L"Are you sure you would like to proceed? (This cannot be undone)\n");
    wprintf(L"Type \"y\" to continue:");
    scanf_s("%c",&yn,1);
    if(yn == 'y')
    {
        wprintf(L"Unregistering...\n");
        HRESULT hr = WslUnregisterDistribution(TargetName);
        return hr;
    }
    else
    {
        fwprintf(stderr,L"Accepting is required to proceed.\n\n");
        return S_OK;
    }
}

void show_usage()
{
    wprintf(L"Usage :\n");
    wprintf(L"    <no args>\n");
    wprintf(L"      - Open a new shell with your default settings.\n\n");
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
    wprintf(L"    backup\n");
    wprintf(L"      - Output backup.tar.gz to the current directory using tar command.\n\n");
    wprintf(L"    clean\n");
    wprintf(L"      - Uninstall the distro.\n\n");
    wprintf(L"    help\n");
    wprintf(L"      - Print this usage message.\n\n");
    
}

void show_version()
{
    wprintf(L"%s, version %s\n",SW_NAME,SW_VER);
    wprintf(L"%s\n",SW_URL);
}
