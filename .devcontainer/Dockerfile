FROM cgifr/lino:latest-debian AS lino

FROM adrienaury/go-devcontainer:v3.1

USER root


ADD cgi_ca_root.crt /usr/local/share/ca-certificates/cgi_ca_root.crt
RUN update-ca-certificates

RUN apk add --update --progress --no-cache make gomplate jq

ARG VERSION_GOLICENSE=0.2.0
ARG VERSION_MILLER=6.2.0

RUN    wget -nv -O- https://github.com/mitchellh/golicense/releases/download/v${VERSION_GOLICENSE}/golicense_${VERSION_GOLICENSE}_linux_x86_64.tar.gz | tar xz -C /usr/bin golicense \
    && wget -nv -O- https://github.com/johnkerl/miller/releases/download/v${VERSION_MILLER}/miller-${VERSION_MILLER}-linux-amd64.tar.gz | tar xz --strip-components 1 -C /usr/bin miller-${VERSION_MILLER}-linux-amd64/mlr \
    && chmod +x /usr/bin/golicense /usr/bin/mlr

COPY --from=lino /usr/bin/pimo /usr/bin/pimo

RUN     wget https://github.com/sgreben/jp/releases/download/1.1.12/jp_1.1.12_linux_x86_64.zip -nv -O/tmp/jp.zip \
    && unzip /tmp/jp.zip  -d /usr/bin/ \
    && rm /tmp/jp.zip
USER vscode
