# sgv (Simple Go Version) - Windows Installation Script
#
# This script downloads and installs the latest version of sgv for Windows.
# It should be run in a PowerShell terminal with administrator privileges.
#
# Usage:
#   powershell -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/fun7257/sgv/main/install.ps1 | iex"
#

# --- Helper Functions ---

function Write-Color($Text, $Color) {
    Write-Host $Text -ForegroundColor $Color
}

function Info($Message) {
    Write-Color "[INFO] $Message" "Green"
}

function Warn($Message) {
    Write-Color "[WARN] $Message" "Yellow"
}

function Error($Message) {
    Write-Color "[ERROR] $Message" "Red"
    exit 1
}

# --- Main Installation Logic ---

function Main {
    # 1. Check for Administrator Privileges
    if (-NOT ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        Error "This script must be run with Administrator privileges to set environment variables."
    }

    # 2. Detect Architecture
    $archType = $env:PROCESSOR_ARCHITECTURE
    $os = "windows"
    $arch = ""

    if ($archType -eq "AMD64") {
        $arch = "amd64"
    } elseif ($archType -eq "ARM64") {
        $arch = "arm64"
    } else {
        Error "Unsupported architecture: $archType. Only amd64 and arm64 are supported."
    }

    Info "Detected OS: $os, Architecture: $arch"

    # 3. Get the latest version from GitHub Releases
    $repo = "fun7257/sgv"
    $latestReleaseUrl = "https://api.github.com/repos/$repo/releases/latest"
    try {
        $releaseInfo = Invoke-RestMethod -Uri $latestReleaseUrl
        $latestTag = $releaseInfo.tag_name
        $latestVersion = $latestTag.TrimStart('v')
    } catch {
        Error "Could not fetch the latest version tag from GitHub. Please check the repository path and release status."
    }

    Info "Latest version is $latestTag"

    # 4. Download the pre-compiled binary
    $downloadUrl = "https://github.com/$repo/releases/download/$latestTag/sgv_${latestVersion}_${os}_${arch}.zip"
    $tempDir = Join-Path $env:TEMP "sgv_install"
    $zipFile = Join-Path $tempDir "sgv.zip"

    if (Test-Path $tempDir) {
        Remove-Item -Recurse -Force $tempDir
    }
    New-Item -ItemType Directory -Force -Path $tempDir | Out-Null

    Info "Downloading from $downloadUrl..."
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $zipFile
    } catch {
        Error "Download failed. Please check the URL and your network connection."
    }

    # 5. Unpack and install
    $installDir = "$env:ProgramFiles\sgv"
    $sgvExe = Join-Path $installDir "sgv.exe"

    Info "Installing sgv to $installDir..."
    Expand-Archive -Path $zipFile -DestinationPath $installDir -Force

    if (-NOT (Test-Path $sgvExe)) {
        Error "The downloaded archive does not contain 'sgv.exe' in the expected location."
    }

    # 6. Update System Environment Variables
    Info "Updating environment variables..."

    # Set GOROOT
    $goRoot = Join-Path $env:USERPROFILE ".sgv\current"
    Info "Setting GOROOT to $goRoot"
    [Environment]::SetEnvironmentVariable("GOROOT", $goRoot, "Machine")

    # Unset GOPATH
    Info "Unsetting GOPATH"
    try {
        [Environment]::SetEnvironmentVariable("GOPATH", $null, "Machine")
    } catch {}

    # Update PATH
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    $goBin = Join-Path $goRoot "bin"
    $userGoBin = Join-Path $env:USERPROFILE "go\bin"

    $pathItems = $currentPath -split ';' | Where-Object { $_ -ne '' }
    $newPathItems = @()

    # Add our paths if they don't exist
    if ($pathItems -notcontains $installDir) {
        $newPathItems += $installDir
    }
    if ($pathItems -notcontains $goBin) {
        $newPathItems += $goBin
    }
    if ($pathItems -notcontains $userGoBin) {
        $newPathItems += $userGoBin
    }

    $newPathItems += $pathItems
    $newPath = $newPathItems -join ';'

    if ($newPath -ne $currentPath) {
        Info "Adding sgv and Go paths to the system PATH."
        [Environment]::SetEnvironmentVariable("Path", $newPath, "Machine")
    } else {
        Info "System PATH already seems to be configured. Skipping."
    }

    # Cleanup
    Remove-Item -Recurse -Force $tempDir

    # --- Final Instructions ---
    Info "Installation successful!"
    Write-Color "sgv has been installed to: $sgvExe" "Yellow"
    Write-Color "Please restart your terminal for the environment variable changes to take effect." "Yellow"
}

# Run the main function
Main
