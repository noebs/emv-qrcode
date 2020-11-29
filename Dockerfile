FROM golang:latest


COPY . /app
WORKDIR /app
RUN go build

CMD ["emv-qrcode"]
EXPOSE 8012