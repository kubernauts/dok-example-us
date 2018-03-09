FROM golang:latest 
MAINTAINER Michael Hausenblas <mhausenb@redhat.com>
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
RUN go build -o main .
USER nobody:nobody
CMD ["/app/main"]