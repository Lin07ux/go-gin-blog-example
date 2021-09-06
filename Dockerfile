FROM scratch

ENV GOPATH="/home"
WORKDIR $GOPATH/src/github.com/lin07ux/go-gin-example
COPY ./conf ./conf
COPY ./go-gin-example .
EXPOSE 8000
CMD ["./go-gin-example"]