FROM alpine

EXPOSE 80

COPY bin/sadwave-events-api-v2 /bin

CMD ["/bin/sadwave-events-api-v2"]