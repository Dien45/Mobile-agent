$ErrorActionPreference = "Stop"
$GoVersion = "1.24.2"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
$Arch = if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -eq "Arm64") { "arm64" } else { "amd64" }
$Go = Get-Command go -ErrorAction SilentlyContinue

if (-not $Go) {
  $Cache = Join-Path $env:LOCALAPPDATA "Arka\toolchains"
  New-Item -ItemType Directory -Force -Path $Cache | Out-Null
  $Archive = "go$GoVersion.windows-$Arch.zip"
  $Url = "https://go.dev/dl/$Archive"
  $Zip = Join-Path $Cache $Archive
  Write-Host "Downloading pinned Go $GoVersion for windows/$Arch..."
  Invoke-WebRequest $Url -OutFile $Zip
  $Expected = (Invoke-WebRequest "$Url.sha256").Content.Trim()
  $Actual = (Get-FileHash -Algorithm SHA256 $Zip).Hash.ToLowerInvariant()
  if ($Expected.ToLowerInvariant() -ne $Actual) { throw "Go checksum verification failed" }
  $GoDir = Join-Path $Cache "go"
  Remove-Item -Recurse -Force -ErrorAction SilentlyContinue $GoDir
  Expand-Archive -Path $Zip -DestinationPath $Cache -Force
  $GoExe = Join-Path $GoDir "bin\go.exe"
} else {
  $GoExe = $Go.Source
}

$Bin = Join-Path $env:LOCALAPPDATA "Arka\bin"
New-Item -ItemType Directory -Force -Path $Bin | Out-Null
Write-Host "Building Arka locally..."
Push-Location $Root
try { & $GoExe build -trimpath -ldflags "-s -w" -o (Join-Path $Bin "arka.exe") ./cmd/arka } finally { Pop-Location }
if ($LASTEXITCODE -ne 0) { throw "Arka build failed" }

$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (($userPath -split ';') -notcontains $Bin) {
  [Environment]::SetEnvironmentVariable("Path", (($userPath.TrimEnd(';') + ';' + $Bin).TrimStart(';')), "User")
  Write-Host "Added $Bin to your user PATH. Open a new terminal if needed."
}
Write-Host "Installed: $Bin\arka.exe"
Write-Host "Run: arka start"
