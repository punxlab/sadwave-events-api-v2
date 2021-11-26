FROM alpine

EXPOSE 80

ADD . /cache

COPY bin/sadwave-events-api-v2 /bin

CMD ["/bin/sadwave-events-api-v2"]