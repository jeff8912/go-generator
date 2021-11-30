@echo off

set program_version=1.0.333
for /f "delims=" %%t in ('go version') do set compiler_version=%%t
set build_time=%DATE% %TIME%
set author=%username%
go build -ldflags "-X 'main.PROGRAM_VERSION=%program_version%' -X 'main.COMPILER_VERSION=%compiler_version%' -X 'main.BUILD_TIME=%build_time%' -X 'main.AUTHOR=%author%'"
