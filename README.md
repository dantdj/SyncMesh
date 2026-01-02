# SyncMesh
Service that allows for file synchronization amongst connected peers. Intended as a toy project to gain some experience with peer-to-peer networking, NAT traversal, etc, while having some real-world utility.

## Rough planned structure

### Signalling server
This provides a separate service for peer discovery and connection establishment. This is effectively a STUN/TURN server, allowing some traversal of NAT.

In a production version of this project, this would be a cloud-hosted service to allow all of the clients to connect. 

### Local client
This is the main client, which runs on the user's machine. It handles the actual file synchronization and peer-to-peer communication. It will connect to the signalling server to discover other peers and establish connections.

For the sake of ease, it will also have a web API layer to allow for controlling various bits of functionality (resyncing clients, adding new files, etc). This saves on implementing a desktop UI, which isn't what I'm trying to learn here.