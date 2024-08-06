# Handling QR Code Output in PowerShell on Windows
=====================================================

## Issue
-----

When running `nftfetch wallet` on Windows using PowerShell, QR code output may not render correctly due to VT/ANSI escape sequences not being supported by default in the console. This problem does not occur in environments like Visual Studio Code's integrated terminal or Cygwin, where VT/ANSI escape sequences are supported by default.

## Solution Options
---------------

### Option 1: Enable VT/ANSI Escape Sequences Support
Enable VT/ANSI escape sequences support globally in Windows by modifying the registry.

#### Registry Method
```
Set-ItemProperty HKCU:\Console VirtualTerminalLevel -Type DWORD 1
```
or
```
reg add HKCU\Console /v VirtualTerminalLevel /t REG_DWORD /d 1
```
Open a new console window for the changes to take effect.
Note: This method activates VT support globally, which may affect the output of other programs.

### Programmatic Method
Use the `SetConsoleMode()` Windows API function for a programmatic approach. This can be complex and may not be feasible for all programs or languages.

### Ad-hoc Workaround
Enclose the external program's calls in PowerShell using parentheses:
```
(.\nftfetch wallet)
```
or pipe the output to `Out-Host`:
```
.\nftfetch wallet | Out-Host
```