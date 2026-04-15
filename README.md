# Simulating Instability

- randomly drop packets

- Bit Flipping: Randomly change a byte in the slice before forwarding.
- Delay: Use time.AfterFunc to forward the packet after a random delay.
- Duplication: Send the same WriteToUDP twice occasionally.

# Packet

Fields:
- id
- timestamp
- len (amount of payload actually used)
- payload (fixed size but large enough)

# Sending packets

- must flatten struct into a bytes buffer
- must also force big endian (use encoding/binary.Write to write)
- payload can be of type []bytes, so no need to encode. write directly

