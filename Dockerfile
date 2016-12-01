FROM daocloud.io/library/centos

MAINTAINER JREAM.LU <ms.08.lu@gmail.com>

RUN     mkdir -p /data/gowww/jkernel \
        && yum install -y git \
        && yum install -y wget \
        && mkdir -p /soft/ \
        && cd /soft \
        && wget http://www.golangtc.com/static/go/1.7.3/go1.7.3.linux-amd64.tar.gz \
        && tar -zxvf go1.7.3.linux-amd64.tar.gz \
        && mv go /usr/loacl/ \

ENV GOROOT /usr/local/go
ENV GOPATH /data/gowww/jkernel
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

WORKDIR /data/gowww/jkernel