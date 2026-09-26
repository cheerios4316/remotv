.PHONY: build build-client build-server build-windows build-client-windows build-server-windows

.PHONY: build-firefox sign-firefox

build-firefox:
	npm --prefix extensions/firefox ci
	npm --prefix extensions/firefox run build

sign-firefox:
	npm --prefix extensions/firefox ci
	npm --prefix extensions/firefox run sign

build: build-client build-server

build-client:
	docker build --target build-client --output type=local,dest=. .

build-server:
	docker build --target build-server --output type=local,dest=. .

build-windows: build-client-windows build-server-windows

build-client-windows:
	docker build --target build-client-windows --output type=local,dest=. .

build-server-windows:
	docker build --target build-server-windows --output type=local,dest=. .
