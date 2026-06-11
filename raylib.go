package r3d

/*
#include <stdlib.h>
#include "raylib.h"

// TraceLog is variadic and cgo cannot call variadic C functions directly,
// so route through a fixed-arity shim.
static void r3d_tracelog(int level, const char *msg) { TraceLog(level, "%s", msg); }
*/
import "C"

import "unsafe"

// Config flags (raylib ConfigFlags) — pass to SetConfigFlags before InitWindow.
const (
	FlagWindowResizable   = uint32(C.FLAG_WINDOW_RESIZABLE)
	FlagWindowUndecorated = uint32(C.FLAG_WINDOW_UNDECORATED)
	FlagWindowHighDPI     = uint32(C.FLAG_WINDOW_HIGHDPI)
	FlagVSyncHint         = uint32(C.FLAG_VSYNC_HINT)
	FlagMSAA4xHint        = uint32(C.FLAG_MSAA_4X_HINT)
	FlagFullscreenMode    = uint32(C.FLAG_FULLSCREEN_MODE)
)

// Trace-log levels (raylib TraceLogLevel).
const (
	LogAll     = int(C.LOG_ALL)
	LogTrace   = int(C.LOG_TRACE)
	LogDebug   = int(C.LOG_DEBUG)
	LogInfo    = int(C.LOG_INFO)
	LogWarning = int(C.LOG_WARNING)
	LogError   = int(C.LOG_ERROR)
	LogFatal   = int(C.LOG_FATAL)
	LogNone    = int(C.LOG_NONE)
)

// Common colors.
var (
	White    = Color{255, 255, 255, 255}
	Black    = Color{0, 0, 0, 255}
	Red      = Color{230, 41, 55, 255}
	Green    = Color{0, 228, 48, 255}
	Blue     = Color{0, 121, 241, 255}
	DarkGray = Color{80, 80, 80, 255}
	RayWhite = Color{245, 245, 245, 255}
)

// SetConfigFlags sets window/context configuration flags (call before InitWindow).
func SetConfigFlags(flags uint32) { C.SetConfigFlags(C.uint(flags)) }

// InitWindow creates the window and OpenGL context.
func InitWindow(width, height int, title string) {
	ct := C.CString(title)
	defer C.free(unsafe.Pointer(ct))
	C.InitWindow(C.int(width), C.int(height), ct)
}

// CloseWindow closes the window and unloads the OpenGL context.
func CloseWindow() { C.CloseWindow() }

// WindowShouldClose reports whether a close was requested (ESC / close button).
func WindowShouldClose() bool { return bool(C.WindowShouldClose()) }

// IsWindowReady reports whether the window was successfully initialized.
func IsWindowReady() bool { return bool(C.IsWindowReady()) }

// SetTargetFPS sets the target frame rate.
func SetTargetFPS(fps int) { C.SetTargetFPS(C.int(fps)) }

// GetFPS returns the current frame rate.
func GetFPS() int { return int(C.GetFPS()) }

// BeginDrawing starts a 2D/canvas drawing pass.
func BeginDrawing() { C.BeginDrawing() }

// EndDrawing ends the drawing pass and swaps buffers.
func EndDrawing() { C.EndDrawing() }

// ClearBackground clears the framebuffer with the given color.
func ClearBackground(col Color) { C.ClearBackground(col.c()) }

// DrawText draws text using the default font.
func DrawText(text string, x, y, fontSize int, col Color) {
	ct := C.CString(text)
	defer C.free(unsafe.Pointer(ct))
	C.DrawText(ct, C.int(x), C.int(y), C.int(fontSize), col.c())
}

// DrawFPS draws the current FPS in the top-left corner.
func DrawFPS(x, y int) { C.DrawFPS(C.int(x), C.int(y)) }

// MeasureText measures default-font text width in pixels.
func MeasureText(text string, fontSize int) int {
	ct := C.CString(text)
	defer C.free(unsafe.Pointer(ct))
	return int(C.MeasureText(ct, C.int(fontSize)))
}

// TraceLog writes a message to the raylib log at the given level.
func TraceLog(level int, text string) {
	ct := C.CString(text)
	defer C.free(unsafe.Pointer(ct))
	C.r3d_tracelog(C.int(level), ct)
}

// UpdateCamera updates the camera in-place according to the given camera mode.
func UpdateCamera(cam *Camera3D, mode int) {
	cc := cam.c()
	C.UpdateCamera(&cc, C.int(mode))
	*cam = goCamera(cc)
}

// GetFrameTime returns the time in seconds for the last frame drawn.
func GetFrameTime() float32 { return float32(C.GetFrameTime()) }

// TakeScreenshot saves the current framebuffer to a PNG file (relative to CWD).
func TakeScreenshot(fileName string) {
	cf := C.CString(fileName)
	defer C.free(unsafe.Pointer(cf))
	C.TakeScreenshot(cf)
}
