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


int main()
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


    res = WslApiInit();
    if (res) {
        fwprintf(stderr,L"ERROR:WslApi.dll load failed(%s).\n",res);
        wprintf(L"Please any key to continue...");
        getchar();
        return 1;
    }


    if(WslIsDistributionRegistered(TargetName))
    {
        unsigned long distributionVersion;
        unsigned long defaultUID;
        int distributionFlags;
        LPSTR defaultEnv;
        unsigned long defaultEnvCnt;
        res = WslGetDistributionConfiguration(TargetName,&distributionVersion,&defaultUID,&distributionFlags,&defaultEnv,&defaultEnvCnt);
        if(res!=0)
        {
            fwprintf(stderr,L"ERROR:Get Configuration failed!\nHRESULT:0x%x\n",res);
            wprintf(L"Please any key to continue...");
            getchar();
            return 1;
        }

        wchar_t LxUID[50] = L"";
        if(WslGetLxUID(TargetName,LxUID) == NULL)
        {
            fwprintf(stderr,L"ERROR:GetLxUID failed!\n");
            wprintf(L"Please any key to continue...");
            getchar();
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
                    (void) WslConfigureDistribution(TargetName,0,distributionFlags); //set default uid to 0(root)
                    FILE *fp;
                    unsigned long uid;
                    wchar_t wcmd[300] = L"wsl.exe ";
                    wcscat_s(wcmd,ARRAY_LENGTH(wcmd),LxUID);
                    wcscat_s(wcmd,ARRAY_LENGTH(wcmd),L"id -u ");
                    wcscat_s(wcmd,ARRAY_LENGTH(wcmd),wargv[3]);
                    if((fp=_wpopen(wcmd,L"r")) ==NULL) {
                        fwprintf(stderr,L"ERROR:Command Excute Failed!\n");
                        return 1;
                    }
                    wchar_t buf[256];
                    if (!feof(fp))
                        fgetws(buf, sizeof(buf), fp);

                    (void) pclose(fp);
                    if(swscanf(buf,L"%d",&uid)==1)
                    {
                        res = WslConfigureDistribution(TargetName,uid,distributionFlags);
                        if(res != 0)
                        {
                            (void) WslConfigureDistribution(TargetName,defaultUID,distributionFlags); //revert uid
                            fwprintf(stderr,L"ERROR:Configure Failed!\nHRESULT:0x%x\n",res);
                            return 1;
                        }
                        return 0;
                    }
                    else
                    {
                        (void) WslConfigureDistribution(TargetName,defaultUID,distributionFlags); //revert uid
                        wprintf(L"\n");
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nFailed to detect user.\n");
                    }
                    return 1;
                }
                else if(wcscmp(wargv[2],L"--default-uid") == 0)
                {
                    unsigned long uid;
                    if(swscanf(wargv[3],L"%d",&uid)==1)
                    {
                        res = WslConfigureDistribution(TargetName,uid,distributionFlags);
                        if(res != 0)
                        {
                            fwprintf(stderr,L"ERROR:Configure Failed!\nHRESULT:0x%x\n",res);
                            return 1;
                        }
                        return 0;
                    }
                    else
                    {
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nInput UID\n");
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
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nInput on/off\n");
                        return 1;
                    }
                    res = WslConfigureDistribution(TargetName,defaultUID,distributionFlags);
                    if(res != 0)
                    {
                        fwprintf(stderr,L"ERROR:Configure Failed!\nHRESULT0x%x\n",res);
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
                        fwprintf(stderr,L"ERROR:Invalid Argument.\nInput on/off\n");
                        return 1;
                    }
                    res = WslConfigureDistribution(TargetName,defaultUID,distributionFlags);
                    if(res != 0)
                    {
                        fwprintf(stderr,L"ERROR:Configure Failed!\nHRESULT:0x%x\n",res);
                        return 1;
                    }
                    return 0;
                }
                else
                {
                    fwprintf(stderr,L"ERROR:Invalid Arguments\n");
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
                fwprintf(stderr,L"ERROR:Invalid Arguments\n");
                return 1;
            }
            else
            {
                fwprintf(stderr,L"ERROR:Invalid Arguments.\n");
                wprintf(L"\n");
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
        int i;
        for (i=2;i<wargc;i++)
        {
            wcscat_s(rArgs,ARRAY_LENGTH(rArgs),L" ");
            wcscat_s(rArgs,ARRAY_LENGTH(rArgs),wargv[i]);
        }

        unsigned long exitcode;
        res = WslLaunchInteractive(TargetName,rArgs,1,&exitcode);
        if(res==0)
            return exitcode;
        else
        {
            fwprintf(stderr,L"ERROR:Launch Interactive mode Failed!\nHRESULT:0x%x\n",res);
            wprintf(L"Please any key to continue...");
            getchar();
            return 1;
        }
    }
    else
    {
        wchar_t tgzname[300] = L"rootfs.tar.gz";
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
        wprintf(L"Installing...\n");
        res = WslRegisterDistribution(TargetName,tgzname);
        if(res != 0)
        {
            fwprintf(stderr,L"ERROR:Installation Failed!\nHRESULT:0x%x\n",res);
            wprintf(L"Please any key to continue...");
            getchar();
            return 1;
        }
        wprintf(L"Installation Complete!\n");
        return 0;
    }
    return 0;
}