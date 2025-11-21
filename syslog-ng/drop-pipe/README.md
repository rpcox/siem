## drop-pipe

A proof of concept example of using a named pipe w/ syslog-ng destination program.

### syslog-ng config


    filter f_siem-drop-1 { match('^[^ ]+ [^ ]+ drop', value("MSG")); };
    destination d_drop { 
    	pipe("/tmp/mypipe"
    		hook-commands(
    			startup("systemctl start drop-pipe")
    			shutdown("systemctl stop drop-pipe")
    		)
    		template("$MSG\n")
    	); 
    };
    log { source(s_syslog); filter(f_siem-drop-1); destination(d_drop); flags(final); };
        

### test

        
Start syslog-ng and check the service status of syslog-ng and drop-pipe. Note that drop-pipe.service remains disabled.
    
    > systemctl status syslog-ng
    ● syslog-ng.service - System Logger Daemon
         Loaded: loaded (/lib/systemd/system/syslog-ng.service; enabled; preset: enabled)
         Active: active (running) since Fri 2025-11-21 06:25:34 MST; 17min ago
           Docs: man:syslog-ng(8)
       Main PID: 2551169 (syslog-ng)
          Tasks: 11 (limit: 3921)
            CPU: 1.515s
         CGroup: /system.slice/syslog-ng.service
                 └─2551169 /usr/sbin/syslog-ng -F
    
    Nov 21 06:25:34 loghost02 systemd[1]: Starting syslog-ng.service - System Logger Daemon...
    Nov 21 06:25:34 loghost02 systemd[1]: Started syslog-ng.service - System Logger Daemon.
    
    > systemctl status drop-pipe
    ● drop-pipe.service - drop-pipe poc
         Loaded: loaded (/etc/systemd/system/drop-pipe.service; disabled; preset: enabled)
         Active: active (running) since Fri 2025-11-21 06:25:34 MST; 17min ago
       Main PID: 2551172 (drop-pipe)
          Tasks: 7 (limit: 3921)
            CPU: 8ms
         CGroup: /system.slice/drop-pipe.service
                 └─2551172 /usr/local/sbin/drop-pipe
    
    Nov 21 06:25:34 loghost02 systemd[1]: Started drop-pipe.service - drop-pipe poc.
        
        
On a remote host execute
        

        > for i in {1..6}; do echo "<123>2025-11-21T12:34:56Z spud drop: DROP THIS" | nc loghost02.lan 514; done
        > echo "<123>2025-11-21:34:56Z spud nodrop: DON'T DROP THIS" | nc loghost02.lan 514
    
        
Check the logs on relay
        

    > tail -f /var/log/drop.log 
      1: <123>2025-11-20T12:34:56Z spud drop: DROP THIS
      2: <123>2025-11-20T12:34:56Z spud drop: DROP THIS
      3: <123>2025-11-20T12:34:56Z spud drop: DROP THIS
      4: <123>2025-11-20T12:34:56Z spud drop: DROP THIS
      5: <123>2025-11-20T12:34:56Z spud drop: DROP THIS
      6: <123>2025-11-20T12:34:56Z spud drop: DROP THIS
    
    > tail -f /mnt/log/network/other.log
    <123>2025-11-21T12:34:56Z spud nodrop: DON'T DROP THIS
    

### bug

When trying to stop syslog-ng (systemctl stop syslog-ng) it hung. I had to find drop-pipe in the process table and send a SIGKILL, then syslog-ng promptly stopped.

Using systemctl start/stop drop-pipe works fine. Issue is related to syslog-ng. May need to change out the shutdown command.
