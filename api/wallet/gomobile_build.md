
//build output of gomobile bind
GOOS=android CGO_ENABLED=1 $GOPATH\bin\gobind.exe -lang=go,java -outdir=$WORK github.com/scryinfo/citychain/api/wallet
GOOS=android GOARCH=arm CC=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\armv7a-linux-androideabi16-clang CXX=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\armv7a-linux-androideabi16-clang++ CGO_ENABLED=1 GOARM=7 GOPATH=$WORK;$GO
PATH GO111MODULE=off go build -x -buildmode=c-shared -o=$WORK\android\src\main\jniLibs\armeabi-v7a\libgojni.so gobind
GOOS=android GOARCH=arm64 CC=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang CXX=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang++ CGO_ENABLED=1 GOPATH=$WORK;$GOPATH GO111MO
DULE=off go build -x -buildmode=c-shared -o=$WORK\android\src\main\jniLibs\arm64-v8a\libgojni.so gobind
GOOS=android GOARCH=386 CC=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\i686-linux-android16-clang CXX=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\i686-linux-android16-clang++ CGO_ENABLED=1 GOPATH=$WORK;$GOPATH GO111MODULE=off
 go build -x -buildmode=c-shared -o=$WORK\android\src\main\jniLibs\x86\libgojni.so gobind
GOOS=android GOARCH=amd64 CC=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\x86_64-linux-android21-clang CXX=$ANDROID_HOME\ndk-bundle\toolchains\llvm\prebuilt\windows-x86_64\bin\x86_64-linux-android21-clang++ CGO_ENABLED=1 GOPATH=$WORK;$GOPATH GO111MODU
LE=off go build -x -buildmode=c-shared -o=$WORK\android\src\main\jniLibs\x86_64\libgojni.so gobind
PWD=$WORK\java javac -d $WORK\javac-output -source 1.7 -target 1.7 -bootclasspath $ANDROID_HOME\platforms\android-29\android.jar go\Seq.java go\Universe.java go\error.java wallet\Wallet.java
jar c -C $WORK\javac-output 

//output file 
