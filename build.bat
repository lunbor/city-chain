

@echo on

cd %~dp0
%~d0
set batPath=%cd%

cd %batPath%/dots/api_service & go build ./...
cd %batPath%/dots/data/ & go build ./...
cd %batPath%/dots/ & go build ./...
cd %batPath%/servers/chain_server/ & go build ./...

cd %batPath%
