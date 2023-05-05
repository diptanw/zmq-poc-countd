FROM golang:1.7

RUN set -x; \
  export DEBIAN_FRONTEND=noninteractive; \
  echo "deb http://ppa.launchpad.net/bpaquet/zeromq4-precise/ubuntu precise main " >> /etc/apt/sources.list \
  && apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys FF160426C4832F92 \
  && apt-get update \
  && apt-get install -y \
    libzmq1=4.0.4-bpa~precise1 \
    libzmq-dev=4.0.4-bpa~precise1 \
    pkg-config

ENV PATH /go/bin:$PATH
ENV GOPROJECTPATH /go/src/emitter
COPY . $GOPROJECTPATH
WORKDIR $GOPROJECTPATH

RUN go get -v

CMD ["bash"]
