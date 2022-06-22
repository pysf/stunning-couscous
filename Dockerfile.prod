FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o /stunning-couscous

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /stunning-couscous /stunning-couscous

ENTRYPOINT ["/stunning-couscous"]