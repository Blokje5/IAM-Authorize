# Design

This project is inspired by AWS IAM & Chef Automate IAM. The goal is to build an open source IAM Policy Authorization implementation that can be used to authorize requests in your infrastructure.

## Terms

Principal - Is an enity (e.g. user, service) that can be authenticated. Once a principal makes a request and is authenticated, the IAM-authorize can be used to validate the principals permissions
Action - An action is something that can be invoked by the principal by making a request. For example, CreateUser could be an action.
Resource - Is an object (e.g. server, user, policy) a principal can be authorized to perform actions on.
Policy - A policy defines the set of rules that determine whether a principal is authorized to perform an action on a resource
Namespace - A scope that is applied to both resources and actions. E.g. An action could be scoped to the users namespace.

## Namespace management

Administrators can create namespaces in this service. This provides an unique scope for a service to create policies in. For example, if you have an user management service, you could create a namespace `users`. Actions would be of the form `namespace:action`, e.g. `users:createUser`. The same holds true for resources, which are of the form `namespace:object_identifier` e.g. `users:1234` to identify the user with id 1234. Namespaces are validated on policy creation. This ensures that no incorrect policies are created.
