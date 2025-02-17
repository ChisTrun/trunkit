# MyKit - A toolkit for microservices

## Context & Scope
As we’re adopting Microservices and migrating from Monolith to Microservices, we continue to look for ways to work faster and deliver qualitative and relevant services quickly and efficiently while maintaining consistency, coordination, and high quality across multiple teams.

## Goals and non-goals
### Goals
MyKit - a framework for building Go microservices. MyKit is designed to create a fully functional microservice scaffolding in seconds, allowing engineers to focus on the business logic straight away!

MyKit provides abstraction from all aspects of distributed system design by simplifying the creation and operation of microservices through scaffolding, using smart library configuration defaults, automatic initialization, context propagation, and runtime framework configuration. Moreover, it standardized communication across services.

We will no longer need to spend long hours generating boilerplate code, initializing common libraries, creating dashboards and alarms, or creating Data Access Objects (DAOs). Instead, we can concentrate on delivering scalable and agile services that are essential for the success of our engineers and in turn delight our consumers.

### Non-goals
At the MyKit level, we will not try to cover these categories
Transportation is not supported more than these transportation layers: RESTful API, gRPC from the first place, but not limit for extension in the near future
Rate Limiting, Circuit Breaking will not cover by the kit at the service level, it should be covered at Gateway level, and when it comes to service level should be cover with a sidecar in mesh
Authentication with a package like Go kit and we will cover this with API Gateway which will come in another proposal also

## The actual design
![alt text](./docs/kit-components.png)

## References
- https://engineering.grab.com/introducing-grab-kit
- https://github.com/go-kit/kit
- https://shijuvar.medium.com/go-microservices-with-go-kit-introduction-43a757398183
- https://microservices.io/patterns/apigateway.html
- http://peter.bourgon.org/go-kit/
- https://dapr.io/# mykit
# mykit
