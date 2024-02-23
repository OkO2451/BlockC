# simple blockchain thingy

## testing
go test -v -run TestPrepareData
go test -v -run TestRun

## stuff to improve

- links between blocks
- use rest instead of cli
- impliment mem pool


## choice of the curve

1.P-256 (also known as prime256v1 or secp256r1):
This is part of the NIST (National Institute of Standards and Technology) suite of curves. It's widely used and supported, including in U.S. government systems. However, it has been criticized because its parameters were selected in a way that isn't fully transparent, leading to suspicions of a potential backdoor.

2.secp256k1:
This curve is used by Bitcoin and Ethereum for their public key cryptography. It's not widely used outside of blockchain technologies. Its main advantage is that it's slightly faster than P-256.

3.Curve25519: 
This curve was designed by Daniel J. Bernstein with speed, security, and simplicity in mind. It's used in many modern protocols, including Signal and WireGuard. It's generally considered to be very secure and efficient.

4.Ed25519: 
This is a specific implementation of EdDSA (Edwards-curve Digital Signature Algorithm) using the Curve25519. It's designed to be fast and to provide a high level of security, even against attackers with physical access to the hardware.

5.P-384 and P-521: 
These are also part of the NIST suite of curves. They provide a higher level of security than P-256, but they're slower and not as widely supported.