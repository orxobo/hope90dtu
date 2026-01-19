# E90-DTU Protocol Specification

## Overview
The E90-DTU communicates over UDP using a custom binary protocol with specific command structure and response format. The device listens on UDP port 8886 and responds to commands with structured binary responses.

## Connection Details
- **Protocol**: UDP
- **Port**: 8886
- **IP Address**: 192.168.68.113 (example)
- **Timeout**: 2-3 seconds for responses

## Command Structure
Commands follow a specific binary format:
```
[c0][c1][c2][c3][Command][Data]
```

Where:
- `c0 c1 c2 c3` = Fixed header (0xc0 0xc1 0xc2 0xc3)
- `Command` = Command identifier (1 byte)
- `Data` = Variable data (0 or more bytes)

## Command List

### 1. RSSI Query (0x01)
- **Command**: `c0 c1 c2 c3 00 01`
- **Description**: Query current RSSI value
- **Response**: `c1 00 01 9f`
- **Interpretation**: 
  - `c1` = Response header
  - `00` = Command type (query)
  - `01` = Command ID (RSSI)
  - `9f` = RSSI value (159 in decimal)

### 2. Device Info (0x02)
- **Command**: `c0 c1 c2 c3 00 02`
- **Description**: Request device information
- **Response**: `c1 00 02 9f 00`
- **Interpretation**:
  - `c1` = Response header
  - `00` = Command type (query)
  - `02` = Command ID (device info)
  - `9f` = First data byte
  - `00` = Second data byte

### 3. Status Query (0x03)
- **Command**: `c0 c1 c2 c3 00 03`
- **Description**: Request device status
- **Response**: Timeout (I/O timeout)
- **Note**: May require specific parameters or timing

### 4. Channel Info (0x04)
- **Command**: `c0 c1 c2 c3 00 04`
- **Description**: Request channel information
- **Response**: Timeout (I/O timeout)
- **Note**: May require specific parameters or timing

### 5. Configuration Query (0x05)
- **Command**: `c0 c1 c2 c3 00 05`
- **Description**: Request configuration
- **Response**: Timeout (I/O timeout)
- **Note**: May require specific parameters or timing

### 6. Send Packet (0x01)
- **Command**: `c0 c1 c2 c3 01 00`
- **Description**: Send packet to network
- **Response**: Timeout (I/O timeout)
- **Note**: May require specific packet data format

## Response Format
Responses follow a consistent format:
```
[Response Header][Command Type][Command ID][Data][Checksum]
```

## Known Issues
1. Most commands except RSSI and device info queries timeout
2. Commands requiring packet data may need specific formatting
3. Timing or parameter requirements for some commands are unknown

## Mesh Node Creation
To create a mesh node that T-Deck can recognize:

1. **Send a packet command**: `c0 c1 c2 c3 01 00`
2. **Follow with mesh packet data**: Properly formatted mesh packet
3. **Node announcement packet**: May need specific mesh protocol format

## Recommendations
1. Focus on the working commands (RSSI and device info) for initial testing
2. Investigate the specific packet format needed for mesh communications
3. Consider using existing mesh protocol tools for proper packet generation
4. The E90-DTU likely has its own packet format specific to its mesh implementation
5. T-Deck requires valid node announcement packets to recognize new nodes

## Notes
- The protocol uses a fixed header `c0 c1 c2 c3` for all commands
- Response headers start with `0xc1`
- Commands with 0x00 as command type appear to be queries
- Commands with 0x01 as command type appear to be for sending data
- Timeout issues suggest the device may require specific parameters or timing for certain operations