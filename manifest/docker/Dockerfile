#FROM alpinelinux/golang:3.18
FROM loads/alpine:3.8

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

###############################################################################
#                                INSTALLATION
###############################################################################

ENV WORKDIR                 /app
ADD resource                $WORKDIR/
ADD ./temp/linux_amd64/main $WORKDIR/main
ADD ./version.txt $WORKDIR/version.txt
#ADD ./i18n  $WORKDIR/i18n
RUN chmod +x $WORKDIR/main

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main
