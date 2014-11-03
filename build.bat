@if exist "main-res.syso" (
    @del "main-res.syso"
)

@if exist "%~dp0Release\ngui.exe" (
    @del "%~dp0Release\ngui.exe"
)

windres -o main-res.syso main.rc
rem go build -o Release\ngui.exe


IF "%1"=="noconsole" (
    go build -ldflags="-H windowsgui" -o Release/ngui.exe
    rem @if %ERRORLEVEL% neq 0 goto end
) else (
    go build -o Release/ngui.exe
    rem @if %ERRORLEVEL% neq 0 goto end
)

cd Release
ngui.exe
cd ..

pause
