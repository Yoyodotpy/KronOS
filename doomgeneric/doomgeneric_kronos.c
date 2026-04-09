#define _POSIX_C_SOURCE 199309L
#include <time.h>

#include "doomkeys.h"
#include "m_argv.h"
#include "doomgeneric.h"

#include <linux/fb.h>
#include <sys/ioctl.h>
#include <sys/mman.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

int fb_fd = 0;
struct fb_var_screeninfo vinfo;
struct fb_fix_screeninfo finfo;
uint8_t *fbp = NULL;
long int screensize = 0;

void DG_Init() {
	fb_fd = open("/dev/fb0", O_RDWR);
	if (fb_fd == -1) {
		perror("no framebuffer at /dev/fb0");
		exit(1);
	}

	ioctl(fb_fd, FBIOGET_FSCREENINFO, &finfo);
	ioctl(fb_fd, FBIOGET_VSCREENINFO, &vinfo);

	screensize = finfo.smem_len;
	fbp = (uint8_t *)mmap(0, screensize, PROT_READ | PROT_WRITE, MAP_SHARED, fb_fd, 0);
	if ((intptr_t)fbp == -1) {
        perror("no mmap");
        exit(1);
    }
}

void drawpix(int x, int y, int r, int g, int b) {
	if (x < 0 || x >= vinfo.xres || y < 0 || y >= vinfo.yres) return;

	long int loc = (x * (vinfo.bits_per_pixel / 8)) + (y * finfo.line_length);

    if (vinfo.bits_per_pixel == 32) {
        fbp[loc + (vinfo.red.offset / 8)] = r;
        fbp[loc + (vinfo.green.offset / 8)] = g;
        fbp[loc + (vinfo.blue.offset / 8)] = b;
        fbp[loc + (vinfo.transp.offset / 8)] = 0xFF;
    }
}

#define DOOM_W 320
#define DOOM_H 200

void DG_DrawFrame() {
    for (int y = 0; y < vinfo.yres; y++) {
        for (int x = 0; x < vinfo.xres; x++) {
            int doom_x = (x * DOOM_W) / vinfo.xres;
            int doom_y = (y * DOOM_H) / vinfo.yres;

            uint32_t pixel = DG_ScreenBuffer[doom_y * DOOM_W + doom_x];

            uint8_t r = (pixel >> 16) & 0xFF;
            uint8_t g = (pixel >> 8) & 0xFF;
            uint8_t b = pixel & 0xFF;

            long int loc = (x * (vinfo.bits_per_pixel / 8)) + (y * finfo.line_length);

            if (loc >= 0 && (loc + 3) < screensize) {
                fbp[loc + (vinfo.red.offset / 8)] = r;
                fbp[loc + (vinfo.green.offset / 8)] = g;
                fbp[loc + (vinfo.blue.offset / 8)] = b;

                if (vinfo.bits_per_pixel == 32) {
                    fbp[loc + (vinfo.transp.offset / 8)] = 0xFF;
                }
            }
        }
    }
}

void DG_SleepMs(uint32_t ms) {
	return;
}

void DG_SetWindowTitle(const char * title) {
}

uint32_t DG_GetTicksMs() {
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (ts.tv_sec * 1000) + (ts.tv_nsec / 1000000);
}

int DG_GetKey(int* pressed, unsigned char* doomKey) {
	return 0;
}

int main(int argc, char **argv)
{
    doomgeneric_Create(argc, argv);

    while (1)
    {
        doomgeneric_Tick();
    }

    return 0;
}
