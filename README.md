Netcat for onion router
=======================

Netcat-like tool for tor network

Features
--------
* Autodetect SOCKS port by control socket in client mode
* Creates ephemeral hidden service in listen mode

Requirements
------------
* Tor daemon with control socket enabled (tcp/unix)

Usage
-----

```
# Create ephemeral hidden service and wait for incoming connections
$ torcat -v -l 9999
(stderr) iwwcfmu5ncqwnncw.onion
(stderr) [Waiting]
```

```
# Connect to hidden service, send "Hello world!" and close connection
$ torcat -v iwwcfmu5ncqwnncw.onion 9999
(stderr) [Connected]
(stdin)  Hello world!
^C
```
