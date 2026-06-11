# Build helpers for the r3d-go binding.
#
# On Apple Silicon every executable must carry a valid code signature. The
# signature the Go/ld toolchain embeds during cgo external linking can come out
# invalid ("linker-signed" but failing `codesign -v`), which makes dyld kill the
# process at launch with "Code Signature Invalid". Re-signing ad-hoc fixes it,
# so every binary target re-signs after building.

BIN := r3ddemo
PKG := ./cmd/demo

.PHONY: all demo lib run smoke shot clean

all: demo

# Build the showcase and produce a valid ad-hoc signature.
demo:
	go build -o $(BIN) $(PKG)
	codesign --force --sign - $(BIN)
	@codesign -v $(BIN) && echo "signature OK"

# Type-check / build the library package only (no signing needed for a lib).
lib:
	go build ./...

# Build then run the showcase interactively.
run: demo
	./$(BIN)

# Headless smoke test: render 150 frames then auto-exit.
smoke: demo
	R3D_DEMO_FRAMES=150 ./$(BIN)

# Render headlessly and save a screenshot to screenshots/showcase.png.
shot: demo
	R3D_DEMO_FRAMES=180 R3D_DEMO_SHOT=screenshots/showcase.png ./$(BIN)

clean:
	rm -f $(BIN) showcase.png
