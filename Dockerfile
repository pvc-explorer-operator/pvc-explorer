# Stage 1: build the Vue SPA
FROM node:26-alpine@sha256:7c6af15abe4e3de859690e7db171d0d711bf37d27528eddfe625b2fe89e097f8 AS ui-builder
WORKDIR /ui
COPY ui/package*.json ./
RUN npm ci
COPY ui/ ./
RUN npm run build

# Stage 2: build the Go controller binary (with embedded UI)
FROM golang:1.26@sha256:2d6c80227255c3112a4d08e67ba98e58efd3846daf15d9d7d4c389565d881b1a AS builder
ARG TARGETOS
ARG TARGETARCH
# VERSION is injected by the build system (make docker-build / GitHub Actions).
# Defaults to "dev" if not provided.
ARG VERSION=dev

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy source explicitly — avoids Podman path-ambiguity with broad COPY . .
COPY api/ ./api/
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY hack/ ./hack/
COPY ui/embed.go ./ui/embed.go

# Overlay built UI assets from the ui-builder stage
COPY --from=ui-builder /ui/dist ./ui/dist

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -ldflags "-X main.Version=${VERSION}" -o manager cmd/main.go

# Stage 3: minimal runtime image
FROM gcr.io/distroless/static:nonroot@sha256:963fa6c544fe5ce420f1f54fb88b6fb01479f054c8056d0f74cc2c6000df5240
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
