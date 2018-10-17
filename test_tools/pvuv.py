#!/usr/bin/env python3
#_*_ coding:utf8 _*_

from  urllib import request
from urllib import parse
import json


host = '192.168.8.211'
post_key = 'ActionData'

msg_pv = {
    "category": "LivPV",
    "appid": "nnnnnnnn",
    "key": "xxx-xxx-xxx",
    "globalid": "55310D34-C95C-76FA-4255-BB5A8AE726AA",
    "uid": "test_user",
    "time": 1532934840
}

msg_uv = {
    "category": "LiveUV",
    "appid": "nnnnnnnn",
    "key": "xxx-xxx-xxx",
    "globalid": "55310D34-C95C-76FA-4255-BB5A8AE726AA",
    "time": 1532934840
}


def postMsg(url, msg):
    data = {post_key: json.dumps(msg)}
    data_enc = parse.urlencode(data).encode()
    req = request.Request(url, data_enc)
    resp = request.urlopen(req).read().decode()
    return resp


def testPV():
    url = "http://{h}/counter/pv".format(h=host)
    return postMsg(url, msg_pv)


def testUV():
    url = "http://{h}/counter/uv".format(h=host)
    return postMsg(url, msg_uv)


if __name__ == "__main__":
    print(testPV())
    print(testUV())
