//go:build windows

package r3d

// Windows / amd64 link flags. The native archives must be rebuilt for MinGW
// (gcc) and vendored under lib/windows_amd64/ (libr3d.a, libraylib.a,
// libassimp.a). raylib drives the window/GL context through the Win32 system
// libraries below instead of the macOS frameworks.
//
// Notes:
//   - -lstdc++ replaces macOS's -lc++ (Assimp is C++; MinGW's runtime is libstdc++).
//   - opengl32/gdi32/winmm/user32/shell32 are what raylib 5.5 links on Windows;
//     winmm is required for raylib's timer, gdi32 for the GL pixel format.
//   - zlib is vendored as libzlibstatic.a (Assimp's own static zlib build);
//     it must come AFTER -lassimp, which references its inflate/deflate symbols.
//   - System libraries come AFTER the static archives so the static linker can
//     resolve their symbols (GNU ld is order-sensitive).

/*
#cgo LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -lr3d -lraylib -lassimp -lzlibstatic
#cgo LDFLAGS: -lopengl32 -lgdi32 -lwinmm -luser32 -lshell32 -lstdc++ -lm
*/
import "C"
