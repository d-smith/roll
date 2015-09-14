# Roll Design Notes

Roll is a project that provides some basic functionality for registering users
and applications in a developer portal, for registered applications to obtain
access tokens for making API calls, and for environments providing API services to
grant access only to applications that have obtained access tokens.

These notes will talk about the use of this software in the context of the XTRAC development
portal and APIs, but the roll project is general enough for reuse in a variety of contexts.

## Overview

As part of our UI and API strategies, XTRAC will be building a developer portal for consumers of
UI components and APIs. The scope of this document is to discuss application registration and
access from the viewpoint of the API project.

From the API perspective, we want to create a developer portal to make it easy for developers
to create applications that integrate with XTRAC through our API set. Such as portal provides
the following:

* A simple mechanism to register as an XTRAC developer and obtain instant access to an
environment that can be used to experiment with API calls.
* The ability to register XTRAC applications and obtain an API key and API secret for the
application.
* An interactive quick start guide that introduces the developer to XTRAC workflow as
experienced via the API.
* Live API reference documentation that enables developers to form and execute API calls
from the documentation.
* Full documentation that includes best practices, reference architectures and code, and
task-oriented guides that help developers achieve certain classes of activities, such as
building web and mobile applications, integrating line of business systems with XTRAC,
and submitting XTRAC transactions on behalf of others.
* Access to a collaboration platform to connect developers to each other to facilitate
building a community of XTRAC developers.

The scope of this document is to discuss:

* Services for registering users and applications.
* Mechanisms for registered applications to obtain access tokens.
* Mechanisms for protecting API endpoints to restrict access to callers who have valid access 
tokens.



## Quality Attribute Considerations

### Security

### Performance and Scale

### Portability

### Standards Compliance


 