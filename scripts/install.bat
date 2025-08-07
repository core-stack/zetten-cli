@echo off
setlocal enabledelayedexpansion

set CLI_NAME=zetten-cli.exe
set SERVICE_NAME=zetten-service.exe
set INSTALL_DIR=%ProgramFiles%\Zetten
set ROOT_DIR=%USERPROFILE%\.zetten
set ROOT_CONFIG=%ROOT_DIR%\config.yml

echo Installing Zetten CLI and Service...

IF NOT EXIST "%CLI_NAME%" (
  echo %CLI_NAME% not found.
  exit /b 1
)

IF NOT EXIST "%SERVICE_NAME%" (
  echo %SERVICE_NAME% not found.
  exit /b 1
)

echo Creating install directory: %INSTALL_DIR%
mkdir "%INSTALL_DIR%" >nul 2>&1

echo Copying binaries to %INSTALL_DIR%...
copy /Y "%CLI_NAME%" "%INSTALL_DIR%\"
copy /Y "%SERVICE_NAME%" "%INSTALL_DIR%\"

echo Adding install directory to PATH...
setx PATH "%PATH%;%INSTALL_DIR%" >nul

IF NOT EXIST "%ROOT_DIR%" (
  echo Creating configuration directory: %ROOT_DIR%
  mkdir "%ROOT_DIR%"
)

IF NOT EXIST "%ROOT_CONFIG%" (
  echo Creating config root default in %ROOT_CONFIG%
  echo zettenProjects: [] > "%ROOT_CONFIG%"
  echo mirror: [] >> "%ROOT_CONFIG%"
)

echo Installing service %SERVICE_NAME%...
"%INSTALL_DIR%\%SERVICE_NAME%" install

echo Starting service %SERVICE_NAME%...
"%INSTALL_DIR%\%SERVICE_NAME%" start

echo.
echo Installation complete!
echo Configuration in: %ROOT_CONFIG%
echo Use "%SERVICE_NAME% [start|stop|uninstall]" to manage the service.

endlocal
