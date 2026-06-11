// Command demo is a feature showcase for the r3d-go binding. It exercises most
// of what R3D offers in a single scene — HDR skybox + image-based lighting, a
// procedural sky, PBR models and primitives, a local reflection probe, skeletal
// animation, GPU instancing, billboards, an animated sprite, projected decals, a
// CPU particle system, transparent surfaces, a stencil outline, three light
// types with shadows, optional split-screen multi-view, and the full
// post-processing stack (bloom, SSAO, SSGI, SSR, fog, depth-of-field,
// tonemapping) — all toggled live from the keyboard.
package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"

	r3d "github.com/nicolasfournier/r3d-go"
)

// raylib/GLFW drive Cocoa, which must run on the process main thread.
func init() { runtime.LockOSThread() }

const assets = "assets/"

type particle struct {
	pos, vel r3d.Vector3
	life     float32
}

func main() {
	const width, height = 1600, 900

	// NB: no FLAG_MSAA_4X_HINT — an MSAA backbuffer breaks R3D_BeginPro
	// sub-viewport rendering (split-screen). R3D does its own AA at the
	// internal resolution, so a multisampled backbuffer adds nothing anyway.
	r3d.SetConfigFlags(r3d.FlagWindowResizable | r3d.FlagVSyncHint | r3d.FlagWindowHighDPI)
	r3d.InitWindow(width, height, "R3D — Go showcase")
	r3d.SetTargetFPS(60)
	if !r3d.IsWindowReady() {
		r3d.TraceLog(r3d.LogError, "window init failed")
		return
	}
	r3d.Init(r3d.GetRenderWidth(), r3d.GetRenderHeight())
	r3d.SetTextureFilter(r3d.FilterAnisotropic4x)

	// --- Camera --------------------------------------------------------------
	cam := r3d.Camera3D{
		Position:   r3d.Vector3{X: 8, Y: 5, Z: 8},
		Target:     r3d.Vector3{X: 0, Y: 1, Z: 0},
		Up:         r3d.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       60,
		Projection: int32(r3d.CameraPerspective),
	}

	// --- Skybox: HDR (IBL) + a procedural alternative -----------------------
	hdrSky := r3d.LoadCubemap(assets+"panorama/outdoor.hdr", r3d.CubemapPanorama)
	ibl := r3d.GenAmbientMap(hdrSky, r3d.AmbientIllumination|r3d.AmbientReflection)

	procParams := r3d.DefaultProceduralSky()
	procParams.SetSun(r3d.Vector3{X: -1, Y: -0.6, Z: -0.4}, r3d.Color{R: 255, G: 240, B: 210, A: 255}, 3)
	procParams.SetSkyColors(r3d.Color{R: 60, G: 120, B: 220, A: 255}, r3d.Color{R: 180, G: 210, B: 245, A: 255}, 2)
	procSky := r3d.GenProceduralSky(1024, procParams)

	r3d.SetSky(hdrSky)
	r3d.SetAmbientMap(ibl)
	r3d.SetBackgroundEnergy(1.0)
	r3d.SetAmbientEnergy(1.0)
	r3d.SetTonemap(r3d.TonemapACES, 1.0)

	// --- Local reflection probe (for the chrome sphere) ---------------------
	probe := r3d.CreateProbe(r3d.ProbeReflection | r3d.ProbeIllumination)
	probe.SetPosition(r3d.Vector3{X: 3, Y: 1.2, Z: 3})
	probe.SetRange(8)
	probe.SetUpdateMode(r3d.ProbeUpdateOnce)
	probe.SetActive(true)

	// --- Ground plane (slightly glossy so SSR has something to reflect) -----
	ground := r3d.GenMeshPlane(40, 40, 1, 1)
	groundMat := r3d.GetDefaultMaterial()
	groundMat.SetAlbedoColor(r3d.Color{R: 60, G: 64, B: 70, A: 255})
	groundMat.SetRoughness(0.5)
	groundMat.SetMetalness(0.1)

	// --- PBR model: Damaged Helmet ------------------------------------------
	helmet := r3d.LoadModel(assets + "models/DamagedHelmet.glb")

	// --- Animated character: Cesium Man -------------------------------------
	character := r3d.LoadModel(assets + "models/CesiumMan.glb")
	anims := r3d.LoadAnimationLib(assets + "models/CesiumMan.glb")
	player := r3d.LoadAnimationPlayer(character, anims)
	if anims.Count() > 0 {
		player.SetLoop(0, true)
		player.Play(0)
	}

	// --- Meshes shared by several effects -----------------------------------
	sphere := r3d.GenMeshSphere(0.5, 48, 48)
	cube := r3d.GenMeshCube(0.6, 0.6, 0.6)

	// --- Chrome sphere (mirror, lit by the probe) ---------------------------
	chromeMat := r3d.GetDefaultMaterial()
	chromeMat.SetAlbedoColor(r3d.White)
	chromeMat.SetMetalness(1.0)
	chromeMat.SetRoughness(0.05)

	// --- Transparent glass cubes --------------------------------------------
	glassColors := []r3d.Color{{R: 120, G: 200, B: 255, A: 110}, {R: 255, G: 140, B: 160, A: 110}, {R: 160, G: 255, B: 170, A: 110}}

	// --- Stencil outline materials (solid writes mask, ring draws outline) ---
	outlineSolid := r3d.GetDefaultMaterial()
	outlineSolid.SetAlbedoColor(r3d.Color{R: 40, G: 40, B: 50, A: 255})
	outlineSolid.SetStencil(r3d.CompareAlways, 0x02, 0xFF, r3d.StencilKeep, r3d.StencilKeep, r3d.StencilReplace)
	outlineRing := r3d.GetDefaultMaterial()
	outlineRing.SetStencil(r3d.CompareNotEqual, 0x02, 0xFF, r3d.StencilKeep, r3d.StencilKeep, r3d.StencilKeep)
	outlineRing.SetCullMode(r3d.CullFront)
	outlineRing.SetEmission(r3d.Color{R: 255, G: 180, B: 0, A: 255}, 3)
	outlineRing.SetUnlit(true)

	// --- Instanced ring of cubes --------------------------------------------
	const ringCount = 24
	ringMesh := r3d.GenMeshCube(0.5, 0.5, 0.5)
	ringMat := r3d.GetDefaultMaterial()
	ringMat.SetRoughness(0.3)
	ringMat.SetMetalness(0.9)
	ring := r3d.LoadInstanceBuffer(ringCount,
		r3d.InstancePosition|r3d.InstanceRotation|r3d.InstanceScale|r3d.InstanceColor)
	{
		pos := ring.MapPositions(true)
		rot := ring.MapRotations(true)
		scl := ring.MapScales(true)
		col := ring.MapColors(true)
		for i := 0; i < ringCount; i++ {
			a := float32(i) / ringCount * 2 * math.Pi
			pos[i] = r3d.Vector3{X: 11 * cosf(a), Y: 0.5, Z: 11 * sinf(a)}
			rot[i] = r3d.QuaternionFromAxisAngle(r3d.Vector3{X: 0, Y: 1, Z: 0}, a)
			scl[i] = r3d.Vector3{X: 1, Y: 1, Z: 1}
			col[i] = r3d.ColorFromHSV(float32(i)/ringCount*360, 0.7, 1.0)
		}
		ring.Unmap(r3d.InstancePosition | r3d.InstanceRotation | r3d.InstanceScale | r3d.InstanceColor)
	}

	// --- Billboard trees -----------------------------------------------------
	quad := r3d.GenMeshQuad(2, 2, 1, 1, r3d.Vector3{X: 0, Y: 0, Z: 1})
	treeMat := r3d.GetDefaultMaterial()
	treeMat.SetAlbedoMap(r3d.LoadAlbedoMap(assets+"images/tree.png", r3d.White))
	treeMat.SetBillboardMode(r3d.BillboardYAxis)
	treeMat.SetTransparencyMode(r3d.TransparencyAlpha)
	treePositions := []r3d.Vector3{{X: -6, Y: 1, Z: -3}, {X: -7, Y: 1, Z: 2}, {X: 5, Y: 1, Z: 6}}

	// --- Animated sprite -----------------------------------------------------
	spriteMat := r3d.GetDefaultMaterial()
	spriteMat.SetAlbedoMap(r3d.LoadAlbedoMap(assets+"images/spritesheet.png", r3d.White))
	spriteMat.SetBillboardMode(r3d.BillboardFront)
	spriteMat.SetTransparencyMode(r3d.TransparencyAlpha)
	spriteMat.SetUnlit(true)
	const spriteFrames = 4

	// --- Decal projected on the ground --------------------------------------
	decal := r3d.NewDecal()
	decal.SetAlbedoMap(r3d.LoadAlbedoMap(assets+"images/decal.png", r3d.White))
	decal.SetNormalMap(r3d.LoadNormalMap(assets+"images/decal_normal.png", 1.0))
	decal.SetNormalThreshold(45)
	decal.SetFadeWidth(20)

	// --- Particle system (additive emissive billboards) ---------------------
	const maxParticles = 2048
	particles := make([]particle, 0, maxParticles)
	partPositions := make([]r3d.Vector3, maxParticles)
	partMesh := r3d.GenMeshSphere(0.06, 6, 6)
	partMat := r3d.GetDefaultMaterial()
	partMat.SetBillboardMode(r3d.BillboardFront)
	partMat.SetBlendMode(r3d.BlendAdditive)
	partMat.SetEmission(r3d.Color{R: 255, G: 130, B: 40, A: 255}, 1.5)
	partMat.SetUnlit(true)
	partBuf := r3d.LoadInstanceBuffer(maxParticles, r3d.InstancePosition)
	emitter := r3d.Vector3{X: 6, Y: 0.2, Z: -6}

	// --- Lights --------------------------------------------------------------
	sun := r3d.CreateLight(r3d.LightDir)
	sun.SetDirection(r3d.Vector3{X: -1, Y: -1.4, Z: -0.8})
	sun.SetEnergy(1.4)
	sun.SetShadowSoftness(1.5)
	sun.EnableShadow()
	sun.SetActive(true)

	const omniCount = 3
	omni := make([]r3d.Light, omniCount)
	for i := range omni {
		l := r3d.CreateLight(r3d.LightOmni)
		l.SetColor(r3d.ColorFromHSV(float32(i)/omniCount*360, 1, 1))
		l.SetEnergy(2.0)
		l.SetRange(14)
		l.SetActive(true)
		omni[i] = l
	}

	spot := r3d.CreateLight(r3d.LightSpot)
	spot.LookAt(r3d.Vector3{X: 0, Y: 9, Z: 4}, r3d.Vector3{X: 0, Y: 1, Z: 0})
	spot.SetColor(r3d.Color{R: 255, G: 240, B: 220, A: 255})
	spot.SetEnergy(3.0)
	spot.SetRange(40)
	spot.SetInnerCutOff(14)
	spot.SetOuterCutOff(24)
	spot.SetShadowSoftness(1.0)
	spot.EnableShadow()
	spot.SetActive(true)

	// --- Everything drawn inside an R3D pass (reused for split-screen) -------
	drawScene := func(t float32) {
		r3d.DrawMesh(ground, groundMat, r3d.Vector3{X: 0, Y: 0, Z: 0}, 1)

		spin := r3d.QuaternionFromEuler(math.Pi/2, t*0.4, 0)
		r3d.DrawModelEx(helmet, r3d.Vector3{X: 0, Y: 1.6, Z: 0}, spin, r3d.Vector3{X: 1.2, Y: 1.2, Z: 1.2})

		face := r3d.QuaternionFromAxisAngle(r3d.Vector3{X: 0, Y: 1, Z: 0}, math.Pi)
		r3d.DrawAnimatedModelEx(character, player, r3d.Vector3{X: -3, Y: 0, Z: 0}, face, r3d.Vector3{X: 1, Y: 1, Z: 1})

		// PBR roughness sweep (fully metallic)
		for i := 0; i < 6; i++ {
			m := r3d.GetDefaultMaterial()
			m.SetAlbedoColor(r3d.Color{R: 220, G: 180, B: 90, A: 255})
			m.SetMetalness(1.0)
			m.SetRoughness(float32(i) / 5)
			r3d.DrawMeshEx(sphere, m,
				r3d.Vector3{X: float32(i)*1.3 - 3.2, Y: 0.5, Z: 4},
				r3d.QuaternionIdentity(), r3d.Vector3{X: 1, Y: 1, Z: 1})
		}

		// chrome sphere reflecting the local probe
		r3d.DrawMeshEx(sphere, chromeMat, r3d.Vector3{X: 3, Y: 1.2, Z: 3},
			r3d.QuaternionIdentity(), r3d.Vector3{X: 1.4, Y: 1.4, Z: 1.4})

		// stencil-outlined sphere (solid pass writes the mask, ring pass outlines)
		outlinePos := r3d.Vector3{X: -3.5, Y: 0.8, Z: 4.5}
		r3d.DrawMeshEx(sphere, outlineSolid, outlinePos, r3d.QuaternionIdentity(), r3d.Vector3{X: 1, Y: 1, Z: 1})
		r3d.DrawMeshEx(sphere, outlineRing, outlinePos, r3d.QuaternionIdentity(), r3d.Vector3{X: 1.1, Y: 1.1, Z: 1.1})

		// instanced ring + particles
		r3d.DrawMeshInstanced(ringMesh, ringMat, ring, ringCount)
		r3d.DrawMeshInstanced(partMesh, partMat, partBuf, len(particles))

		// billboard trees
		for _, p := range treePositions {
			r3d.DrawMeshEx(quad, treeMat, p, r3d.QuaternionIdentity(), r3d.Vector3{X: 2, Y: 2, Z: 2})
		}

		// animated sprite
		r3d.DrawMeshEx(quad, spriteMat, r3d.Vector3{X: 3, Y: 1.5, Z: -2},
			r3d.QuaternionIdentity(), r3d.Vector3{X: 1, Y: 1, Z: 1})

		// decals on the ground
		r3d.DrawDecal(decal, r3d.Vector3{X: 1.5, Y: 0.05, Z: 2}, 2)
		r3d.DrawDecalEx(decal, r3d.Vector3{X: -2, Y: 0.05, Z: -2},
			r3d.QuaternionFromAxisAngle(r3d.Vector3{X: 0, Y: 1, Z: 0}, 0.6), r3d.Vector3{X: 2.5, Y: 2.5, Z: 2.5})

		// transparent glass cubes (drawn last; alpha-blended)
		for i, gc := range glassColors {
			m := r3d.GetDefaultMaterial()
			m.SetAlbedoColor(gc)
			m.SetTransparencyMode(r3d.TransparencyAlpha)
			m.SetRoughness(0.05)
			r3d.DrawMeshEx(cube, m,
				r3d.Vector3{X: float32(i)*1.4 - 1.4, Y: 0.4, Z: -3.5},
				r3d.QuaternionFromAxisAngle(r3d.Vector3{X: 0, Y: 1, Z: 0}, t*0.6+float32(i)),
				r3d.Vector3{X: 1, Y: 1, Z: 1})
		}
	}

	// --- Effect state (toggled live) ----------------------------------------
	bloomMode := r3d.BloomMix
	r3d.SetBloom(bloomMode, 0.1, 1.0)
	ssaoOn := true
	r3d.SetSSAO(ssaoOn, 1.0, 0.5, 1.5)
	fogMode := r3d.FogDisabled
	dofOn := false
	ssrOn := false
	ssgiOn := false
	procOn := false
	split := false
	tonemap := r3d.TonemapACES
	paused := false
	showHelp := true

	// Optional initial states for headless testing/screenshots.
	if os.Getenv("R3D_DEMO_SPLIT") != "" {
		split = true
	}
	if os.Getenv("R3D_DEMO_PROC") != "" {
		procOn = true
		r3d.SetSky(procSky)
	}
	if os.Getenv("R3D_DEMO_SSR") != "" {
		ssrOn = true
		r3d.SetSSR(true, 6, 0.2)
	}
	if os.Getenv("R3D_DEMO_SSGI") != "" {
		ssgiOn = true
		r3d.SetSSGI(true, 1.5)
	}

	maxFrames, _ := strconv.Atoi(os.Getenv("R3D_DEMO_FRAMES"))
	frame := 0
	var orbit float32

	for !r3d.WindowShouldClose() {
		if maxFrames > 0 && frame >= maxFrames {
			break
		}
		frame++
		dt := r3d.GetFrameTime()
		t := float32(r3d.GetTime())

		// --- input -----------------------------------------------------------
		if r3d.IsKeyPressed(r3d.KeyB) {
			bloomMode = r3d.BloomMode((int(bloomMode) + 1) % 4)
			r3d.SetBloom(bloomMode, 0.1, 1.0)
		}
		if r3d.IsKeyPressed(r3d.KeyO) {
			ssaoOn = !ssaoOn
			r3d.SetSSAO(ssaoOn, 1.0, 0.5, 1.5)
		}
		if r3d.IsKeyPressed(r3d.KeyG) {
			fogMode = r3d.FogMode((int(fogMode) + 1) % 4)
			r3d.SetFog(fogMode, r3d.Color{R: 160, G: 170, B: 190, A: 255}, 5, 45, 0.04)
		}
		if r3d.IsKeyPressed(r3d.KeyF) {
			dofOn = !dofOn
			r3d.SetDoF(dofOn, 8, 2, 20)
		}
		if r3d.IsKeyPressed(r3d.KeyT) {
			tonemap = r3d.TonemapMode((int(tonemap) + 1) % r3d.TonemapModeCount)
			r3d.SetTonemap(tonemap, 1.0)
		}
		if r3d.IsKeyPressed(r3d.KeyR) {
			ssrOn = !ssrOn
			r3d.SetSSR(ssrOn, 6, 0.2)
		}
		if r3d.IsKeyPressed(r3d.KeyI) {
			ssgiOn = !ssgiOn
			r3d.SetSSGI(ssgiOn, 1.5)
		}
		if r3d.IsKeyPressed(r3d.KeyK) {
			procOn = !procOn
			if procOn {
				r3d.SetSky(procSky)
			} else {
				r3d.SetSky(hdrSky)
			}
		}
		if r3d.IsKeyPressed(r3d.KeyV) {
			split = !split
		}
		if r3d.IsKeyPressed(r3d.KeySpace) {
			paused = !paused
		}
		if r3d.IsKeyPressed(r3d.KeyH) {
			showHelp = !showHelp
		}

		if !paused {
			orbit += dt * 0.5
		}

		// --- animate scene ---------------------------------------------------
		r3d.UpdateCamera(&cam, r3d.CameraOrbital)
		player.Update(dt)
		r3d.SetBackgroundRotation(r3d.QuaternionFromAxisAngle(r3d.Vector3{X: 0, Y: 1, Z: 0}, t*0.02))

		for i, l := range omni {
			a := orbit + float32(i)/omniCount*2*math.Pi
			l.SetPosition(r3d.Vector3{X: 5 * cosf(a), Y: 2.5, Z: 5 * sinf(a)})
		}

		for s := 0; s < 10 && len(particles) < maxParticles; s++ {
			ang := float32(s)*0.61 + t*3
			particles = append(particles, particle{
				pos:  emitter,
				vel:  r3d.Vector3{X: cosf(ang) * 0.8, Y: 3.0 + 0.5*sinf(t*7+ang), Z: sinf(ang) * 0.8},
				life: 1.2,
			})
		}
		alive := particles[:0]
		for _, p := range particles {
			p.vel.Y -= 9.81 * dt * 0.35
			p.pos = p.pos.Add(p.vel.Scale(dt))
			p.life -= dt
			if p.life > 0 {
				alive = append(alive, p)
			}
		}
		particles = alive
		for i, p := range particles {
			partPositions[i] = p.pos
		}
		partBuf.UploadPositions(partPositions, len(particles))

		spriteFrame := int(t*8) % spriteFrames
		spriteMat.SetUVScale(1.0/spriteFrames, 1)
		spriteMat.SetUVOffset(float32(spriteFrame)/spriteFrames, 0)

		// --- render ----------------------------------------------------------
		r3d.BeginDrawing()
		if split {
			rw, rh := float32(r3d.GetRenderWidth()), float32(r3d.GetRenderHeight())

			left := r3d.CameraFromRL(cam)
			r3d.BeginPro(r3d.View{Camera: left, Viewport: r3d.Rectangle{X: 0, Y: 0, Width: rw / 2, Height: rh}})
			drawScene(t)
			r3d.End()

			mirror := cam
			mirror.Position = r3d.Vector3{X: -cam.Position.X, Y: cam.Position.Y, Z: -cam.Position.Z}
			right := r3d.CameraFromRL(mirror)
			r3d.BeginPro(r3d.View{Camera: right, Viewport: r3d.Rectangle{X: rw / 2, Y: 0, Width: rw / 2, Height: rh}})
			drawScene(t)
			r3d.End()
		} else {
			r3d.Begin(cam)
			drawScene(t)
			r3d.End()
		}

		uiScale := r3d.GetRenderWidth() / width
		if uiScale < 1 {
			uiScale = 1
		}
		drawHUD(uiScale, hud{showHelp, bloomMode, ssaoOn, fogMode, dofOn, ssrOn, ssgiOn, procOn, split, tonemap, len(particles)})

		r3d.EndDrawing()

		// Optional headless screenshot, taken after EndDrawing so the batched
		// 2D HUD is flushed into the captured frame.
		if shot := os.Getenv("R3D_DEMO_SHOT"); shot != "" && maxFrames > 0 && frame == maxFrames-1 {
			r3d.TakeScreenshot(shot)
		}
	}

	// --- cleanup -------------------------------------------------------------
	probe.Destroy()
	partBuf.Unload()
	ring.Unload()
	player.Unload()
	anims.Unload()
	character.Unload(true)
	helmet.Unload(true)
	procSky.Unload()
	ibl.Unload()
	hdrSky.Unload()
	r3d.Close()
	r3d.CloseWindow()
	fmt.Println("showcase exited cleanly")
}

type hud struct {
	showHelp bool
	bloom    r3d.BloomMode
	ssao     bool
	fog      r3d.FogMode
	dof      bool
	ssr      bool
	ssgi     bool
	proc     bool
	split    bool
	tonemap  r3d.TonemapMode
	parts    int
}

func drawHUD(s int, h hud) {
	r3d.DrawRectangle(0, 0, 380*s, 280*s, r3d.Color{R: 0, G: 0, B: 0, A: 150})
	r3d.DrawText("R3D Go showcase", 12*s, 10*s, 22*s, r3d.RayWhite)
	r3d.DrawFPS(320*s, 12*s)
	y := 44 * s
	line := func(str string, on bool) {
		col := r3d.Color{R: 160, G: 160, B: 160, A: 255}
		if on {
			col = r3d.Green
		}
		r3d.DrawText(str, 12*s, y, 18*s, col)
		y += 24 * s
	}
	line(fmt.Sprintf("[B] Bloom: %s", bloomName(h.bloom)), h.bloom != r3d.BloomDisabled)
	line(fmt.Sprintf("[O] SSAO: %v", h.ssao), h.ssao)
	line(fmt.Sprintf("[I] SSGI: %v", h.ssgi), h.ssgi)
	line(fmt.Sprintf("[R] SSR (reflections): %v", h.ssr), h.ssr)
	line(fmt.Sprintf("[G] Fog: %s", fogName(h.fog)), h.fog != r3d.FogDisabled)
	line(fmt.Sprintf("[F] Depth of Field: %v", h.dof), h.dof)
	line(fmt.Sprintf("[T] Tonemap: %s", tonemapName(h.tonemap)), true)
	line(fmt.Sprintf("[K] Procedural sky: %v", h.proc), h.proc)
	line(fmt.Sprintf("[V] Split-screen: %v", h.split), h.split)
	r3d.DrawText(fmt.Sprintf("particles: %d", h.parts), 12*s, y, 18*s, r3d.RayWhite)

	if h.showHelp {
		r3d.DrawText("[Space] pause lights    [H] hide help    [ESC] quit",
			12*s, r3d.GetScreenHeight()-30*s, 18*s, r3d.RayWhite)
	}
}

func bloomName(b r3d.BloomMode) string {
	return []string{"off", "mix", "additive", "screen"}[int(b)]
}
func fogName(f r3d.FogMode) string {
	return []string{"off", "linear", "exp2", "exp"}[int(f)]
}
func tonemapName(t r3d.TonemapMode) string {
	return []string{"linear", "reinhard", "filmic", "aces", "agx"}[int(t)]
}

func cosf(x float32) float32 { return float32(math.Cos(float64(x))) }
func sinf(x float32) float32 { return float32(math.Sin(float64(x))) }
