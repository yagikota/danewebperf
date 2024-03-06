#!/bin/sh

sudo tcpdump -i any  -w /captured/firefox.pcap &

sudo python3 /home/seluser/measure/pageload_measure.py $@
