 This is a test tool to help with debugging and communication with E90-DTU devices.

 I have an E90-DTU(900SL22)E fro Ebyte with the attached user manual. I want to be able to communicate through it with my Lillygo T-deck using meshtastic.
The E90-DTU(900SL22)E does not have a serial interface nor dip switches, only an Ethernet port.
On the E90:
* I have set up a UDP server on socket A configured for port 8886.
* In the wireless settings:
    * Channel 69 (919.125Mhz)
    * Net ID 0
    * Air baud 2.4Kbps
    * Packet length 240 Byte
    * Packet RSSI enabled
    * Channel RSSI enabled
    * All other settings are default ()
* It is currently sitting on IP 192.168.68.113
On the T-Deck:
* Device role: client
* Must see a node before it will try to transmit.
* frequency slot 17 919.125Mhz [250khz]
* Channel LONGFAST

I need to make the E90 transmit a meshtastic "Node" so i can examine the packets and communicate via meshtastic.
I need to determine how to use the UDP interface past getting RSSI.
I need you to create a golang app with me that:
* Initially tests to see if we can communicate with the E90,
* Determine more functionality from the UDP server in the E90.
* Create a node that the T-deck will see so i can transmit packts back to the E90,
* Help me examine the returned packets.
 Lets start with checking to see if we can communicate with the E90 and return something.

Note: get RSSI by sending 
TX: c0c1c2c30001
RX: c100019f
   Current noise: -79.5 dBm

https://github.com/meshtastic
https://meshtastic.org/


on GC 
192.168.68.113

bne was 192.168.0.125

Ebyte E90-DTU
current IP 192.168.0.125
MAC 2C-BC-BB-34-25-DF
Serial S5200489S


Firmware FW-9181-0-10

E90-DTU(900SL22)E
https://www.cdebyte.com/products/E90-DTU(900SL22)E

https://www.cdebyte.com/Resources-Download


lillygo t-deck with the meshtastic firmware