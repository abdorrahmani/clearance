@echo off
powershell -NoExit -Command "Set-Item Property HKCU:\Console VirtualTerminalLevel -Value 1; Start-Process -FilePath 'clearance.exe' -Verb RunAs -PassThru | Wait-Process; exit $LASTEXITCODE" 