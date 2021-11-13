FROM alpine

EXPOSE 3000

COPY bin/sadwave-events-api-v2 /bin

CMD ["/bin/sadwave-events-api-v2"]