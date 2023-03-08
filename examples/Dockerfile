FROM ghcr.io/foundry-rs/foundry

# workdir
RUN mkdir -p /examples
WORKDIR /examples

# deps
RUN apk add curl bash

COPY . ./
RUN for init in $(ls */contracts/init.sh); do (cd $(dirname $init) && forge build); done

ENTRYPOINT ["/examples/entrypoint.sh"]
CMD []
