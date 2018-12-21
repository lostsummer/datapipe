#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
data_down

Usage:
  data_down --help | --version
  data_down [(-a|-d <date>)] [-p]

Options:
  -h --help         show help
  -v --version      show version
  -a                downland all couner data
  -d <date>         specify date
  -p                print data
"""

import redis
import pickle
import config
import time
from docopt import docopt


__author__ = 'wangyx'
__version__ = '0.1.0'

rds_host = config.redises["readfrom"]["host"]
rds_port = config.redises["readfrom"]["port"]
rds_db = config.redises["readfrom"]["db"]

today = time.strftime('%Y%m%d')

if __name__ == '__main__':
    arguments = docopt(__doc__, version="data_down {0}".format(__version__))
    date = arguments["-d"]
    if date == None:
        date = today
    fetchall = arguments["-a"]
    if fetchall:
        key_pattern = 'EMoney.DataPipe:Counter:*'
    else:
        key_pattern = 'EMoney.DataPipe:Counter:*:{}'.format(date)

    shouldprint = arguments["-p"]

    cli = redis.Redis(host=rds_host,
                port=rds_port,
                db=rds_db)

    db_data = {}
    for key in cli.scan_iter(match=key_pattern):
        hash_data = {}
        for field, value in cli.hscan_iter(key):
            hash_data[field] = value

        db_data[key] = hash_data

    storage = open('rds.pkl', 'wb')
    pickle.dump(db_data, storage)
    storage.close()
    if shouldprint:
        for k in db_data:
            fv = db_data[k]
            for f in fv:
                v = fv[f]
                print((k, f, v))

