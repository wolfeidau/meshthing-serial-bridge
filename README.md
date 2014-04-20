# meshthing-serial-bridge

This program reads packets from the serial port and writes them to a unix pipe.

# overview

This program expects packets to be encoded with a single byte length prefix, then the packet data.

This program creates a fifo which is located at "/tmp/wireshark" by default, this can be used by wireshark as a capture source.

# usage 

```
Usage of ./sniffer-bridge:
  -port="": optional path to serial device
  -version=false: print the version information
```

Run `sniffer-bridge` in a terminal, if you omit the `-port` it will search for a usbserial device and use that (OSX only at the moment).

```
./sniffer-bridge -port=/dev/tty.usbmodem1411
```

```
tshark -i /tmp/wireshark
```

# licence

Copyright (c) 2013 Mark Wolfe released under the MIT license.