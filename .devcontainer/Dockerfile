FROM adrienaury/go-devcontainer:v0.5

USER root

RUN apk add --update --progress --no-cache make gomplate

ARG VERSION_GOLICENSE=0.2.0
ARG VERSION_MILLER=5.10.2
RUN    wget -nv -O- https://github.com/mitchellh/golicense/releases/download/v${VERSION_GOLICENSE}/golicense_${VERSION_GOLICENSE}_linux_x86_64.tar.gz | tar xz -C /usr/bin golicense \
    && wget -nv -O/usr/bin/mlr https://github.com/johnkerl/miller/releases/download/v${VERSION_MILLER}/mlr.linux.x86_64 \
    && chmod +x /usr/bin/golicense /usr/bin/mlr

USER vscode
