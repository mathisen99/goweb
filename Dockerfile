FROM golang:1.19-alpine
LABEL website="Grit:lab forum"
#todo add info
LABEL desc=""
LABEL authors=""
RUN apk update && apk add bash \
    && apk --update-cache add sqlite \
    && rm -rf /var/cache/apk/*
WORKDIR /app
COPY . .
RUN go build -o /forum-bin
EXPOSE 8080
CMD ["/forum-bin"]