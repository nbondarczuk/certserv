FROM alpine

LABEL name="certsrv"

ENV PORT=8080
ENV WORK=/work
ENV PATH=$WORK:$PATH

RUN mkdir -p $WORK

COPY certsrv $WORK

EXPOSE $PORT

CMD ["/work/certsrv"]


