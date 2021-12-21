FROM gcr.io/distroless/base
ARG BIN
COPY /bin/cli /cli
CMD ["/cli"]
