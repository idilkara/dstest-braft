#!/bin/bash
BIN=../../../braft/example/counter_modified
cd ${BIN}


. ${BIN}/../shflags

# define command-line flags
DEFINE_boolean clean 1 'Remove old "runtime" dir before running'
DEFINE_integer add_percentage 100 'Percentage of fetch_add operation'
DEFINE_integer bthread_concurrency '8' 'Number of worker pthreads'
DEFINE_integer server_port 8100 "Port of the first server"
DEFINE_integer server_num '3' 'Number of servers'
DEFINE_integer thread_num 1 'Number of sending thread'
DEFINE_string crash_on_fatal 'true' 'Crash on fatal log'
DEFINE_string log_each_request 'true' 'Print log for each request'
DEFINE_string valgrind 'false' 'Run in valgrind'
DEFINE_string use_bthread "true" "Use bthread to send request"


# hostname prefers ipv6
IP=`hostname -i | awk '{print $NF}'`
#IP=127.0.0.1

FLAGS_server_num=3
FLAGS_port=8100



raft_peers="172.17.0.2:8100:0,172.17.0.2:8101:0,172.17.0.2:8102:0,"

#raft_peers="127.0.0.1:8100:0,127.0.0.1:8101:0,127.0.0.1:8102:0,"

echo ${raft_peers}



export TCMALLOC_SAMPLE_PARAMETER=524288

mkdir -p runtime_cli/2
cp ./counter_client runtime_cli/2 #it exists error
echo runtime_cli/2
cd runtime_cli/2

#run this counter client with these flags
./counter_client \
        --add_percentage=${FLAGS_add_percentage} \
        --bthread_concurrency=${FLAGS_bthread_concurrency} \
        --conf="${raft_peers}" \
        --crash_on_fatal_log=${FLAGS_crash_on_fatal} \
        --log_each_request=${FLAGS_log_each_request} \
        --thread_num=${FLAGS_thread_num} \
        --use_bthread=${FLAGS_use_bthread} 

#is it possible that since each use the same client server - it messes up  the output? 
