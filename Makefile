# Build helpers for the r3d-go binding.
#
# Cross-platform: macOS/arm64 and Windows/amd64. The native archives live under
# lib/<goos>_<goarch>/ and the per-OS cgo link flags under cgo_<goos>.go.
#
# macOS note: on Apple Silicon every executable must carry a valid code
# signature. The signature the Go/ld toolchain embeds during cgo external
# linking can come out invalid ("linker-signed" but failing `codesign -v`),
# which makes dyld kill the process at launch with "Code Signature Invalid".
# Re-signing ad-hoc fixes it, so the macOS build re-signs after building.

PKG := ./cmd/demo

# Detect the host OS so the same Makefile works on both platforms.
ifeq ($(OS),Windows_NT)
  BIN  := r3ddemo.exe
  SIGN := @echo "(windows: no code signing needed)"
  RUN  := ./$(BIN)
else
  BIN  := r3ddemo
  SIGN := codesign --force --sign - $(BIN) && codesign -v $(BIN) && echo "signature OK"
  RUN  := ./$(BIN)
endif

.PHONY: all demo lib run smoke shot clean

all: demo

# Build the showcase (and, on macOS, produce a valid ad-hoc signature).
demo:
	go build -o $(BIN) $(PKG)
	$(SIGN)

# Type-check / build the library package only (no signing needed for a lib).
lib:
	go build ./...

# Build then run the showcase interactively.
run: demo
	$(RUN)

# Headless smoke test: render 150 frames then auto-exit.
smoke: demo
	R3D_DEMO_FRAMES=150 $(RUN)

# Render headlessly and save a screenshot to screenshots/showcase.png.
shot: demo
	R3D_DEMO_FRAMES=180 R3D_DEMO_SHOT=screenshots/showcase.png $(RUN)

clean:
	rm -f r3ddemo r3ddemo.exe showcase.png
