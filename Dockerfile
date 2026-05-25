# Stage 1: build the Vue SPA
FROM node:26-alpine AS ui-builder
WORKDIR /ui
COPY ui/package*.json ./
RUN npm ci
COPY ui/ ./
RUN npm run build

# Stage 2: build the Go controller binary (with embedded UI)
FROM golang:1.26 AS builder
ARG TARGETOS
ARG TARGETARCH
# VERSION is injected by the build system (make docker-build / GitHub Actions).
# Defaults to "dev" if not provided.
ARG VERSION=dev

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
COPY --from=ui-builder /ui/dist ./ui/dist

RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -ldflags "-X main.Version=${VERSION}" -o manager cmd/main.go

# Stage 3: minimal runtime image
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
