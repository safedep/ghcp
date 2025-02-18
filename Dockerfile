FROM golang:1.23-bookworm AS build

WORKDIR /build

COPY go.mod go.sum Makefile ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN make

FROM golang:1.23-bookworm

WORKDIR /app

COPY --from=build /build/bin/ghcp /app/ghcp

ENV PATH "${PATH}:/app/"
EXPOSE 8000

RUN useradd -ms /bin/bash nonroot

USER nonroot:nonroot

CMD ["ghcp", "server", "--address", "0.0.0.0:8000"]
