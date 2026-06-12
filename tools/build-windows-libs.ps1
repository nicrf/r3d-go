#requires -Version 5.1
<#
.SYNOPSIS
    Build the vendored Windows/amd64 native archives for r3d-go.

.DESCRIPTION
    Fetches the r3dStarter CMake project (which pulls raylib 5.5, R3D and
    Assimp via FetchContent), builds the three static libraries with MinGW-w64,
    and copies them — plus Assimp's static zlib — into lib/windows_amd64/.

    The libraries MUST be built with the SAME MinGW-w64 gcc that cgo links
    with (w64devkit by default), otherwise Assimp's libstdc++ ABI won't match
    at link time. Do not mix MinGW distributions.

    Requirements on the machine:
      - A MinGW-w64 toolchain (gcc/g++/mingw32-make), e.g. w64devkit.
      - CMake >= 3.15.
      - Python 3 (R3D embeds its shaders at configure time).
      - git.

.PARAMETER W64DevkitBin
    Path to the MinGW-w64 bin directory (must contain gcc.exe, g++.exe,
    mingw32-make.exe). Default: C:\w64devkit\w64devkit\bin if present, else
    whatever gcc.exe is found on PATH.

.PARAMETER BuildDir
    Scratch directory for the CMake build. Default: <repo>\.build-win
    (git-ignored). Pass -Clean to wipe it first.

.PARAMETER Clean
    Remove the build directory before configuring (forces a full rebuild).

.EXAMPLE
    pwsh tools/build-windows-libs.ps1

.EXAMPLE
    pwsh tools/build-windows-libs.ps1 -W64DevkitBin C:\tools\w64devkit\bin -Clean
#>
[CmdletBinding()]
param(
    [string]$W64DevkitBin,
    [string]$BuildDir,
    [switch]$Clean
)

$ErrorActionPreference = 'Stop'

function Find-Tool {
    param([string]$Name, [string[]]$Candidates)
    foreach ($c in $Candidates) { if ($c -and (Test-Path $c)) { return (Resolve-Path $c).Path } }
    $cmd = Get-Command $Name -ErrorAction SilentlyContinue
    if ($cmd) { return $cmd.Source }
    throw "Could not locate '$Name'. Pass it explicitly or add it to PATH."
}

# --- Resolve the repo root (this script lives in <repo>/tools/) ------------
$repo = Split-Path -Parent $PSScriptRoot
$libOut = Join-Path $repo 'lib\windows_amd64'

# --- Locate the toolchain --------------------------------------------------
if (-not $W64DevkitBin) {
    $default = 'C:\w64devkit\w64devkit\bin'
    if (Test-Path $default) {
        $W64DevkitBin = $default
    } else {
        $gcc = Get-Command gcc -ErrorAction SilentlyContinue
        if (-not $gcc) { throw "No MinGW gcc found. Install w64devkit or pass -W64DevkitBin." }
        $W64DevkitBin = Split-Path -Parent $gcc.Source
    }
}
foreach ($exe in 'gcc.exe','g++.exe','mingw32-make.exe') {
    if (-not (Test-Path (Join-Path $W64DevkitBin $exe))) {
        throw "$exe not found in $W64DevkitBin — is this a complete MinGW-w64 toolchain?"
    }
}

$cmake  = Find-Tool 'cmake'  @("C:\Program Files\CMake\bin\cmake.exe")
$python = Find-Tool 'python' @("$env:LOCALAPPDATA\Programs\Python\Python313\python.exe")
# Prefer a real interpreter over the Windows Store alias.
if ($python -like '*WindowsApps*') {
    $real = (& py -c "import sys; print(sys.executable)") 2>$null
    if ($real -and (Test-Path $real)) { $python = $real }
}
$pythonDir = Split-Path -Parent $python

if (-not $BuildDir) { $BuildDir = Join-Path $repo '.build-win' }
$src = Join-Path $BuildDir 'r3dStarter'

Write-Host "repo       : $repo"
Write-Host "toolchain  : $W64DevkitBin"
Write-Host "cmake      : $cmake"
Write-Host "python     : $python"
Write-Host "build dir  : $BuildDir"
Write-Host ""

if ($Clean -and (Test-Path $BuildDir)) {
    Write-Host "Cleaning $BuildDir ..."
    Remove-Item -Recurse -Force $BuildDir
}
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

# --- Isolated PATH: ONLY the MinGW we want, plus CMake and Python. ---------
# Crucially excludes git-bash's sh.exe (breaks the MinGW Makefiles generator)
# and any other gcc on the system.
$cmakeDir = Split-Path -Parent $cmake
$env:PATH = "$W64DevkitBin;$cmakeDir;$pythonDir;C:\Windows\System32;C:\Windows"
$env:CC  = 'gcc'
$env:CXX = 'g++'

# --- Fetch r3dStarter ------------------------------------------------------
if (-not (Test-Path $src)) {
    Write-Host "Cloning r3dStarter ..."
    & git clone --depth 1 https://github.com/jensroth-git/r3dStarter $src
    if ($LASTEXITCODE -ne 0) { throw "git clone failed ($LASTEXITCODE)" }
}

$cmBuild = Join-Path $src 'build'

# --- Configure -------------------------------------------------------------
# Image decoders (PNG/HDR/JPG/QOI) are re-enabled (r3dStarter turns them off);
# `-include assert.h` works around raylib rmodels.c using assert() without it.
Write-Host "Configuring (CMake / MinGW Makefiles) ..."
& $cmake -S $src -B $cmBuild -G "MinGW Makefiles" `
    -DCMAKE_BUILD_TYPE=Release `
    -DCMAKE_C_COMPILER=gcc.exe `
    -DCMAKE_CXX_COMPILER=g++.exe `
    -DCMAKE_MAKE_PROGRAM=mingw32-make.exe `
    "-DPython3_EXECUTABLE=$python" `
    "-DPYTHON_EXECUTABLE=$python" `
    -DSUPPORT_FILEFORMAT_PNG=ON `
    -DSUPPORT_FILEFORMAT_HDR=ON `
    -DSUPPORT_FILEFORMAT_JPG=ON `
    -DSUPPORT_FILEFORMAT_QOI=ON `
    "-DCMAKE_C_FLAGS=-Wno-error -Wno-implicit-function-declaration -include assert.h" `
    "-DCMAKE_CXX_FLAGS=-Wno-error"
if ($LASTEXITCODE -ne 0) {
    throw "CMake configure failed. If it said 'Permission denied' reading the " +
          "compiler-id test exe, that's Windows Defender locking the fresh " +
          "binary — just re-run this script."
}

# --- Build the three library targets (NOT the template exe, whose src/main.cpp
#     can lag behind R3D's API — the libraries themselves still build). ------
Write-Host "Building raylib, r3d, assimp ..."
& $cmake --build $cmBuild --target raylib r3d assimp -j
if ($LASTEXITCODE -ne 0) { throw "CMake build failed ($LASTEXITCODE)" }

# --- Collect the archives --------------------------------------------------
$artifacts = @{
    'libr3d.a'        = Join-Path $cmBuild 'lib\libr3d.a'
    'libraylib.a'     = Join-Path $cmBuild '_deps\raylib-build\raylib\libraylib.a'
    'libassimp.a'     = Join-Path $cmBuild '_deps\assimp-build\lib\libassimp.a'
    'libzlibstatic.a' = Join-Path $cmBuild '_deps\assimp-build\contrib\zlib\libzlibstatic.a'
}
New-Item -ItemType Directory -Force -Path $libOut | Out-Null
foreach ($name in $artifacts.Keys) {
    $from = $artifacts[$name]
    if (-not (Test-Path $from)) { throw "Expected artifact not found: $from" }
    Copy-Item -Force $from (Join-Path $libOut $name)
    Write-Host ("  -> lib/windows_amd64/{0}  ({1:N0} bytes)" -f $name, (Get-Item $from).Length)
}

Write-Host ""
Write-Host "Done. Now build the demo:  go build -o r3ddemo.exe ./cmd/demo"
