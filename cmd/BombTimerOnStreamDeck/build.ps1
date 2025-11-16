<#
.Description
Build and deploy file for running the plugin locally
#>

param(
    [switch]$build
)

$destFolderName = "uk.co.ltroberts.streamdeck.bombtimer" # do NOT include the `.sdPlugin` sufix here
$buildDir = Join-Path -Path $env:APPDATA -ChildPath "Elgato\StreamDeck\Plugins\$destFolderName.sdPlugin"

$localBinaryName = "BombTimerOnStreamDeck.exe"  # this should match what is in your manifsets "codePath" field
$localManifsetName = "manifest.json"
$localImagesFolder = "images"

# Build directory if exists, else do nothing, don't need to both reporting this to user tho
New-Item -Path $buildDir -ItemType Directory -Force | Out-Null


class BundledFile {
    [string]$sourceName
    [string]$destinationNameWithPath
    [switch]$Recursive = $false
}

# Any other files you need to copy to the plugin directory should be added here
# Dont forget to add comma as you cannot leave a trailing comma
$filesToCopy = @(
    [BundledFile]@{sourceName = $localManifsetName; destinationNameWithPath = (Join-Path $buildDir "manifest.json")},
    [BundledFile]@{sourceName = $localBinaryName; destinationNameWithPath = (Join-Path $buildDir $localBinaryName)},
    [BundledFile]@{sourceName = $localImagesFolder; destinationNameWithPath = $buildDir; Recursive = $true}
)

if ($build) {
    Write-Host "Received build flag, building binary..."
    go build -o $localBinaryName
} else {

    Write-Host "Skipping building the binary"
}

# Wipe the directory so there are no stale files there.
# Need something there? Add it to the copy loop above
Write-Host "Clearing out $buildDir"
Remove-Item (Join-Path $buildDir "*") -Recurse -Force

foreach ($fileToCopy in $filesToCopy){
    $copyParams = @{
        Path = $fileToCopy.sourceName
        Destination = $fileToCopy.destinationNameWithPath
        Force = $true
        Recurse = $fileToCopy.Recursive
    }
    
    Copy-Item @copyParams
    Write-Host "Successfully copied $($fileToCopy.sourceName) to $($fileToCopy.destinationNameWithPath)"
}

