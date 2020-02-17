FROM golang:1.13.6
WORKDIR /service/
COPY . /service/
RUN go get github.com/go-chi/chi
EXPOSE 8081
ENTRYPOINT ["go", "run", "infura-mock.go"]