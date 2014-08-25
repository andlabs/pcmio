// 25 august 2014
// main_windows.go 24 august 2014
#include <windows.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

void checkerr(MMRESULT err, char *when)
{
	if (err != MMSYSERR_NOERROR) {
		fprintf(stderr, "error %s: %d\n", when, err);		// TODO
		exit(EXIT_FAILURE);
	}
}

int main(void)
{
	HWAVEOUT hwo;
	WAVEFORMATEX format;
	WAVEHDR header;
	uint8_t buffer[256 * 440];
	int i, j, k;
	HANDLE event;

	event = CreateEvent(NULL, TRUE, TRUE, NULL);			// start off signaled just in case
	if (event == NULL) {
		fprintf(stderr, "error creating event: %d", GetLastError());
		return EXIT_FAILURE;
	}

	ZeroMemory(&format, sizeof (WAVEFORMATEX));
	format.wFormatTag = WAVE_FORMAT_PCM;
	format.nChannels = 1;
	format.nSamplesPerSec = 44100;
	format.wBitsPerSample = 8;
	format.nBlockAlign = format.nChannels * format.wBitsPerSample;
	format.nAvgBytesPerSec = format.nSamplesPerSec * (DWORD) format.nBlockAlign;
	format.cbSize = 0;
	checkerr(waveOutOpen(&hwo, WAVE_MAPPER, &format,
		(DWORD_PTR) event, 0, CALLBACK_EVENT | WAVE_FORMAT_DIRECT), "opening PCM stream");

	k = 0;
	for (i = 0; i < 440; i++) {
		int min, max, step;

		min = 0;
		max = 256;
		step = 1;
		if (i % 2 == 1) {
			min = 255;
			max = -1;
			step = -1;
		}
		for (j = min; j != max; j += step) {
			buffer[k] = (uint8_t) (j & 0xFF);
			k++;
		}
	}

	for (i = 0; i < 20; i++) {
		header.lpData = (LPSTR) buffer;
		header.dwBufferLength = 256 * 440;
		if (header.dwBufferLength <= 0)
			break;
		header.dwFlags = 0;
		header.dwLoops = 0;
		if (ResetEvent(event) == 0) {
			fprintf(stderr, "error resetting event: %d", GetLastError());
			return EXIT_FAILURE;
		}
		checkerr(waveOutPrepareHeader(hwo, &header, sizeof header),
			"pareparing PCM data for playack");
		checkerr(waveOutWrite(hwo, &header, sizeof header),
			"playing PCM data");
		if (WaitForSingleObject(event, INFINITE) == WAIT_FAILED) {
			fprintf(stderr, "error waiting for event", GetLastError());
			return EXIT_FAILURE;
		}
	}

	checkerr(waveOutClose(hwo), "closing PCM stream");
	return 0;
}
