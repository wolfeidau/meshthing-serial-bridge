# meshthing-serial-bridge [![Build Status](https://drone.io/github.com/wolfeidau/meshthing-serial-bridge/status.png)](https://drone.io/github.com/wolfeidau/meshthing-serial-bridge/latest)

This program reads packets from the serial port and writes them to a unix pipe.

# overview

This program expects packets to be encoded with a single byte length prefix, then the packet data.

This program creates a fifo which is located at `/tmp/wireshark` at the moment, this can be used by wireshark as a capture source.

# building

Ensure you have golang 1.2 installed and just run make.

```
make
```

Binary is moved to `bin/sniffer-bridge`

# usage 

```
Usage of sniffer-bridge:
  -port="": optional path to serial device
  -version=false: print the version information
```

Run `sniffer-bridge` in a terminal, if you omit the `-port` it will search for a usbserial device and use that (OSX only at the moment).

```
sniffer-bridge -port=/dev/tty.usbmodem1411
```

Then run wireshark in another terminal.

```
tshark -i /tmp/wireshark
```

# links

Meshthing sniffer code https://github.com/geekscape/contiki/tree/example-meshthing-sniffer/examples/meshthing/sniffer

# licence

Copyright (c) 2013 Mark Wolfe released under the MIT license.