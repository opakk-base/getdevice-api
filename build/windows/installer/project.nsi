Unicode true
RequestExecutionLevel admin
Name "Pegasus"
!ifndef PEGASUS_EXE
  !define PEGASUS_EXE "bin\Pegasus.exe"
!endif
!ifndef INSTALLER_OUT
  !define INSTALLER_OUT "dist\Pegasus-windows-installer.exe"
!endif
OutFile "${INSTALLER_OUT}"
InstallDir "$PROGRAMFILES64\Pegasus"
SetCompressor /SOLID lzma

Page directory
Page instfiles
UninstPage uninstConfirm
UninstPage instfiles

Section "Install"
  SetOutPath "$INSTDIR"
  File "${PEGASUS_EXE}"
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  CreateDirectory "$SMPROGRAMS\Pegasus"
  CreateShortcut "$SMPROGRAMS\Pegasus\Pegasus.lnk" "$INSTDIR\Pegasus.exe"
  CreateShortcut "$SMPROGRAMS\Pegasus\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
SectionEnd

Section "Uninstall"
  Delete "$SMPROGRAMS\Pegasus\Pegasus.lnk"
  Delete "$SMPROGRAMS\Pegasus\Uninstall.lnk"
  RMDir "$SMPROGRAMS\Pegasus"
  Delete "$INSTDIR\Pegasus.exe"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
SectionEnd
