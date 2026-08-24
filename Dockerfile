# Stage 1: build the Vue SPA (native to build platform — no QEMU needed)
FROM --platform=$BUILDPLATFORM node:26-alpine@sha256:144769ec3f32e8ee36b3cfde91e82bee25d9367b20f31a151f3f7eea3a2a8541 AS ui-builder
WORKDIR /ui
COPY ui/package*.json ./
RUN npm ci
COPY ui/ ./
RUN npm run build

# Stage 2: build the Go controller binary native, cross-compile for target
FROM --platform=$BUILDPLATFORM golang:1.27@sha256:65b6f280bf050ec5af12716857e8ea8439d694dbba8f31ceeb7630670071f2bb AS builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY api/ ./api/
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY hack/ ./hack/
COPY ui/embed.go ./ui/embed.go

COPY --from=ui-builder /ui/dist ./ui/dist

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -ldflags "-X main.Version=${VERSION}" -o manager cmd/main.go

# Stage 3: minimal runtime image
FROM gcr.io/distroless/static:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
