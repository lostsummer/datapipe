#!/usr/bin/env python
# -*- coding: utf-8 -*-

rds_nanhui = {
    "host" : '172.28.1.118',
    "port" : 6379,
    "db"   : 0
}

rds_wanguo = {
    "host" : '172.31.23.16',
    "port" : 6395,
    "db"   : 0
}

redises = {
    "readfrom" : rds_nanhui,
    "writeto"  : rds_wanguo
}
