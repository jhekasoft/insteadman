[Setup]
AppName=InsteadMan
AppVersion=3.5
DefaultDirName={commonpf}\InsteadMan
DefaultGroupName=Games
UninstallDisplayIcon={app}\insteadman-fyne.exe
OutputDir=.
OutputBaseFilename=insteadman-windows-setup
AllowNoIcons=true
SetupIconFile=.\resources\images\logo.ico
ChangesAssociations=yes

[Languages]
Name: en; MessagesFile: compiler:Default.isl
Name: ru; MessagesFile: compiler:Languages\Russian.isl
Name: uk; MessagesFile: compiler:Languages\Ukrainian.isl

[Files]
Source: .\insteadman\*; DestDir: {app}; Flags: replacesameversion recursesubdirs

[CustomMessages]
CreateDesktopIcon=Create a &desktop icon
LaunchGame=Launch &InsteadMan
UninstallMsg=Uninstall InsteadMan
RmSettingsMsg=Would you like to remove settings?

[Tasks]
Name: desktopicon; Description: {cm:CreateDesktopIcon}

[Run]
Filename: {app}\insteadman-fyne.exe; Description: {cm:LaunchGame}; WorkingDir: {app}; Flags: postinstall

[Icons]
Name: {commondesktop}\InsteadMan; Filename: {app}\insteadman-fyne.exe; WorkingDir: {app}; Tasks: desktopicon
Name: {group}\InsteadMan; Filename: {app}\insteadman-fyne.exe; WorkingDir: {app}
Name: {group}\{cm:UninstallMsg}; Filename: {uninstallexe}

[UninstallDelete]
Name: {app}; Type: dirifempty

[Code]
procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  case CurUninstallStep of
    usPostUninstall:
      begin
        if MsgBox(CustomMessage('RmSettingsMsg'), mbConfirmation, MB_YESNO or MB_DEFBUTTON2) = idYes then
        begin
          // remove settings and saved games manually
          DelTree(ExpandConstant('{localappdata}\insteadman'), True, True, True);
        end;
      end;
  end;
end;