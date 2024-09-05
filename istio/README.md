# Istio

This directory contains configurations related to the circuit breaker setup in Istio.

## Circuit Breaker

The `circuit-breaker` directory includes:

- **DestinationRule**: Defines the policies for controlling traffic to services, including circuit-breaking behavior.
- **Sidecar Manifests**: Does the same like DestinationRule, but only for the server.

These manifests are used to apply rate-limiting and other traffic control policies in the Istio environment.
