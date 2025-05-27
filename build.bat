@echo off
setlocal

if "%1"=="" (
    set TARGET=build
) else (
    set TARGET=%1
)

powershell -ExecutionPolicy Bypass -File build.ps1 %TARGET%
