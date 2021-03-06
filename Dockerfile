FROM golang:1.17.0-alpine AS build
ADD . /src/
WORKDIR /src/cmd/server

RUN CGO_ENABLED=0 go build -o /out/mintsql

FROM scratch
COPY --from=build /out/mintsql /app/mintsql

EXPOSE 7384

ENTRYPOINT [ "/app/mintsql" ]
