#!/bin/bash

# Enable or disable ipv6. (Default: "yes", Possible Values: "yes, no")
DO_IPV6=${DO_IPV6:-yes}

# Enable or disable ipv4. (Default: "yes", Possible Values: "yes, no")
DO_IPV4=${DO_IPV4:-yes}

# Enable or disable udp. (Default: "yes", Possible Values: "yes, no")
DO_UDP=${DO_UDP:-yes}

# Enable or disable tcp. (Default: "yes", Possible Values: "yes, no")
DO_TCP=${DO_TCP:-yes}

# Verbosity number, 0 is least verbose. (Default: "0", Possible Values: "<integer>")
VERBOSITY=${VERBOSITY:-0}

# Number of threads to create. 1 disables threading. (Default: "1", Possible Values: "<integer>")
NUM_THREADS=${NUM_THREADS:-1}

# Buffer size for UDP port 53 incoming. Use 4m to catch query spikes for busy servers. (Default: "0", Possible Values: "<integer>")
SO_RCVBUFF=${SO_RCVBUFF:-0}

# Buffer size for UDP port 53 outgoing. Use 4m to handle spikes on very busy servers. (Default: "0", Possible Values: "<integer>")
SO_SNDBUF=${SO_SNDBUF:-0}

# Use SO_REUSEPORT to distribute queries over threads. (Default: "no", Possible Values: "yes, no")
SO_REUSEPORT=${SO_REUSEPORT:-no}

# EDNS reassembly buffer to advertise to UDP peers. 1480 can solve fragmentation (timeouts). (Default: "4096", Possible Values: "<integer>")
EDNS_BUFFER_SIZE=${EDNS_BUFFER_SIZE:-4096}

# The amount of memory to use for the message cache. Plain value in bytes or you can append k, m or G. (Default: "4m", Possible Values: "<integer>")
MSG_CACHE_SIZE=${MSG_CACHE_SIZE:-4m}

# The amount of memory to use for the RRset cache. Plain value in bytes or you can append k, m or G. (Default: "4m", Possible Values: "<integer>")
RRSET_CACHE_SIZE=${RRSET_CACHE_SIZE:-4m}

# The time to live (TTL) value lower bound, in seconds. If more than an hour could easily give trouble due to stale data. (Default: "0", Possible Values: "<integer>")
CACHE_MIN_TTL=${CACHE_MIN_TTL:-86400}

# The time to live (TTL) value cap for RRsets and messages in the cache. Items are not cached for longer. In seconds. (Default: "86400", Possible Values: "<integer>")
CACHE_MAX_TTL=${CACHE_MAX_TTL:-86400}

# The time to live (TTL) value cap for negative responses in the cache. (Default: "3600", Possible Values: "<integer>")
CACHE_MAX_NEGATIVE_TTL=${CACHE_MAX_NEGATIVE_TTL:-3600}

# Enable to automatically re-fetch cached records before they expire. (Default: "no", Possible Values: "yes, no")
PREFETCH=${PREFTECH:-no}

# Enable to not answer id.server and hostname.bind queries. (Default: "no", Possible Values: "yes, no")
HIDE_IDENTITY=${HIDE_IDENTITY:-no}

# Enable to not answer version.server and version.bind queries. (Default: "no", Possible Values: "yes, no")
HIDE_VERSION=${HIDE_VERSION:-no}

# print statistics to the log (for every thread) every N seconds. (Default: "0", Possible Values: "0, 1")
STATISTICS_INTERVAL=${STATISTICS_INTERVAL:-0}

# enable cumulative statistics, without clearing them after printing. (Default: "no", Possible Values: "yes, no")
STATISTICS_CUMULATIVE=${STATISTICS_CUMULATIVE:-no}

# enable extended statistics (query types, answer codes, status) printed from unbound-control. (Default: "no", Possible Values: "yes, no")
EXTENDED_STATISTICS=${EXTENDED_STATISTICS:-yes}

# Sets the interface to listen on useful when using --net=host (Default 0.0.0.0, Possible Values: "", "@")
INTERFACE=${INTERFACE:-0.0.0.0}

# Enable the remote control feature (Default "yes", Possible Values: "yes, no")
REMOTE_CONTROL_ENABLE=${REMOTE_CONTROL_ENABLE:-yes}


sed 's/{{DO_IPV6}}/'"${DO_IPV6}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{DO_IPV4}}/'"${DO_IPV4}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{DO_UDP}}/'"${DO_UDP}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{DO_TCP}}/'"${DO_TCP}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{VERBOSITY}}/'"${VERBOSITY}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{NUM_THREADS}}/'"${NUM_THREADS}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{SO_RCVBUFF}}/'"${SO_RCVBUFF}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{SO_SNDBUF}}/'"${SO_SNDBUF}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{SO_REUSEPORT}}/'"${SO_REUSEPORT}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{EDNS_BUFFER_SIZE}}/'"${EDNS_BUFFER_SIZE}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{MSG_CACHE_SIZE}}/'"${MSG_CACHE_SIZE}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{RRSET_CACHE_SIZE}}/'"${RRSET_CACHE_SIZE}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{CACHE_MIN_TTL}}/'"${CACHE_MIN_TTL}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{CACHE_MAX_TTL}}/'"${CACHE_MAX_TTL}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{CACHE_MAX_NEGATIVE_TTL}}/'"${CACHE_MAX_NEGATIVE_TTL}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{PREFETCH}}/'"${PREFETCH}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{HIDE_IDENTITY}}/'"${HIDE_IDENTITY}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{HIDE_VERSION}}/'"${HIDE_VERSION}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{STATISTICS_INTERVAL}}/'"${STATISTICS_INTERVAL}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{STATISTICS_CUMULATIVE}}/'"${STATISTICS_CUMULATIVE}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{EXTENDED_STATISTICS}}/'"${EXTENDED_STATISTICS}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{INTERFACE}}/'"${INTERFACE}"'/' -i /usr/local/etc/unbound/unbound.conf
sed 's/{{REMOTE_CONTROL_ENABLE}}/'"${REMOTE_CONTROL_ENABLE}"'/' -i /usr/local/etc/unbound/unbound.conf

echo "Starting unbound..."
/usr/local/sbin/unbound -c /usr/local/etc/unbound/unbound.conf -d -v