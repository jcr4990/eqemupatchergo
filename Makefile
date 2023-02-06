NAME ?= eqemupatchergo
VERSION ?= 0.0.2
ICON_PNG ?= icon.png
PACKAGE_NAME ?= com.xackery.eqemupatcher

run:
	@-mkdir -p bin
	cd bin && go run ../main.go
run-mobile:
	go run -tags mobile main.go
test:
	go test ./...
bundle:
	@echo "bundle: creating client/bundle.go..."
	fyne bundle --package client -name VersionText assets/version.txt > client/bundle.go
	fyne bundle --package client -name NameText --append assets/name.txt >> client/bundle.go
	fyne bundle --package client -name UrlText --append assets/url.txt >> client/bundle.go
	fyne bundle --package client -name RoFImage --append assets/rof.png >> client/bundle.go
	echo ${VERSION} > "assets/version.txt"
build-all: build-darwin build-ios build-linux build-windows build-android
build-darwin:
	@echo "build-darwin: compiling"
	@-mkdir -p bin
	@-rm -rf bin/${NAME}-darwin.zip
	@-rm -rf bin/orcspawn.app
	@time fyne package -os darwin -icon ${ICON_PNG} --appVersion ${VERSION} --tags main.Version=${VERSION}
	@zip -mvr bin/${NAME}-${VERSION}-darwin.zip ${NAME}.app -x "*.DS_Store"
build-linux:
	@echo "build-linux: compiling"
	@-mkdir -p bin
	@-rm -rf bin/${NAME}-linux
	@time fyne-cross linux -icon ${ICON_PNG}
	@mv fyne-cross/bin/linux-amd64/${NAME} bin/${NAME}-linux
	@-rm -rf fyne-cross/
build-windows:
	@echo "build-windows: compiling"
	-mkdir -p bin
	-rm bin/${NAME}-*-windows.zip
	fyne-cross windows -icon ${ICON_PNG}
	mv fyne-cross/bin/windows-amd64/${NAME}.exe bin/
	@#-rm -rf fyne-cross/
	@#cd bin && zip -mv ${NAME}-${VERSION}-windows.zip ${NAME}.exe
build-ios:
	@echo "build-ios: compiling"
	@-mkdir -p bin
	@-rm -rf bin/${NAME}-ios.zip
	@DISABLE_MANUAL_TARGET_ORDER_BUILD_WARNING=1 time fyne package -os ios -appID ${PACKAGE_NAME} -icon ${ICON_PNG}
	@zip -mvr bin/${NAME}-ios.zip ${NAME}.app -x "*.DS_Store"
build-android:
	@echo "build-android: compiling"
	@-mkdir -p bin
	@-rm bin/${NAME}.apk
	@ANDROID_NDK_HOME=~/Library/Android/sdk/ndk-bundle fyne package -os android -appID ${PACKAGE_NAME} -icon ${ICON_PNG}
	@mv ${NAME}.apk bin/${NAME}.apk
build-web:
	@echo "build-web: compiling"
	@-mkdir -p bin
	@#-rm -rf bin/${NAME}-darwin.zip
	@time fyne package -os web -icon ${ICON_PNG}
	@#zip -mvr bin/${NAME}-darwin.zip ${NAME}.app -x "*.DS_Store"