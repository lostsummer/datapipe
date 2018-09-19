#!/bin/bash

DIR=$(dirname $(realpath $0))
EXE=$(basename $DIR)    # 目录名同程序名相同
export RunEnv=production && cd $DIR && nohup ./$EXE > $DIR/logs/stdout.log &

