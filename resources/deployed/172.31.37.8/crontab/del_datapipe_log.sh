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

process1srvlog /emoney/dp_api2/datapipe
process1srvlog /emoney/dp_api2/datapipe_task
process1srvlog /emoney/dp_aliapi/datapipe2
process1srvlog /emoney/dp_aliapi/datapipe_task2

