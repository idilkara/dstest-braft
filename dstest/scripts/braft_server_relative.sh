#!/bin/bash
BIN=../../../braft/example/counter_modified
SHFLAGSDIR=../../../braft/example

cd ${BIN}

#echo "Listing contents of ${SHFLAGSDIR}:"
#ls ${SHFLAGSDIR}

#.  ${SHFLAGSDIR}/shflags
#echo "Sourced shflags successfully"

. ${SHFLAGSDIR}/shflags


# define command-line flags
DEFINE_string crash_on_fatal 'true' 'Crash on fatal log'
DEFINE_integer bthread_concurrency '18' 'Number of worker pthreads'
DEFINE_string sync 'true' 'fsync each time'
DEFINE_string valgrind 'false' 'Run in valgrind'
DEFINE_integer max_segment_size '8388608' 'Max segment size'
DEFINE_integer server_num '3' 'Number of servers'
DEFINE_boolean clean 1 'Remove old "runtime" dir before running'
DEFINE_integer port $5 "Port of the first server"


# The alias for printing to stderr
#alias error=" echo counter: "

# hostname prefers ipv6
IP=`hostname -i | awk '{print $NF}'`
#IP=127.0.0.1

#SERVER NUM AND FIRST PORT
FLAGS_server_num=3

raft_peers="172.17.0.2:$1:0,172.17.0.2:$2:0,172.17.0.2:$3:0,"


#raft_peers="172.17.0.2:8100:0,172.17.0.2:8101:0,172.17.0.2:8102:0,"

#raft_peers="127.0.0.1:$1:0,127.0.0.1:$2:0,127.0.0.1:$3:0,"
export TCMALLOC_SAMPLE_PARAMETER=524288


echo ${raft_peers}

echo "IN SERVER "

mkdir -p runtime/$4
cp ./counter_server runtime/$4 #it exists error
echo runtime/$4
cd runtime/$4

port_num=$((${FLAGS_port}))
echo $port_num
echo 0

ls -l

./counter_server \
        -bthread_concurrency=${FLAGS_bthread_concurrency}\
        -crash_on_fatal_log=${FLAGS_crash_on_fatal} \
        -raft_max_segment_size=${FLAGS_max_segment_size} \
        -raft_sync=${FLAGS_sync} \
        -port=${port_num}  -conf="${raft_peers}" &
cd ../..


#keep the script running
while true; do
    sleep 1  
done