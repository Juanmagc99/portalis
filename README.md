# Portalis

**Portalis** is a simple service discovery tool with built-in load balancing using a round-robin algorithm. It actively monitors and sanitizes service instances, automatically removing those that are unreachable or fail health checks, ensuring high availability and reliability for your applications.

---

## Features

- **Round Robin Load Balancing**: Efficiently distributes client requests evenly across available service instances.
- **Automatic Instance Sanitization**: Periodically checks and removes service instances that are down or unhealthy.
- **Go Compatibility**: Seamless integration with Go applications through an intuitive client https://github.com/Juanmagc99/portalis.
