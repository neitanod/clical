#requires -Version 5.1
# Compila clical.exe en el directorio actual (equivalente Windows del Makefile 'build').
$ErrorActionPreference = 'Stop'
Set-Location $PSScriptRoot

Write-Host "Compilando clical.exe..."
& go build -o clical.exe ./cmd/clical
if ($LASTEXITCODE -ne 0) { throw "go build fallo con codigo $LASTEXITCODE" }
Write-Host "Compilacion exitosa: $PSScriptRoot\clical.exe"
