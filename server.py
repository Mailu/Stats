from scapy.all import *

import uuid
import re
import time


class Server(object):

    def __init__(self, file_path, domain):
        self.logfile = open(file_path, "a+")
        self.domain = domain

    def handle(self, packet):
        """ Handle a DNS query, then log the result.
        """
        if (DNS in packet and packet[DNS].opcode == 0L and packet[DNS].ancount == 0):
            name = packet[DNS].qd.qname.lower()
            # Check the name format
            if not name.endswith(self.domain):
                return
            split_name = name.split(".")
            if len(split_name) < 2:
                return
            version, instanceid = split_name[:2]
            try:
                uuid.UUID(instanceid, version=4)
            except ValueError:
                return
            if len(version) > 2 or not version.isdigit():
                return
            # Actually log an entry
            self.log(version, instanceid)

    def log(self, version, instanceid):
        """ Log an instance pinging the stats server.
        """
        timestamp = str(time.time())
        self.logfile.write("{},{},{}\n".format(version, instanceid, timestamp))
        self.logfile.flush()


    def serve(self):
        """ Start the server loop.
        """
        sniff(filter="udp port 53", store=0, prn=self.handle)


if __name__ == "__main__":
    server = Server("output.log", "stats.mailu.io")
    server.serve()
