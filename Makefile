NAME ?= eqemupatchergo
VERSION ?= 0.0.3
ICON_PNG ?= icon.png
PACKAGE_NAME ?= com.xackery.eqemupatcher
GITHUB_ACTION_PATH ?= xackery/eqemupatcher

run:
	@-mkdir -p bin
	cd bin && go run ../main.go
run-mobile:
	go run -tags mobile main.go
test:
	go test ./...
bundle:
	@echo "bundle: creating client/bundle.go..."
	echo ${VERSION} > "assets/version.txt"
	echo https://github.com/${GITHUB_ACTION_PATH}/releases/download/latest/ > "assets/update_url.txt"
	fyne bundle --package client -name VersionText assets/version.txt > client/bundle.go
	fyne bundle --package client -name UpdateUrlText --append assets/update_url.txt >> client/bundle.go
	fyne bundle --package client -name NameText --append assets/name.txt >> client/bundle.go
	fyne bundle --package client -name UrlText --append assets/url.txt >> client/bundle.go
	fyne bundle --package client -name RoFImage --append assets/rof.png >> client/bundle.go
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

#go install golang.org/x/tools/cmd/goimports@latest
#go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
#go install golang.org/x/lint/golint@latest
#go install honnef.co/go/tools/cmd/staticcheck@v0.2.2

sanitize:
	rm -rf vendor/
	go vet -tags ci ./...
	test -z $(goimports -e -d . | tee /dev/stderr)
	gocyclo -over 30 .
	golint -set_exit_status $(go list -tags ci ./...)
	staticcheck -go 1.14 ./...
	go test -tags ci -covermode=atomic -coverprofile=coverage.out ./...
    coverage=`go tool cover -func coverage.out | grep total | tr -s '\t' | cut -f 3 | grep -o '[^%]*'`	