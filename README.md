# AutoDiscoverServer

AutoDiscoverServer is a lightweight UDP server designed to simplify service discovery in networked applications. It listens for discovery requests from clients and responds with essential information, allowing clients to establish reliable connections.

## Features

Service Discovery: AutoDiscoverServer facilitates easy service discovery within a networked environment. Clients can send a discovery packet to the server, and it responds with the necessary information to connect to a TCP server.

## Usage

Build: Build the AutoDiscoverServer using your preferred Go build tools.

```bash
go build autodiscover.go
```
Run: Execute the compiled binary to start the server.

```bash
./autodiscover
```
1. Discovery: Clients can discover the TCP server's address and port by sending a structured UDP packet of type DiscoveryRequest with the Command field set to "REQUEST" to the server's UDP port (default is 10001).

2. Upon receiving a valid DiscoveryRequest, the server will respond with a structured message of type DiscoveryResponse, containing the TCP server's IP address and port. Clients can use this information to establish a reliable connection.

## Configuration

You can modify the following constants in the code to suit your specific requirements:

* maxDatagramSize: Maximum size of incoming UDP packets.
* servicePort: UDP port on which the server listens for discovery requests.
* requestString: The string that clients must send in a discovery request.

## Security Considerations

AutoDiscoverServer focuses on simplicity and ease of use. However, it's important to note that it does not include authentication or encryption. If security is a concern, consider implementing additional security measures based on your specific needs.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

AutoDiscoverServer was created to simplify service discovery in networked applications. We appreciate the open-source community and contributions from developers worldwide.

Feel free to contribute, report issues, or suggest improvements to make AutoDiscoverServer even better!
