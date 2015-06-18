cd browser
go build -o ../Release/browser.exe
cd ..

@if exist "main-res.syso" (
    @del "main-res.syso"
)

@if exist "%~dp0Release\example1.exe" (
    @del "%~dp0Release\example1.exe"
)

cd examples\example1

windres -o main-res.syso main.rc


IF "%1"=="noconsole" (
    go build -ldflags="-H windowsgui" -o ../../Release/example1.exe
    rem @if %ERRORLEVEL% neq 0 goto end
) else (
    go build -o ../../Release/example1.exe
    rem @if %ERRORLEVEL% neq 0 goto end
)

cd ../../Release
example1.exe
cd ..

pause
