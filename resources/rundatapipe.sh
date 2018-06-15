#!/bin/bash

export RunEnv=production && cd /emoney/datapipe && nohup ./datapipe > /emoney/datapipe/logs/stdout.log &
