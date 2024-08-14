FROM node:18.16.1 AS build-resources

WORKDIR /

COPY resources/ resources
COPY static static
COPY views views

RUN cd resources && npm ci
RUN cd resources && npm run build-css

FROM golang:1.22 AS build-worker

ARG COMMIT_SHA=0.0.1

WORKDIR /

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=$COMMIT_SHA" -mod=readonly -v -o worker cmd/worker/main.go

FROM golang:1.22 AS build-app

ARG COMMIT_SHA=0.0.1

WORKDIR /

COPY . .

COPY --from=build-resources static static

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=$COMMIT_SHA" -mod=readonly -v -o app cmd/app/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=$COMMIT_SHA" -mod=readonly -v -o healthcheck cmd/healthcheck/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=$COMMIT_SHA" -mod=readonly -v -o migrate cmd/migrate/main.go

FROM scratch AS worker

WORKDIR /

COPY --from=build-worker /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-worker worker worker

CMD ["./worker"]

FROM scratch AS app

WORKDIR /

COPY --from=build-app /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-app app app
COPY --from=build-app migrate migrate
COPY --from=build-app migrations migrations
COPY --from=build-app healthcheck healthcheck
COPY --from=build-app static static 
COPY --from=build-app resources/seo resources/seo

HEALTHCHECK --interval=30s --timeout=3s CMD ["./healthcheck"]

CMD ["./app", "./migrate"]
