# E90-DTU(900SL22)E Test Tool
 This is a personal test tool to help with debugging and communication with E90-DTU devices on Linux/Windows/Android, as only a Windows tool exists from the manufacturer.
 It is in no way associated with the device manufacturer.

## Notes:
* Frequency id 850.125 to 930.125 in 1Mhz increments, ie. channels 1-80
* Module address id for connecting to other E90 modules, like a filter 0-65535 where 65535 is broadcast/monitor all (filterless)
* PP25-PP29 of LoRa Alliance regional parameters (LARP) frequency guide AU915-928
* info https://meshtastic.au/wp/?page_id=47

## Device details
* E90-DTU(900SL22)E
* No serial interface nor dip switches, only an Ethernet port.
* The serial option is added as future functionality as most people have the serial interface. 
https://www.cdebyte.com/products/E90-DTU(900SL22)E
https://www.cdebyte.com/Resources-Download

## Trying to set up as meshtastic node
* This does not seem possible.
https://github.com/meshtastic
https://meshtastic.org/

### Frequencies:
```
AU915-928

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

AS923
Use two frequency plans based on country/region. OTAA devices use two common channels: 923.2MHz and 923.4MHz.
They will receive the additional channels on a successful join.
```