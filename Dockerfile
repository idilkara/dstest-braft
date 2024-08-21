# syntax=docker/dockerfile:1

# We use a multi-stage build setup.
# (https://docs.docker.com/build/building/multi-stage/)

###############################################################################
# Stage 1 (to create a "build" image, ~850MB)                                 #
###############################################################################
# Image from https://hub.docker.com/_/golang
FROM golang:1.22.2 AS builder
# smoke test to verify if golang is available
RUN go version

ARG PROJECT_VERSION

COPY ./dstest /go/src/dstest/
WORKDIR /go/src/dstest/
RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Version=${PROJECT_VERSION}'" \
    -o cmd/dstest/main cmd/dstest/main.go
RUN go test -cover -v ./...

###############################################################################
# Stage 2 (braft)                                                             #
###############################################################################
FROM ubuntu:24.04 AS braft-builder

WORKDIR /src/braft_builder

#build brpc-install dependencies
#for out of sync host machines -> RUN apt-get -o Acquire::Check-Valid-Until=false update

RUN apt-get update && apt-get install -y \
    git \
    g++ \
    make \
    libssl-dev \
    libgflags-dev \
    libprotobuf-dev \
    libprotoc-dev \
    protobuf-compiler \
    libleveldb-dev \
    libsnappy-dev \
    libgoogle-perftools-dev \
    cmake \
    libgtest-dev \
    dos2unix
    
    
#build brpc
COPY ./brpc /src/braft_builder/brpc
WORKDIR /src/braft_builder/brpc/build
RUN cmake .. && make -j6 && make install


#WORKDIR /src/braft_builder/brpc/example/echo_c++
#RUN cmake -B build && cmake --build build -j4

#build braft
WORKDIR /src/braft_builder

COPY ./braft /src/braft_builder/braft
WORKDIR /src/braft_builder/braft/bld
RUN cmake -DBRPC_DIR=/src/braft_builder/brpc .. && make -j6

#build the counter example of braft 
#from braft-example - this is a special example choose another name to not confuse with braft examples 

COPY ./braft-example/counter_modified /src/braft_builder/braft/example/counter_modified
WORKDIR /src/braft_builder/braft/example/counter_modified
RUN cmake -DBRPC_DIR=/src/braft_builder . && make -j6

#need to do the following for some reason.
RUN dos2unix /src/braft_builder/braft/example/shflags

###############################################################################
# Stage 3 (to create a downsized "container executable", ~5MB)                #
###############################################################################
FROM ubuntu:24.04

RUN apt-get update
RUN apt-get install -y openjdk-17-jre-headless golang-go
RUN apt-get update && apt-get install -y \
    git \
    g++ \
    make \
    libssl-dev \
    libgflags-dev \
    libprotobuf-dev \
    libprotoc-dev \
    protobuf-compiler \
    libleveldb-dev \
    libsnappy-dev \
    libgoogle-perftools-dev \
    cmake \
    psmisc 

WORKDIR /root
COPY --from=builder /go/src/dstest /root/dstest
COPY --from=braft-builder /src/braft_builder/brpc /root/brpc
COPY --from=braft-builder /src/braft_builder/braft /root/braft

# root 
    # brpc 
    # braft 
        # example/counter_modified
    # dstest
#

WORKDIR /root/dstest/cmd/dstest/
ENTRYPOINT ["./main"]
CMD ["run", "-c", "/configs/braft2.yml"]
