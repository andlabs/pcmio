// 30 august 2014

#include "winapi_windows.h"

DWORD xCreateEvent(HANDLE *h)
{
	*h = CreateEvent(NULL, TRUE, TRUE, NULL);			// start off signaled just in case
	if (*h == NULL)
		return GetLastError();
	return 0;
}

DWORD xResetEvent(HANDLE h)
{
	if (ResetEvent(h) == 0)
		return GetLastError();
	return 0;
}

DWORD xWaitForSingleObject(HANDLE h)
{
	if (WaitForSingleObject(h, INFINITE) == WAIT_FAILED)
		return GetLastError();
	return 0;
}
