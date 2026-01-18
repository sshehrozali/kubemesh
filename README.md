# Ethernet Packet Offset Reference

This guide provides the byte indexes for slicing raw network data starting from the Ethernet header (Layer 2).

### Packet Header Map

| Component | Byte Index | Size (Bytes) | Description |
| :--- | :--- | :--- | :--- |
| **MAC Destination** | `0:6` | 6 | Destination hardware address |
| **MAC Source** | `6:12` | 6 | Source hardware address |
| **EtherType** | `12:14` | 2 | Protocol type (IPv4 is `0x0800`) |
| **IP Header Size** | `14` | 1 | (Byte & 0x0F) * 4 |
| **Source IP** | `26:30` | 4 | Sender IP address |
| **Destination IP** | `30:34` | 4 | Receiver IP address |
| **Source Port** | `34:36` | 2 | Sender TCP port |
| **Destination Port** | `36:38` | 2 | Receiver TCP port |
| **Sequence No** | `38:42` | 4 | TCP sequence tracking number |
| **TCP Header Size** | `46` | 1 | (Byte >> 4) * 4 |
| **Flags** | `47` | 1 | ACK, PSH, FIN, SYN bits |



---

### Slicing Logic

Because IP and TCP headers can have "Options," their sizes are dynamic. To find the exact start of your **Payload** (Data), use the following index calculation:

**Payload Start Index** = `14 (Ethernet)` + `IP_Header_Size` + `TCP_Header_Size`

---