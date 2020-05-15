# # # # # Build go backend # # # # #
FROM golang:alpine as gobuilder

# Install tools needed to get go projects and build
RUN apk --no-cache add build-base git bzr mercurial gcc
RUN go get nhooyr.io/websocket
RUN mkdir -p /go/src/github.com/monodop/devlog
ADD . /go/src/github.com/monodop/devlog
WORKDIR /go/src/github.com/monodop/devlog
RUN go build -o devlog_server .

# # # # # Build frontend # # # # #
FROM node as nodebuilder
# Install OpenJDK-8
RUN apt-get update && \
    apt-get install -y openjdk-8-jdk && \
    apt-get install -y ant && \
    apt-get clean;
# Fix certificate issues
RUN apt-get update && \
    apt-get install ca-certificates-java && \
    apt-get clean && \
    update-ca-certificates -f;
# Setup JAVA_HOME -- useful for docker commandline
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/
RUN export JAVA_HOME

RUN mkdir /build
ADD ./frontend /build
WORKDIR /build
RUN npm install && npm run antlr && npm run build

# Build final docker image
FROM alpine
RUN addgroup -S appuser && adduser -S -D -H -h /app appuser
USER appuser
COPY --from=gobuilder /go/src/github.com/monodop/devlog/devlog_server /app/
COPY --from=nodebuilder --chown=appuser:appuser /build/dist /app/frontend
WORKDIR /app
EXPOSE 9090:9090
EXPOSE 9091:9091
CMD ["./devlog_server"]