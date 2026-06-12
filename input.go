package r3d

/*
#include "raylib.h"
*/
import "C"

// Keyboard keys (raylib KeyboardKey) used by demos/apps.
const (
	KeySpace = int(C.KEY_SPACE)
	KeyEnter = int(C.KEY_ENTER)
	KeyTab   = int(C.KEY_TAB)
	KeyLeft  = int(C.KEY_LEFT)
	KeyRight = int(C.KEY_RIGHT)
	KeyUp    = int(C.KEY_UP)
	KeyDown  = int(C.KEY_DOWN)

	KeyA = int(C.KEY_A)
	KeyB = int(C.KEY_B)
	KeyD = int(C.KEY_D)
	KeyF = int(C.KEY_F)
	KeyG = int(C.KEY_G)
	KeyH = int(C.KEY_H)
	KeyI = int(C.KEY_I)
	KeyK = int(C.KEY_K)
	KeyO = int(C.KEY_O)
	KeyP = int(C.KEY_P)
	KeyR = int(C.KEY_R)
	KeyS = int(C.KEY_S)
	KeyT = int(C.KEY_T)
	KeyV = int(C.KEY_V)

	KeyY = int(C.KEY_Y)
	KeyZ = int(C.KEY_Z)

	KeyLeftControl = int(C.KEY_LEFT_CONTROL)
	KeyLeftSuper   = int(C.KEY_LEFT_SUPER)

	KeyOne   = int(C.KEY_ONE)
	KeyTwo   = int(C.KEY_TWO)
	KeyThree = int(C.KEY_THREE)
	KeyFour  = int(C.KEY_FOUR)
	KeyFive  = int(C.KEY_FIVE)
)

// IsKeyPressed reports whether a key was pressed this frame.
func IsKeyPressed(key int) bool { return bool(C.IsKeyPressed(C.int(key))) }

// IsKeyDown reports whether a key is currently held.
func IsKeyDown(key int) bool { return bool(C.IsKeyDown(C.int(key))) }

// GetTime returns elapsed time in seconds since InitWindow.
func GetTime() float64 { return float64(C.GetTime()) }

// GetMouseWheelMove returns the mouse wheel delta this frame.
func GetMouseWheelMove() float32 { return float32(C.GetMouseWheelMove()) }

// Mouse buttons (raylib MouseButton).
const (
	MouseButtonLeft   = int(C.MOUSE_BUTTON_LEFT)
	MouseButtonRight  = int(C.MOUSE_BUTTON_RIGHT)
	MouseButtonMiddle = int(C.MOUSE_BUTTON_MIDDLE)
	MouseButtonSide   = int(C.MOUSE_BUTTON_SIDE)
	MouseButtonExtra  = int(C.MOUSE_BUTTON_EXTRA)
)

// GetMousePosition returns the cursor position in logical screen pixels.
func GetMousePosition() (x, y float32) {
	p := C.GetMousePosition()
	return float32(p.x), float32(p.y)
}

// IsMouseButtonDown reports whether a mouse button is currently held.
func IsMouseButtonDown(button int) bool { return bool(C.IsMouseButtonDown(C.int(button))) }

// IsMouseButtonPressed reports whether a mouse button was pressed this frame.
func IsMouseButtonPressed(button int) bool { return bool(C.IsMouseButtonPressed(C.int(button))) }

// GetScreenWidth returns the current window width.
func GetScreenWidth() int { return int(C.GetScreenWidth()) }

// GetScreenHeight returns the current window height (logical pixels).
func GetScreenHeight() int { return int(C.GetScreenHeight()) }

// GetRenderWidth returns the framebuffer width in physical pixels (HiDPI-aware).
func GetRenderWidth() int { return int(C.GetRenderWidth()) }

// GetRenderHeight returns the framebuffer height in physical pixels (HiDPI-aware).
func GetRenderHeight() int { return int(C.GetRenderHeight()) }

// ColorFromHSV builds an RGBA color from hue (deg), saturation and value [0,1].
func ColorFromHSV(hue, saturation, value float32) Color {
	c := C.ColorFromHSV(C.float(hue), C.float(saturation), C.float(value))
	return Color{uint8(c.r), uint8(c.g), uint8(c.b), uint8(c.a)}
}

// DrawRectangle draws a filled 2D rectangle (for HUD backgrounds).
func DrawRectangle(x, y, width, height int, col Color) {
	C.DrawRectangle(C.int(x), C.int(y), C.int(width), C.int(height), col.c())
}
