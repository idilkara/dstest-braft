#!/bin/bash

BIN=../../../braft/example/counter_modified
cd ${BIN}

#killall -9 counter_client
killall -9 counter_server


rm -rf runtime