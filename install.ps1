#requires -Version 5.1
# Instala clical en el directorio de binarios de Go del usuario (equivalente Windows de 'make install').
# Por defecto Go instala en %USERPROFILE%\go\bin, que conviene tener en PATH.
$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

Write-Host "Instalando clical via 'go install'..."
& go install ./cmd/clical
if ($LASTEXITCODE -ne 0) { throw "go install fallo con codigo $LASTEXITCODE" }

# Resolver destino para mostrarlo al usuario.
$gobin = $env:GOBIN
if ([string]::IsNullOrEmpty($gobin)) {
    $gopath = $env:GOPATH
    if ([string]::IsNullOrEmpty($gopath)) {
        $gopath = Join-Path $env:USERPROFILE 'go'
    }
    $gobin = Join-Path $gopath 'bin'
}

$binary = Join-Path $gobin 'clical.exe'
Write-Host "Instalado en: $binary"

# Avisar si el directorio no esta en PATH.
$paths = $env:Path -split ';' | Where-Object { $_ -ne '' }
$inPath = $paths | Where-Object { $_.TrimEnd('\') -ieq $gobin.TrimEnd('\') }
if (-not $inPath) {
    Write-Warning "El directorio '$gobin' no esta en tu PATH actual."
    Write-Warning "Para usar 'clical' desde cualquier lado agregalo con:"
    Write-Warning "  [Environment]::SetEnvironmentVariable('Path', `"`$env:Path;$gobin`", 'User')"
    Write-Warning "Y abri una terminal nueva."
}
