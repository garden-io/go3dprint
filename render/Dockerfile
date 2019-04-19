ARG BASE_IMAGE_VERSION
FROM base-image:${BASE_IMAGE_VERSION}

COPY testfile.test ./
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["./app"]