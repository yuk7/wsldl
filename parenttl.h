/*
 * Copyright (c) 2020 yuk7
 * Author: yuk7 <yukx00@gmail.com>
 *
 * Released under the MIT license
 * http://opensource.org/licenses/mit-license.php
 */

#ifndef PARENTTL_H_
#define PARENTTL_H_

#include <stdbool.h>
#include <stdio.h>
#include <windows.h>
#include <tlhelp32.h>

#ifdef __cplusplus
extern "C" {
#endif

#define PROC_LIST_SIZE 20000

bool isParentCmdLine()
{
    HANDLE hSnapshot;
    DWORD procs[PROC_LIST_SIZE];
    int procsCnt = 0;
    DWORD procID = GetCurrentProcessId();
    
    if((hSnapshot = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS,0)) != INVALID_HANDLE_VALUE)
    {
        PROCESSENTRY32 pe32;
        pe32.dwSize = sizeof(PROCESSENTRY32);
        
        if(Process32First(hSnapshot,&pe32))
        {
            do
            {
                if ((strcmp(pe32.szExeFile, "cmd.exe") == 0) ||
                (strcmp(pe32.szExeFile, "powershell.exe") == 0) ||
                (strcmp(pe32.szExeFile, "wsl.exe") == 0) )
                {
                    if(procsCnt < PROC_LIST_SIZE)
                    {
                        procs[procsCnt] = pe32.th32ProcessID;
                        procsCnt++;
                    }
                }
                if (pe32.th32ProcessID == procID)
                {
                    for(int i = 0; i <= procsCnt; i++)
                    {
                        if(procs[i] == pe32.th32ParentProcessID)
                        {
                            return true;
                        }
                    }
                    return false;
                }
            }
            while(Process32Next(hSnapshot, &pe32));
        }
        CloseHandle(hSnapshot);
    }
    return false;
}

#ifdef __cplusplus
}
#endif

#endif
