###############################################################################
#                                BUILD
###############################################################################
FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
#RUN go get
RUN GOOS=linux GOARCH=amd64 go build -o temp/linux_amd64/main main.go
RUN #rm -rf version.txt && echo "buildtime:"$(date +%Y%m%d%H%M)"" > version.txt


###############################################################################
#                                PREPARE
###############################################################################
FROM loads/alpine:3.8
RUN apk add --update ttf-dejavu fontconfig && rm -rf /var/cache/apk/*
COPY ./manifest/fonts/*.ttc /usr/share/fonts/
COPY ./manifest/fonts/*.ttf /usr/share/fonts/
RUN apk add --update curl && rm -rf /var/cache/apk/*
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

###############################################################################
#                                INSTALLATION
###############################################################################
ENV WORKDIR                 /app
ADD resource                $WORKDIR/
COPY --from=builder /app/temp/linux_amd64/main $WORKDIR/main
ADD ./version.txt $WORKDIR/version.txt
ADD ./resource/email/template.html $WORKDIR/resource/email/template.html
ADD ./manifest/i18n  $WORKDIR/i18n
RUN chmod +x $WORKDIR/main
###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main
