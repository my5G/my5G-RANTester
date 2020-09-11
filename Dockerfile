FROM golang:1.14.4-stretch

COPY . /my5G-RANTester

WORKDIR /my5G-RANTester

RUN go mod download

RUN cd src && go build -o my5g-rantester

ENTRYPOINT ["/my5G-RANTester/src/my5g-rantester"]



