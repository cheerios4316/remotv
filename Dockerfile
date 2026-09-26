FROM golang:1.27.1-alpine AS source

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/

FROM source AS client-compile
RUN CGO_ENABLED=0 go build -o /out/bin/client ./cmd/client

FROM scratch AS build-client
COPY --from=client-compile /out/ /

FROM source AS server-compile
RUN CGO_ENABLED=0 go build -o /out/bin/server ./cmd/server

FROM scratch AS build-server
COPY --from=server-compile /out/ /

FROM source AS client-windows-compile
RUN GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /out/bin/client.exe ./cmd/client

FROM scratch AS build-client-windows
COPY --from=client-windows-compile /out/ /

FROM source AS server-windows-compile
RUN GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o /out/bin/server.exe ./cmd/server

FROM scratch AS build-server-windows
COPY --from=server-windows-compile /out/ /
