FROM golang:latest

# updates sources
RUN apt-get --quiet=2 --download-only update

# install libzmq
RUN apt-get install --yes libzmq3-dev pkg-config
ENV PKG_CONFIG_PATH /usr/lib/x86_64-linux-gnu/pkgconfig

# install node
RUN wget https://deb.nodesource.com/setup_6.x -O setup-node.sh \
 && chmod +x setup-node.sh \
 && ./setup-node.sh \
 && rm setup-node.sh \
 && apt-get install --yes nodejs
ENV PATH /go/app/node_modules/.bin:$PATH

ENTRYPOINT ["./test.sh"]
