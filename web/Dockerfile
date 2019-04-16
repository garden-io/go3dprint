ARG BASE_IMAGE_VERSION
FROM base-image:${BASE_IMAGE_VERSION}

COPY index.html ./
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["find . | entr -r ./hotreload.sh"]