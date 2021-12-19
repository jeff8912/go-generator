@echo off

set programVersion=1.0.0
for /f "delims=" %%t in ('go version') do set compilerVersion=%%t
set buildTime=%DATE% %TIME%
set author=%username%
go build -ldflags "-X 'main.ProgramVersion=%programVersion%' -X 'main.CompileVersion=%compilerVersion%' -X 'main.BuildTime=%buildTime%' -X 'main.Author=%author%'"
