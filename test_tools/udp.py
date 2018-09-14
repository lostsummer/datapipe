import socket

def udpsend(msg, ipaddr):
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.sendto(msg.encode('utf-8'), ipaddr)
    s.close()

if __name__ == '__main__':
    i = 0
    while i<10:
        udpsend("[10000]2018-09-13 12:25:80 我阿斯蒂芬", ('192.168.42.86', 4444))
        udpsend("{\"Appid\":\"1000000\", \"Level\":\"debug\", \"Message\":\"2018-09-13 12:25:80 我阿斯蒂芬\"}", ('192.168.42.86', 4445))
        i += 1


