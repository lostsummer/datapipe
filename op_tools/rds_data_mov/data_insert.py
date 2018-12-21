#!/usr/bin/env python
# -*- coding: utf-8 -*-

import redis
import pickle
import config


__author__ = 'wangyx'
__version__ = '0.1.0'

rds_host = config.redises["writeto"]["host"]
rds_port = config.redises["writeto"]["port"]
rds_db = config.redises["writeto"]["db"]


if __name__ == "__main__":
    storage = open('rds.pkl', 'rb')
    db_data = pickle.load(storage)
    cli = redis.Redis(host=rds_host,
                    port=rds_port,
                    db=rds_db)
    pipe = cli.pipeline()
    for key in db_data:
        hash_data = db_data[key]
        for field in hash_data:
            value = hash_data[field]
            pipe.hset(key, field, value)
        pipe.execute()

