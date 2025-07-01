## keepalived

* Install keepalived
* keepalived.conf files in loghost0[12] directories are respective to the relay (/etc/keepalived/)
* Copy kad-* to /usr/local/bin


### normal operation

    loghost01:/etc/syslog-ng > kad-status
    keepalived: PIDs => 42906,42908

                 VIP    STATUS     STATE     PRI            INSTANCE     IFACE  ROUTER ID                 TRACKING
      ------------------------------------------------------------------------------------------------------------
     192.168.245.100    ACTIVE    MASTER     240    loghost01-MASTER      eth0        100   syslog-ng-track-MASTER
     192.168.245.101              BACKUP     230    loghost02-BACKUP      eth0        101   syslog-ng-track-BACKUP

    loghost01:/etc/syslog-ng > ip a | grep eth0
    2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
        inet 192.168.245.52/24 brd 192.168.245.255 scope global dynamic noprefixroute eth0
        inet 192.168.245.100/32 scope global eth0


    loghost02:/etc/keepalived > kad-status
    keepalived: PIDs => 18873,18874

                 VIP    STATUS     STATE     PRI            INSTANCE     IFACE  ROUTER ID                 TRACKING
      ------------------------------------------------------------------------------------------------------------
     192.168.245.101    ACTIVE    MASTER     240    loghost02-MASTER      eth0        101   syslog-ng-track-MASTER
     192.168.245.100              BACKUP     230    loghost01-BACKUP      eth0        100   syslog-ng-track-BACKUP

    loghost02:/etc/keepalived > ip a | grep eth0
    2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
        inet 192.168.245.53/24 brd 192.168.245.255 scope global dynamic noprefixroute eth0
        inet 192.168.245.101/32 scope global eth0


## syslog-ng down on loghost02

    loghost01:/etc/syslog-ng > kad-status
    keepalived: PIDs => 42906,42908

                 VIP    STATUS     STATE     PRI            INSTANCE     IFACE  ROUTER ID                 TRACKING
      ------------------------------------------------------------------------------------------------------------
     192.168.245.100    ACTIVE    MASTER     240    loghost01-MASTER      eth0        100   syslog-ng-track-MASTER
     192.168.245.101    ACTIVE    BACKUP     230    loghost02-BACKUP      eth0        101   syslog-ng-track-BACKUP

     loghost01:/etc/syslog-ng > ip a | grep eth0
    2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
        inet 192.168.245.52/24 brd 192.168.245.255 scope global dynamic noprefixroute eth0
        inet 192.168.245.100/32 scope global eth0
        inet 192.168.245.101/32 scope global eth0

    loghost02:/etc/keepalived > kad-status
    keepalived: PIDs => 18873,18874

                 VIP    STATUS     STATE     PRI            INSTANCE     IFACE  ROUTER ID                 TRACKING
      ------------------------------------------------------------------------------------------------------------
     192.168.245.101              MASTER     240    loghost02-MASTER      eth0        101   syslog-ng-track-MASTER
     192.168.245.100              BACKUP     230    loghost01-BACKUP      eth0        100   syslog-ng-track-BACKUP

    loghost02:/etc/keepalived > ip a | grep eth0
    2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
        inet 192.168.245.53/24 brd 192.168.245.255 scope global dynamic noprefixroute eth0


