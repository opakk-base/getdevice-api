Unicode true

####
## GetDevice NSIS installer.
##
## Based on the Wails v2 template. Customizations:
##   - PRODUCT_EXECUTABLE pinned to GetDevice.exe (wails.json outputfilename)
##   - Optional desktop shortcut + "run after install" on the finish page
##   - Uninstaller asks whether to remove user config (~\AppData\Roaming\GetDevice)
##
## Wails generates wails_tools.nsh next to this file during
## `wails build -platform windows/amd64 -nsis`, then invokes makensis with:
##   -DARG_WAILS_AMD64_BINARY=..\..\bin\GetDevice.exe
####
## The following defines may be overridden here (wails_tools.nsh skips them
## if already defined):
## !define INFO_PROJECTNAME    "MyProject"
## !define INFO_COMPANYNAME    "MyCompany"
## !define INFO_PRODUCTNAME    "MyProduct"
## !define INFO_PRODUCTVERSION "1.0.0"
## !define INFO_COPYRIGHT      "Copyright"
####

# Pin the executable name to the actual build output (GetDevice.exe);
# wails_tools.nsh would otherwise assume "${INFO_PROJECTNAME}.exe".
!define PRODUCT_EXECUTABLE "GetDevice.exe"

!include "wails_tools.nsh"

# The version information for these two must consist of 4 parts
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

# Enable HiDPI support. https://nsis.sourceforge.io/Reference/ManifestDPIAware
ManifestDPIAware true

!include "MUI.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_ABORTWARNING           # Warn when the user aborts the installer

# Finish page: optionally run the app and/or create a desktop shortcut.
!define MUI_FINISHPAGE_RUN "$INSTDIR\${PRODUCT_EXECUTABLE}"
!define MUI_FINISHPAGE_RUN_TEXT "Run ${INFO_PRODUCTNAME} now"
!define MUI_FINISHPAGE_SHOWREADME ""
!define MUI_FINISHPAGE_SHOWREADME_NOTCHECKED
!define MUI_FINISHPAGE_SHOWREADME_TEXT "Create a desktop shortcut"
!define MUI_FINISHPAGE_SHOWREADME_FUNCTION CreateDesktopShortcut

!insertmacro MUI_PAGE_WELCOME     # Welcome page
# !insertmacro MUI_PAGE_LICENSE "resources\eula.txt"
!insertmacro MUI_PAGE_DIRECTORY   # Install location page
!insertmacro MUI_PAGE_INSTFILES   # Installing page
!insertmacro MUI_PAGE_FINISH      # Finish page

!insertmacro MUI_UNPAGE_CONFIRM   # Confirm before uninstalling
!insertmacro MUI_UNPAGE_INSTFILES # Uninstalling page

!insertmacro MUI_LANGUAGE "English"

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe" # Installer output file.

!ifdef WAILS_INSTALL_SCOPE
  !if "${WAILS_INSTALL_SCOPE}" == "user"
    InstallDir "$LOCALAPPDATA\Programs\${INFO_PRODUCTNAME}"
  !else
    InstallDir "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
  !endif
!else
  InstallDir "$PROGRAMFILES64\${INFO_COMPANYNAME}\${INFO_PRODUCTNAME}"
!endif

ShowInstDetails show   # Always show installation details
ShowUnInstDetails show # Always show uninstall details

Function .onInit
   !insertmacro wails.checkArchitecture
FunctionEnd

Function CreateDesktopShortcut
  CreateShortcut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"
FunctionEnd

Section
    !insertmacro wails.setShellContext

    !insertmacro wails.webview2runtime

    SetOutPath $INSTDIR

    !insertmacro wails.files

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\${PRODUCT_EXECUTABLE}"

    !insertmacro wails.associateFiles
    !insertmacro wails.associateCustomProtocols

    !insertmacro wails.writeUninstaller
SectionEnd

Section "un.install"
    !insertmacro wails.setShellContext

    MessageBox MB_YESNO|MB_ICONQUESTION \
        "Also remove your GetDevice settings and device identity?$\r$\n$\r$\n(DEVICE_ID and CLIENT_KEY are lost if you choose Yes.)" \
        /SD IDNO IDNO keep_config

    RMDir /r "$AppData\GetDevice"          # Config directory (.env written by the app)
keep_config:

    RMDir /r "$LOCALAPPDATA\GetDevice"     # WebView2 user data folder

    RMDir /r "$INSTDIR"

    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"

    !insertmacro wails.unassociateFiles
    !insertmacro wails.unassociateCustomProtocols

    !insertmacro wails.deleteUninstaller    # Also removes the ARP registry key
SectionEnd
