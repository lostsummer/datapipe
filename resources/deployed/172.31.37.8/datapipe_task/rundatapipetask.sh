#!/bin/bash

export RunEnv=production && cd /emoney/datapipe_task && nohup ./datapipe_task > /emoney/datapipe_task/logs/stdout.log &

