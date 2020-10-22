FROM golang:1.12.1-alpine3.9
RUN mkdir /app
#RUN mkdir /app/cli
ADD . /app
#ADD ./cli /app/cli
WORKDIR /app
RUN go build -o cliapp ./cli
RUN go build -o main .