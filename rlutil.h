#pragma once
/**
 * https://github.com/tapio/rlutil
 *
 * File: rlutil.h
 *
 * About: Description
 * This file provides some useful utilities for console mode
 * roguelike game development with C and C++. It is aimed to
 * be cross-platform (at least Windows and Linux).
 *
 * About: Copyright
 * (C) 2010 Tapio Vierros
 *
 * About: Licensing
 * See <License>
 */

#ifndef RLUTIL_INLINE
	#define RLUTIL_INLINE __inline__
#endif

#include <stdio.h> // for getch()
#include <termios.h> // for getch() and kbhit()
#include <unistd.h> // for getch(), kbhit() and (u)sleep()
#include <sys/ioctl.h> // for getkey()
#include <sys/types.h> // for kbhit()
#include <sys/time.h> // for kbhit()

/// Function: getch
/// Get character without waiting for Return to be pressed.
/// Windows has this in conio.h
RLUTIL_INLINE int getch(void) {
        // Here be magic.
        struct termios oldt, newt;
        int ch;
        tcgetattr(STDIN_FILENO, &oldt);
        newt = oldt;
        newt.c_lflag &= ~(ICANON | ECHO);
        tcsetattr(STDIN_FILENO, TCSANOW, &newt);
        ch = getchar();
        tcsetattr(STDIN_FILENO, TCSANOW, &oldt);
        return ch;
}

/// Function: kbhit
/// Determines if keyboard has been hit.
/// Windows has this in conio.h
RLUTIL_INLINE int kbhit(void) {
        // Here be dragons.
        static struct termios oldt, newt;
        int cnt = 0;
        tcgetattr(STDIN_FILENO, &oldt);
        newt = oldt;
        newt.c_lflag    &= ~(ICANON | ECHO);
        newt.c_iflag     = 0; // input mode
        newt.c_oflag     = 0; // output mode
        newt.c_cc[VMIN]  = 1; // minimum time to wait
        newt.c_cc[VTIME] = 1; // minimum characters to wait for
        tcsetattr(STDIN_FILENO, TCSANOW, &newt);
        ioctl(0, FIONREAD, &cnt); // Read count
        struct timeval tv;
        tv.tv_sec  = 0;
        tv.tv_usec = 100;
        select(STDIN_FILENO+1, NULL, NULL, NULL, &tv); // A small time delay
        tcsetattr(STDIN_FILENO, TCSANOW, &oldt);
        return cnt; // Return number of characters
}
