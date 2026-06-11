# r3d-go

A self-contained [cgo](https://pkg.go.dev/cmd/cgo) binding for
[R3D](https://github.com/Bigfoot71/r3d) — the 3D rendering extension for
[raylib](https://www.raylib.com/) — built from the
[r3dStarter](https://github.com/jensroth-git/r3dStarter) toolchain.

Target: **macOS / arm64** (Apple Silicon). The native static libraries are
vendored, so a clone builds with no extra dependencies beyond a C toolchain.

![showcase](screenshots/showcase.png)

The `cmd/demo` showcase exercises most of R3D in one scene: HDR skybox + IBL, a
procedural sky, PBR models and primitives, a local reflection probe, skeletal
animation, GPU instancing, billboards, an animated sprite, projected decals, a
CPU particle system, transparent surfaces, a stencil outline, three light types
with shadows, optional split-screen multi-view, and the full post-processing
stack — bloom, SSAO, SSGI, SSR, fog, depth-of-field and tonemapping — toggled
live from the keyboard.

## Layout

```
r3d/
├── go.mod
├── Makefile             # build + ad-hoc codesign helper
├── r3d.go               # core: init, camera, lights, meshes, materials, instances, draw
├── material.go          # PBR material + texture maps (albedo/normal/emission/ORM), stencil/depth
├── model.go             # model loading, skeletal animation
├── environment.go       # skybox, cubemaps, IBL, procedural sky, post-processing (bloom/SSAO/SSGI/SSR/fog/DoF/tonemap)
├── camera.go            # R3D camera + multi-view (BeginPro, split-screen)
├── probe.go             # local reflection/illumination probes
├── decal.go             # projected decals
├── extra.go             # extra light controls, instance streaming, texture filtering
├── input.go             # raylib input, keys, colors, HUD helpers
├── math.go              # Vector3 / Quaternion helpers
├── raylib.go            # window, drawing, text, screenshot
├── cmd/demo/main.go     # the feature showcase
├── assets/              # models (glb), HDR panorama, images
├── include/             # vendored public headers (raylib + r3d)
└── lib/                 # vendored static archives (Mach-O arm64)
    ├── libr3d.a         # R3D (shaders + lookup textures embedded at compile time)
    ├── libraylib.a      # raylib 5.5 (image formats PNG/HDR/JPG enabled — see below)
    └── libassimp.a      # Assimp (model loading)
```

## Build & run

```bash
make run        # build, ad-hoc sign, run the showcase
# or
make demo       # build + sign -> ./r3ddemo
./r3ddemo       # interactive
```

Showcase controls: **B** bloom · **O** SSAO · **I** SSGI · **R** SSR · **G** fog ·
**F** depth-of-field · **T** tonemap · **K** procedural sky · **V** split-screen ·
**Space** pause lights · **H** toggle help · **ESC** quit.

> **MSAA + multi-view caveat.** Do not set `FlagMSAA4xHint` if you use
> `BeginPro` sub-viewports (split-screen): a multisampled backbuffer makes R3D's
> sub-viewport blit render black. R3D already anti-aliases at its internal
> resolution, so a multisampled backbuffer adds nothing — the showcase omits it.

## Things that bit us on macOS (and the fixes baked in)

1. **Code signing.** With cgo external linking the toolchain's ad-hoc signature
   can verify as invalid, and dyld kills the process at launch with *"Code
   Signature Invalid"* (SIGKILL, zero output). `make` re-signs ad-hoc after
   `go build` (`codesign --force --sign - r3ddemo`). Run that line yourself if
   you `go build`/`go run` by hand.

2. **Main thread.** raylib/GLFW drive Cocoa, which must run on the process main
   thread. Any program using this package must lock it:
   ```go
   import "runtime"
   func init() { runtime.LockOSThread() }
   ```

3. **Retina / HiDPI.** The GL drawable is 2× the logical window. The showcase
   sets `FlagWindowHighDPI` and initializes R3D at `GetRenderWidth()/Height()`
   so the deferred renderer fills the real framebuffer; 2D HUD font sizes are
   scaled by the DPI factor. (Note: `TakeScreenshot` must be called *after*
   `EndDrawing`, once the batched 2D HUD has been flushed.)

4. **Image formats.** r3dStarter builds raylib with image decoders disabled, so
   HDR skyboxes, PNG/JPG textures and embedded glTF textures fail to load. The
   vendored `libraylib.a` here was rebuilt with PNG/HDR/JPG/QOI enabled (see
   "Rebuilding" below).

## Linking

`r3d.go` carries the cgo directives:

```go
#cgo CFLAGS:  -I${SRCDIR}/include -I${SRCDIR}/include/r3d
#cgo LDFLAGS: -L${SRCDIR}/lib -lr3d -lraylib -lassimp -lc++ -lz -lm
#cgo LDFLAGS: -framework OpenGL -framework Cocoa -framework IOKit \
              -framework CoreFoundation -framework CoreVideo
```

## API coverage

Idiomatic Go types (`Vector3`, `Color`, `Camera3D`, `Quaternion`…) and method
receivers (`Light.SetActive`, `Mesh.Unload`, `Material.SetRoughness`,
`InstanceBuffer.MapPositions`, `AnimationPlayer.Update`…). cgo includes the real
`raylib.h`/`r3d.h`, so extending the binding is mechanical — add a wrapper that
converts the value types and calls the matching `C.R3D_*` function. The full C
API lives in `include/r3d/`.

- **Core**: `Init`, `Close`, `Begin`, `End`
- **Window/raylib**: `InitWindow`, `WindowShouldClose`, `UpdateCamera`, input, `DrawText`, `DrawFPS`, `TakeScreenshot`, …
- **Lights**: dir/spot/omni, direction/position/color/energy/range, cones, soft shadows
- **Meshes**: `GenMeshSphere/Cube/Plane/Cylinder/Torus/Quad`
- **Materials (PBR)**: albedo/normal/emission/ORM maps, metalness/roughness/occlusion/specular, billboard/blend/transparency/cull modes, UV transform, unlit
- **Models + animation**: `LoadModel`, `DrawModel(Ex)`, `LoadAnimationLib`, `LoadAnimationPlayer`, `Play/SetLoop/SetSpeed/Update`, `DrawAnimatedModel(Ex)`
- **Skybox / IBL**: `LoadCubemap`, `GenAmbientMap`, `SetSky`, `SetAmbientMap`, background rotation/energy
- **Procedural sky**: `DefaultProceduralSky`, `SetSun`/`SetSkyColors`, `GenProceduralSky`
- **Reflection probes**: `CreateProbe`, position/range/falloff, `SetUpdateMode`
- **Instancing**: `LoadInstanceBuffer`, `MapPositions/Rotations/Scales/Colors`, `UploadPositions`, `DrawMeshInstanced`
- **Decals**: `NewDecal`, maps, `DrawDecal(Ex)`
- **Transparency / stencil**: material transparency/blend/cull modes, `SetStencil`, `SetDepthMode` (outline / x-ray)
- **Multi-view**: `CameraFromRL`, `BeginPro` with viewports (split-screen, minimaps)
- **Post-processing**: `SetBloom`, `SetSSAO`, `SetSSGI`, `SetSSIL`, `SetSSR`, `SetFog`, `SetDoF`, `SetTonemap`, `SetColorAdjustment`
- **Drawing**: `DrawMesh(Ex)`, environment background/ambient color

## Rebuilding the native archives

The vendored `.a` files were produced by building r3dStarter via CMake:

```bash
cmake -S r3dStarter -B build -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_OSX_ARCHITECTURES=arm64 \
  -DASSIMP_WARNINGS_AS_ERRORS=OFF \
  -DCMAKE_C_FLAGS="-Wno-error -include assert.h" \
  -DCMAKE_CXX_FLAGS="-Wno-error" \
  -DSUPPORT_FILEFORMAT_PNG=ON -DSUPPORT_FILEFORMAT_HDR=ON \
  -DSUPPORT_FILEFORMAT_JPG=ON -DSUPPORT_FILEFORMAT_QOI=ON
cmake --build build -j
# then copy build/lib/libr3d.a, the raylib/assimp archives, and the headers here
```

The non-default flags work around current Apple clang issues (Assimp's `-Werror`
on new warnings; raylib `rmodels.c` using `assert()` without `<assert.h>`) and
re-enable raylib's image decoders (PNG/HDR/JPG/QOI) that r3dStarter turns off.
