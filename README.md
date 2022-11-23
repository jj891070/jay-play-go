$ vim /etc/logrotate.d/test
```txt
/etc/redis/*.log {
    daily
    size 0M
    dateext
    missingok
    rotate 14
    nocompress
    notifempty
    noolddir
    copytruncate
    #delaycompress
    #create 0777 root root
    #size 1M
}
```

$ service logrotate status

$ sudo vim /lib/systemd/system/logrotate.service
```txt
[Unit]
Description=Rotate log files
Documentation=man:logrotate(8) man:logrotate.conf(5)
ConditionACPower=true

[Service]
Type=oneshot
ExecStart=/usr/sbin/logrotate /etc/logrotate.conf

# performance options
Nice=19
IOSchedulingClass=best-effort
IOSchedulingPriority=7

# hardening options
#  details: https://www.freedesktop.org/software/systemd/man/systemd.exec.html
#  no ProtectHome for userdir logs
#  no PrivateNetwork for mail deliviery
#  no ProtectKernelTunables for working SELinux with systemd older than 235
#  no MemoryDenyWriteExecute for gzip on i686
PrivateDevices=true
PrivateTmp=true
ProtectControlGroups=true
ProtectKernelModules=true
ProtectSystem=full
RestrictRealtime=true
ReadWritePaths=/etc/redis
```

$ sudo systemctl daemon-reload

$ sudo systemctl restart logrotate