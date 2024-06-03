# https://github.com/docker-library/golang/blob/master/1.20/alpine3.18/Dockerfile
FROM public.ecr.aws/docker/library/golang:1.20.11 as builder
ARG VERSION
ENV GOBIN /service/bin

WORKDIR /service

COPY ./api ./api
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./go.mod ./go.sum ./


RUN go install -ldflags "-X main.build=${VERSION}" ./cmd/...

# https://github.com/alpinelinux/docker-alpine/blob/v3.18/x86_64/Dockerfile
FROM public.ecr.aws/docker/library/alpine:3.18.6
RUN apk add --no-cache doas busybox gcompat libgomp libstdc++ openssl \
    && ln -sfv ld-linux-x86-64.so.2 /lib/libresolv.so.2 \
    && adduser -S admin -D -G wheel; \
    echo 'permit nopass :wheel as root' >> /etc/doas.d/doas.conf; \
    chmod g+rx,o+rx /

COPY --from=builder ./service/api ./api
COPY --from=builder ./service/bin/* ./
