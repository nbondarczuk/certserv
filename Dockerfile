FROM alpine

LABEL name="certserv"

ENV PORT=8080
ENV WORK=/work
ENV PATH=$WORK:$PATH

RUN mkdir -p $WORK

COPY certserv $WORK

EXPOSE $PORT

CMD ["/work/certserv"]


