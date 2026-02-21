# E90-DTU(900SL22)E Test Tool
 This is a test tool to help with debugging and communication with E90-DTU devices.

 ## Notes:
 * Frequency id 850.125 to 930.125 in 1Mhz increments, ie. channels 1-80
 * Module address id for connecting to other E90 modules, like a filter 0-65535 where 65535 is broadcast/monitor all (filterless)
 * PP25-PP29 of LoRa Alliance regional parameters (LARP) frequency guide AU915-928
 * info https://meshtastic.au/wp/?page_id=47

## Goals
 I have an E90-DTU(900SL22)E fro Ebyte with the attached user manual.
The E90-DTU(900SL22)E does not have a serial interface nor dip switches, only an Ethernet port.
On the E90:
* set up a UDP server on socket A configured for port 8886.
* In the wireless settings:
    * Module Address: 65535 (broadcast/monitor)
    * Channel 76 (926.125Mhz)
    * Net ID 0
    * Air baud 2.4Kbps
    * Packet length 240 Byte
    * Packet RSSI enabled
    * Channel RSSI enabled
    * All other settings are default

## Trying to set up as meshtastic node
https://github.com/meshtastic
https://meshtastic.org/

## Device details
E90-DTU(900SL22)E
https://www.cdebyte.com/products/E90-DTU(900SL22)E
https://www.cdebyte.com/Resources-Download

## Tests:
```
90-DTU Meshtastic Control Tool
Device:  192.168.68.113:8886
================================
✓ UDP connection established

=== PROTOCOL ANALYSIS ===

Test 1: RSSI Query
Sending RSSI command: c0c1c2c30001
RSSI command sent successfully
Received RSSI response: c10001a0
Response bytes: [193 0 1 160]
Interpreted: RSSI: 160 dBm

Test 2: Status Query
Sending Status command: c0c1c2c30002
Status command sent successfully
Received Status response: c10002a000
Response bytes: [193 0 2 160 0]
Interpreted: Status: 0

Test 3: Forward/Enable Command
Sending Forward command: c0c1c2c30101
Forward command sent successfully
Received Forward response: c10101a0
Response bytes: [193 1 1 160]
Interpreted: Forward mode enabled (data: 160)
```

## Frequencies:
```
AU915-928

    Uplink:
    916.8 - SF7BW125 to SF12BW125
    917.0 - SF7BW125 to SF12BW125
    917.2 - SF7BW125 to SF12BW125
    917.4 - SF7BW125 to SF12BW125
    917.6 - SF7BW125 to SF12BW125
    917.8 - SF7BW125 to SF12BW125
    918.0 - SF7BW125 to SF12BW125
    918.2 - SF7BW125 to SF12BW125
    917.5 - SF8BW500

    Downlink:
    923.3 - SF7BW500 to SF12BW500 (RX1)
    923.9 - SF7BW500 to SF12BW500 (RX1)
    924.5 - SF7BW500 to SF12BW500 (RX1)
    925.1 - SF7BW500 to SF12BW500 (RX1)
    925.7 - SF7BW500 to SF12BW500 (RX1)
    926.3 - SF7BW500 to SF12BW500 (RX1)
    926.9 - SF7BW500 to SF12BW500 (RX1)
    927.5 - SF7BW500 to SF12BW500 (RX1)
    923.3 - SF12BW500 (RX2)

AS923

Use two frequency plans based on country/region. OTAA devices use two common channels: 923.2MHz and 923.4MHz. They will receive the additional channels on a successful join.
```