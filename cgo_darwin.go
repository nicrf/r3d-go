//go:build darwin

package r3d

// macOS / arm64 link flags. The native archives are vendored as Mach-O arm64
// under lib/darwin_arm64/. raylib drives Cocoa/OpenGL through Apple frameworks.

/*
#cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -lr3d -lraylib -lassimp -lc++ -lz -lm
#cgo LDFLAGS: -framework OpenGL -framework Cocoa -framework IOKit -framework CoreFoundation -framework CoreVideo
*/
import "C"
