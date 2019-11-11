
#chcp 936
#chcp 65001
set JAVA_TOOL_OPTIONS=-Dfile.encoding=utf-8 -Duser.language=en -Duser.country=US
go clean
gomobile bind -target=android github.com/scryinfo/citychain/api/wallet
