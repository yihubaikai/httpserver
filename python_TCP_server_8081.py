#!usr/bin/env python
# -*- coding:utf-8 -*-
from  gevent import monkey;monkey.patch_all()
import gevent
from socket import *
import os
print('start running...%d'%(os.getpid()))
def talk(conn,addr):
    while True:
        data = conn.recv(1024)
        strs = "[" . str(os.getpid()) + "]" + str(addr[0])  + "." + str(addr[1])  + "." + str(data)
        conn.send(strs.upper())
    conn.close()
def server(ip,duankou):
    server = socket(AF_INET, SOCK_STREAM)
    server.setsockopt(SOL_SOCKET, SO_REUSEADDR, 1)
    server.bind((ip,duankou))
    server.listen(5)
    while True:
        conn,addr = server.accept()  #等待链接
        gevent.spawn(talk,conn,addr)  #异步执行 （p =Process(target=talk,args=(coon,addr))
                                                # p.start()）相当于开进程里的这两句
    server.close()
if __name__ == '__main__':
    server('127.0.0.1',8081)

#服务端利用协程