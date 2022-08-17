#!/bin/bash

if [ $IPTABLES_BACKEND = "legacy" ]
then
    update-alternatives --set iptables /usr/sbin/iptables-legacy
    update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy
    update-alternatives --set arptables /usr/sbin/arptables-legacy
    update-alternatives --set ebtables /usr/sbin/ebtables-legacy
elif [ $IPTABLES_BACKEND = "nft" ]
then
    update-alternatives --set iptables /usr/sbin/iptables-nft
    update-alternatives --set ip6tables /usr/sbin/ip6tables-nft
    update-alternatives --set arptables /usr/sbin/arptables-nft
    update-alternatives --set ebtables /usr/sbin/ebtables-nft
else
    echo "environment variable IPTABLES_BACKEND is unrecognized, use -h for more help"
    exit 1
fi

/usr/bin/trav $*
