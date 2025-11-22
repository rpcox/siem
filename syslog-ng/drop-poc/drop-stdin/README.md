## drop-stdin
    
A proof of concept example of using a syslog-ng destination program.
    
### syslog-ng config
    
    filter f_siem-drop-1 { match('^[^ ]+ [^ ]+ drop', value("MSG")); };
    destination d_drop { program("/usr/local/sbin/drop-stdin" template("$MSG\n")); };
    log { source(s_syslog); filter(f_siem-drop-1); destination(d_drop); flags(final); };
    
    
Start syslog-ng and check the service status. Note where drop-stdin shows up in the cgoup (unfortunately, through a shell)
    
    > systemctl status syslog-ng
    ● syslog-ng.service - System Logger Daemon
         Loaded: loaded (/lib/systemd/system/syslog-ng.service; enabled; preset: enabled)
         Active: active (running) since Thu 2025-11-20 11:20:53 MST; 3s ago
           Docs: man:syslog-ng(8)
       Main PID: 2489889 (syslog-ng)
          Tasks: 19 (limit: 3921)
            CPU: 484ms
         CGroup: /system.slice/syslog-ng.service
                 ├─2489889 /usr/sbin/syslog-ng -F
                 ├─2489890 /bin/sh -c /usr/local/sbin/drop-stdin
                 └─2489891 /usr/local/sbin/drop-stdin
    
    Nov 20 11:20:53 loghost02 systemd[1]: Starting syslog-ng.service - System Logger Daemon...
    Nov 20 11:20:53 loghost02 systemd[1]: Started syslog-ng.service - System Logger Daemon.

    
On a remote host execute
    
    > for i in {1..4}; do echo "<123>2025-11-19T12:34:56Z spud drop: DROP THIS" | nc loghost02.lan 514; done
    > echo "<123>2025-11-19T12:34:56Z spud nodrop: DON'T DROP THIS" | nc loghost02.lan 514

    
Check the logs
    
    > tail -f /mnt/log/network/other.log | grep spud
    <123>2025-11-19T12:34:56Z spud nodrop: DON'T DROP THIS
    
    > tail -f /var/log/drop.log 
      1: <123>2025-11-19T12:34:56Z spud drop: DROP THIS
      2: <123>2025-11-19T12:34:56Z spud drop: DROP THIS
      3: <123>2025-11-19T12:34:56Z spud drop: DROP THIS
      4: <123>2025-11-19T12:34:56Z spud drop: DROP THIS
    


