<#
.SYNOPSIS
  apikit installer (Windows, amd64).

.DESCRIPTION
  Downloads a release archive from GitHub, verifies its SHA-256 against
  checksums.txt, installs apikit.exe, and adds the install directory to the
  user PATH.

.PARAMETER Version
  Release to install, e.g. v0.2.0. Default: latest. (env: APIKIT_VERSION)

.PARAMETER InstallDir
  Install directory. Default: %LOCALAPPDATA%\apikit\bin. (env: APIKIT_INSTALL_DIR)

.EXAMPLE
  irm https://raw.githubusercontent.com/dkryvak/apikit/main/install.ps1 | iex

.EXAMPLE
  .\install.ps1 -Version v0.2.0 -InstallDir C:\tools\apikit
#>
[CmdletBinding()]
param(
  [string]$Version = $env:APIKIT_VERSION,
  [string]$InstallDir = $env:APIKIT_INSTALL_DIR
)

$ErrorActionPreference = 'Stop'
$RepoOwner = 'dkryvak'
$RepoName  = 'apikit'
$Bin       = 'apikit'

function Info($m) { Write-Host "==> $m" }

# --- arch check ------------------------------------------------------------
$arch = $env:PROCESSOR_ARCHITECTURE
if ($arch -ne 'AMD64') {
  throw "unsupported architecture: $arch (only windows amd64 is published)"
}
$Os = 'windows'
$Arch = 'amd64'

# --- resolve version -------------------------------------------------------
if ([string]::IsNullOrEmpty($Version)) {
  Info 'resolving latest release'
  $api = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
  $rel = Invoke-RestMethod -Uri $api -Headers @{ 'User-Agent' = 'apikit-installer' }
  $Version = $rel.tag_name
  if ([string]::IsNullOrEmpty($Version)) { throw 'could not determine latest version' }
}
if ($Version -notlike 'v*') { $Version = "v$Version" }
# GoReleaser strips the leading 'v' for asset names (.Version).
$VerNum = $Version.TrimStart('v')

$Asset = "${Bin}_${VerNum}_${Os}_${Arch}.zip"
$Base  = "https://github.com/$RepoOwner/$RepoName/releases/download/$Version"

if ([string]::IsNullOrEmpty($InstallDir)) {
  $InstallDir = Join-Path $env:LOCALAPPDATA 'apikit\bin'
}

# --- download + verify -----------------------------------------------------
$tmp = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $tmp | Out-Null
try {
  $zip = Join-Path $tmp $Asset
  $sums = Join-Path $tmp 'checksums.txt'

  Info "downloading $Asset ($Version)"
  Invoke-WebRequest -Uri "$Base/$Asset" -OutFile $zip -UseBasicParsing
  Invoke-WebRequest -Uri "$Base/checksums.txt" -OutFile $sums -UseBasicParsing

  Info 'verifying sha256'
  $expected = (Get-Content $sums | Where-Object { $_ -match [regex]::Escape($Asset) + '$' } |
               ForEach-Object { ($_ -split '\s+')[0] } | Select-Object -First 1)
  if ([string]::IsNullOrEmpty($expected)) { throw "no checksum entry for $Asset" }
  $actual = (Get-FileHash -Algorithm SHA256 -Path $zip).Hash.ToLower()
  if ($actual -ne $expected.ToLower()) {
    throw "checksum mismatch (expected $expected, got $actual)"
  }

  # --- extract + install ---------------------------------------------------
  Expand-Archive -Path $zip -DestinationPath $tmp -Force
  New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
  $dest = Join-Path $InstallDir "$Bin.exe"
  Copy-Item -Path (Join-Path $tmp "$Bin.exe") -Destination $dest -Force
  Info "installed $Bin $Version -> $dest"

  # --- user PATH -----------------------------------------------------------
  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  if (($userPath -split ';') -notcontains $InstallDir) {
    [Environment]::SetEnvironmentVariable('Path', "$userPath;$InstallDir", 'User')
    Info "added $InstallDir to your user PATH (restart your shell to pick it up)"
  }
}
finally {
  Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
}
