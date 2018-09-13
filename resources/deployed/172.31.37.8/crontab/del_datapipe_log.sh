#!/bin/bash

del7daymore() {
    dir=$1
    find $dir -mtime +7 -name "*.log" -exec rm -rf {} \;
}

process1srvlog() {
    dir=$1
    del7daymore $1/logs/
    del7daymore $1/innerlogs/
}

process1srvlog /emoney/datapipe
process1srvlog /emoney/datapipe_task

