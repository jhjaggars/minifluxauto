FROM alpine:latest
COPY build/minifluxauto /bin/minifluxauto
CMD ["/bin/minifluxauto"]