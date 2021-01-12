# SCIM service

Provision and manage user accounts and groups with the SCIM API.
[SCIM is an open standard](http://www.simplecloud.info/) that is used by Single Sign-On (SSO) services and identity providers to manage people across a variety of tools/apps.

SCIM is RESTful API: it uses HTTP and its verbs to call SCIM methods. You will use `GET` to retrieve information (entity or list of entities), `POST` to create new objects, `PUT` to replace them, `PATCH` to modify objects, and `DELETE` to remove.

## Request body and response formats

This API supports several formats of request body (for `POST`, `PUT` and `PATCH`) and of response:
* JSON. In this case the header `Content-type` of request should be `application/json` or `application/scim+json`. And it's `application/json` for response header.
* XML  In this case the header `Content-type` of request should be `text/xml`, `application/xml` or `application/scim+xml`. And it's `application/xml` is for response header.

## Authorization

This API supports several authorization schemes:

* `basic` - Basic HTTP authorization. User and password should be set in instance configuration.
  - Requests were authorization is needed expect `Authorization` header that looks as `Basic <encoded credentials>` 
* `token` - App Bearer Token. Token should be set in instance configuration.
  - Requests were authorization is needed expect `Authorization` header that looks as `Bearer <token>` 

## New parameters in instance configuration

* `scim.enabled` - SCIM is enabled for this instance. Some of users/groups attributes (that are modified by SCIM) can't be changed via our website.
* `scim.auth.type` - type of authorization (`none`, `basic` or `token`).
  * for `basic` there are additional parameters:
    - `scim.auth.basic.user` - credentials for HTTP Basic Authorization
    - `scim.auth.basic.password` - credentials for HTTP Basic Authorization
  * for `token` there is additional parameter:
    - `scim.auth.token` - bearer token
* `scim.response.format` - format of API response (`json` or `xml`).
* `scim.middleware` - possible middleware applied to the request body before making the standard SCIM processing. If missing, empty or has value `default` the default middleware will be applied.  
    * for Rabobank middleware there is additional parameter:
      - `scim.middleware.rabobank.role-mapping` - mapping of the default roles to the Rabobank groups. The structure of the mapping is `role alias` => array of group IDs.  

## Dependencies

See file `Gopkg.toml`

**To add new dependency**
* Add new dependency via `dep ensure -add <path-to-repository`
* For our libs use constraint with the necessary branch/version
* If new dependency has configuration then:
  * Add configuration file to `conf` directory and to `AWS S3` bucket
  * Add downloading of configuration to `entrypoint.sh` file

## Installation

* Clone repository by running the following commands:
  * `git clone git@gitlab.com:mike-jc/service-scim.git`
  * `cd service-scim`
* Install dependencies.
  * install `dep` if not exist:
    `curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`
  * run `dep ensure`
  * run `git install -v`
* Add configuration files.
    - `cp conf/app.conf.default conf/app.conf`
* Test project if necessary.
  * add test configuration:
    - `mkdir tests/conf`
    - `cp conf/app.conf.test tests/conf/app.conf`
  * run `go test ./tests/...`

## Configurations

Edit `app.conf` file according to your needs:

* `logger.*` Subset of parameters for logger
* `configurator.url` URL to our internal configurator (that manage instances configuration)
* `jwtSecret`, `jwtLeeway` - parameters for generating/parsing JWT tokens
* `sp.*` Subset of settings for SCIM service: enable some operations/features and set their parameters
* `repository.*` Subset of parameters for repository via which we work with data.
  * Possible types of repository engines: `fake` (for test purposes), `file` (for running locally) or `keeper`. 

## Run service

* Add entry in Git CI configuration in `AWS S3`, in `docker-compose.yml` file

# Endpoints

## Health check

Summary: Checks application status. Must be used by monitor.

Endpoint: /healthcheck

Method: **GET**

Auth header: **No Auth**

Example: GET: /healthcheck

Response, status **200**

```json
{
    "status": "success"
}
```

Response, status **500**

```json
{
    "error": "configuration_error",
    "status": "error"
}
```

## Service Provider Configuration

Summary: configuration details for our SCIM API, including which operations are supported.

Endpoint: /ServiceProviderConfigs

Method: **GET**

Auth header: **No Auth**

Example: GET /ServiceProviderConfigs

Response, status **200**

```json
{
  "schemas": [
    "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"
  ],
  "id": "ServiceProviderConfig",
  "patch": {
    "supported":true
  },
  "bulk": {
    "supported":false,
    "maxOperations":0,
    "maxPayloadSize":0
  },
  "filter": {
    "supported":true,
    "maxResults": 250
  },
  "changePassword": {
    "supported":false
  },
  "sort": {
    "supported":false
  },
  "etag": {
    "supported":true
  },
  "authenticationSchemes": [
    {
      "name": "HTTP Basic",
      "description": "Authentication scheme using the HTTP Basic Standard",
      "specUri": "http://www.rfc-editor.org/info/rfc2617",
      "type": "httpbasic",
      "primary": true
    }
  ],
  "meta": {
    "location": "http://127.0.0.1:8101/v2/ServiceProviderConfig",
    "resourceType": "ServiceProviderConfig",
    "created": "2019-04-04T17:43:00Z",
    "lastModified": "2019-04-04T17:43:00Z",
    "version": "v2"
  }
}
```

## Resource Types

Summary: the types of resources available.

Endpoints:
* /resourcetypes
* /resourcetypes/:name
  - where `:name` is one of `user` or `group` 

Method: **GET**

Auth header: **No Auth**

Example: GET /resourcetypes

Response, status **200**

```json
{
  "schemas": [
    "urn:ietf:params:scim:api:messages:2.0:ListResponse"
  ],
  "totalResults": 2,
  "itemsPerPage": 2,
  "startIndex": 1,
  "resources": [
    {
      "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:ResourceType"
      ],
      "id": "User",
      "name": "User",
      "endpoint": "/Users",
      "description": "User account: https://tools.ietf.org/html/rfc7643#section-8.7.1",
      "schema": "urn:ietf:params:scim:schemas:core:2.0:User",
      "meta": {
        "location": "http://127.0.0.1:8101/v2/ResourceTypes/User",
        "resourceType": "ResourceType"
      }
    },
    {
      "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:ResourceType"
      ],
      "id": "Group",
      "name": "Group",
      "endpoint": "/Groups",
      "description": "Group of user: https://tools.ietf.org/html/rfc7643#section-8.7.1",
      "schema": "urn:ietf:params:scim:schemas:core:2.0:Group",
      "meta": {
        "location": "http://127.0.0.1:8101/v2/ResourceTypes/Group",
        "resourceType": "ResourceType"
      }
    }
  ]
}
```

Error, status **400**
```json
{
    "schemas": [
        "urn:ietf:params:scim:api:messages:2.0:Error"
    ],
    "status": "400",
    "scimType": "app.unknown_resource_type",
    "detail": "Unknown resource type"
}
```

## Schemas

Summary: schemas for users and groups. Querying the schemas will provide the most up-to-date rendering of the supported SCIM attributes.

Endpoints:
* /schemas
* /schemas/:urn
  - where `:urn` can be one of `urn:ietf:params:scim:schemas:core:2.0:User` or `urn:ietf:params:scim:schemas:core:2.0:Group`

Method: **GET**

Auth header: **No Auth**

Example: GET /schemas

Response, status **200**

```json
{
  "schemas": [
    "urn:ietf:params:scim:api:messages:2.0:ListResponse"
  ],
  "totalResults": 3,
  "itemsPerPage": 3,
  "startIndex": 1,
  "Resources": [
    {
      "attributes": [
        {
          "name": "schemas",
          "type": "reference",
          "multiValued": true,
          "required": true,
          "caseExact": true,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "An array of Strings containing URIs that are used to indicate the namespaces of the SCIM schemas that define the attributes present in the current structure.",
          "canonicalValues": [
            "urn:ietf:params:scim:schemas:core:2.0:User",
            "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
          ],
          "referenceTypes": [
            "uri"
          ]
        },
        {
          "name": "id",
          "type": "integer",
          "multiValued": false,
          "required": true,
          "caseExact": true,
          "mutability": "readOnly",
          "returned": "always",
          "uniqueness": "global",
          "description": "A unique identifier for a SCIM resource as defined by the service provider."
        },
        {
          "name": "userName",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "server",
          "description": "Unique identifier for the User, typically used by the user to directly authenticate to the service provider. Each User MUST include a non-empty userName value.  This identifier MUST be unique across the service provider's entire set of Users."
        },
        {
          "name": "externalId",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": true,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "global",
          "description": "A String that is an identifier for the resource as defined by the provisioning client."
        },
        {
          "name": "emails",
          "type": "complex",
          "multiValued": true,
          "required": true,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "Email addresses for the user.  The value SHOULD be canonicalized by the service provider, e.g., 'bjensen@example.com' instead of 'bjensen@EXAMPLE.COM'. Canonical type values of 'work', 'home', and 'other'.",
          "subAttributes": [
            {
              "name": "value",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "Email addresses for the user.  The value SHOULD be canonicalized by the service provider, e.g., 'bjensen@example.com' instead of 'bjensen@EXAMPLE.COM'. Canonical type values of 'work', 'home', and 'other'."
            },
            {
              "name": "display",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "A human-readable name, primarily used for display purposes."
            },
            {
              "name": "type",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "A label indicating the attribute's function, e.g., 'work' or 'home'.",
              "canonicalValues": [
                "work",
                "home",
                "other"
              ]
            },
            {
              "name": "primary",
              "type": "boolean",
              "multiValued": false,
              "required": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "A Boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g., the preferred mailing address or primary email address.  The primary attribute value 'true' MUST appear no more than once."
            }
          ]
        },
        {
          "name": "name",
          "type": "complex",
          "multiValued": false,
          "required": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "The components of the user's real name. Providers MAY return just the full name as a single string in the formatted sub-attribute, or they MAY return just the individual component attributes using the other sub-attributes, or they MAY return both.  If both variants are returned, they SHOULD be describing the same name, with the formatted name indicating how the component attributes should be combined.",
          "subAttributes": [
            {
              "name": "formatted",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The full name, including all middle names, titles, and suffixes as appropriate, formatted for display (e.g., 'Ms. Barbara J Jensen, III')."
            },
            {
              "name": "familyName",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The family name of the User, or last name in most Western languages (e.g., 'Jensen' given the full name 'Ms. Barbara J Jensen, III')."
            },
            {
              "name": "givenName",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The given name of the User, or first name in most Western languages (e.g., 'Barbara' given the full name 'Ms. Barbara J Jensen, III')."
            },
            {
              "name": "middleName",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The middle name(s) of the User (e.g., 'Jane' given the full name 'Ms. Barbara J Jensen, III')."
            }
          ]
        },
        {
          "name": "profileUrl",
          "type": "reference",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "A fully qualified URL pointing to a page representing the User's online profile.",
          "referenceTypes": [
            "external"
          ]
        },
        {
          "name": "title",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "The user's job title, such as 'Vice President.'"
        },
        {
          "name": "locale",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "Used to indicate the User's default location for purposes of localizing items such as currency, date time format, or numerical representations."
        },
        {
          "name": "timezone",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "The User's time zone in the 'Olson' time zone database format, e.g., 'America/Los_Angeles'."
        },
        {
          "name": "active",
          "type": "boolean",
          "multiValued": false,
          "required": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "A Boolean value indicating the User's administrative status."
        },
        {
          "name": "password",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "writeOnly",
          "returned": "never",
          "uniqueness": "none",
          "description": "The User's cleartext password. This attribute is intended to be used as a means to specify an initial password when creating a new User or to reset an existing User's password."
        },
        {
          "name": "phoneNumbers",
          "type": "complex",
          "multiValued": true,
          "required": false,
          "mutability": "readWrite",
          "returned": "default",
          "description": "Phone numbers for the User. The value SHOULD be canonicalized by the service provider according to the format specified in RFC 3966, e.g., 'tel:+1-201-555-0123'. Canonical type values of 'work', 'home', 'mobile', 'fax', 'pager', and 'other'.",
          "subAttributes": [
            {
              "name": "value",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "Phone number of the User."
            },
            {
              "name": "display",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "A human-readable name, primarily used for display purposes."
            },
            {
              "name": "type",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "A label indicating the attribute's function, e.g., 'work', 'home', 'mobile'.",
              "canonicalValues": [
                "work",
                "home",
                "mobile",
                "fax",
                "other"
              ]
            },
            {
              "name": "primary",
              "type": "boolean",
              "multiValued": false,
              "required": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "A Boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g., the preferred phone number or primary phone number.  The primary attribute value 'true' MUST appear no more than once."
            }
          ]
        },
        {
          "name": "photos",
          "type": "complex",
          "multiValued": true,
          "required": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "URLs of photos of the User.",
          "subAttributes": [
            {
              "name": "value",
              "type": "reference",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "URL of a photo of the User.",
              "referenceTypes": [
                "external"
              ]
            },
            {
              "name": "type",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "A label indicating the attribute's function, i.e., 'photo' or 'thumbnail'.",
              "canonicalValues": [
                "photo",
                "thumbnail"
              ]
            }
          ]
        },
        {
          "name": "addresses",
          "type": "complex",
          "multiValued": true,
          "required": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "A physical mailing address for this User.  Canonical type values of 'work', 'home', and 'other'.  This attribute is a complex type with the following sub-attributes.",
          "subAttributes": [
            {
              "name": "streetAddress",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The full street address component, which may include house number, street name, P.O. box, and multi-line extended street address information. This attribute MAY contain newlines."
            },
            {
              "name": "locality",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The city or locality component."
            },
            {
              "name": "region",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The state or region component."
            },
            {
              "name": "postalCode",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The zip code or postal code component."
            },
            {
              "name": "country",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The country name component."
            },
            {
              "name": "type",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "A label indicating the attribute's function, e.g., 'work' or 'home'.",
              "canonicalValues": [
                "work",
                "home",
                "other"
              ]
            },
            {
              "name": "primary",
              "type": "boolean",
              "multiValued": false,
              "required": false,
              "mutability": "readWrite",
              "returned": "default",
              "description": "A Boolean value indicating the 'primary' or preferred attribute value for this attribute, e.g., the preferred messenger or primary messenger.  The primary attribute value 'true' MUST appear no more than once."
            }
          ]
        },
        {
          "name": "groups",
          "type": "complex",
          "multiValued": true,
          "required": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "A list of groups to which the user belongs.",
          "subAttributes": [
            {
              "name": "value",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The identifier of the User's group."
            },
            {
              "name": "$ref",
              "type": "reference",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The URI of the corresponding 'Group' resource to which the user belongs.",
              "referenceTypes": [
                "Group"
              ]
            },
            {
              "name": "display",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "A human-readable name, primarily used for display purposes."
            }
          ]
        },
        {
          "name": "roles",
          "type": "complex",
          "multiValued": true,
          "required": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "A list of roles for the User that collectively represent who the User is, e.g., 'Manager', 'Operator'.",
          "subAttributes": [
            {
              "name": "value",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The value of a role."
            },
            {
              "name": "display",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "A human-readable name, primarily used for display purposes."
            }
          ]
        },
        {
          "name": "meta",
          "type": "complex",
          "multiValued": false,
          "required": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "Resource metadata.",
          "subAttributes": [
            {
              "name": "resourceType",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": true,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The name of the resource type of the resource."
            },
            {
              "name": "location",
              "type": "reference",
              "multiValued": false,
              "required": false,
              "caseExact": true,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The URI of the resource being returned.",
              "referenceTypes": [
                "uri"
              ]
            },
            {
              "name": "created",
              "type": "datetime",
              "multiValued": false,
              "required": false,
              "caseExact": true,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The 'DateTime' that the resource was added to the service provider."
            }
          ]
        }
      ],
      "id": "urn:ietf:params:scim:schemas:core:2.0:User",
      "name": "User",
      "description": "User Schema",
      "meta": {
        "location": "http://41e106ae.ngrok.io/v2/Schemas/urn:ietf:params:scim:schemas:core:2.0:User",
        "resourceType": "Schema"
      }
    },
    {
      "attributes": [
        {
          "name": "employeeNumber",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "Numeric or alphanumeric identifier assigned to a person, typically based on order of hire or association with an organization."
        },
        {
          "name": "organization",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "Identifies the name of an organization."
        },
        {
          "name": "division",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "Identifies the name of a division."
        },
        {
          "name": "department",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "Identifies the name of a department."
        },
        {
          "name": "manager",
          "type": "complex",
          "multiValued": false,
          "required": false,
          "mutability": "readWrite",
          "returned": "default",
          "description": "The User's manager.  A complex type that optionally allows service providers to represent organizational hierarchy by referencing the 'id' attribute of another User.",
          "subAttributes": [
            {
              "name": "value",
              "type": "string",
              "multiValued": false,
              "required": true,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The id of the SCIM resource representing the User's manager."
            },
            {
              "name": "$ref",
              "type": "reference",
              "multiValued": false,
              "required": true,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "The URI of the SCIM resource representing the User's manager.",
              "referenceTypes": [
                "User"
              ]
            },
            {
              "name": "displayName",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The displayName of the User's manager."
            }
          ]
        }
      ],
      "id": "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
      "name": "EnterpriseUser",
      "description": "Enterprise User Schema",
      "meta": {
        "location": "http://41e106ae.ngrok.io/v2/Schemas/urn:ietf:params:scim:schemas:extension:enterprise:2.0:User",
        "resourceType": "Schema"
      }
    },
    {
      "attributes": [
        {
          "name": "schemas",
          "type": "reference",
          "multiValued": true,
          "required": true,
          "caseExact": true,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "An array of Strings containing URIs that are used to indicate the namespaces of the SCIM schemas that define the attributes present in the current structure.",
          "canonicalValues": [
            "urn:ietf:params:scim:schemas:core:2.0:Group"
          ],
          "referenceTypes": [
            "uri"
          ]
        },
        {
          "name": "id",
          "type": "integer",
          "multiValued": false,
          "required": true,
          "caseExact": true,
          "mutability": "readOnly",
          "returned": "always",
          "uniqueness": "global",
          "description": "A unique identifier for a SCIM resource as defined by the service provider."
        },
        {
          "name": "externalId",
          "type": "string",
          "multiValued": false,
          "required": false,
          "caseExact": true,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "global",
          "description": "A String that is an identifier for the resource as defined by the provisioning client."
        },
        {
          "name": "displayName",
          "type": "string",
          "multiValued": false,
          "required": true,
          "caseExact": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "server",
          "description": "A human-readable name for the Group."
        },
        {
          "name": "members",
          "type": "complex",
          "multiValued": true,
          "required": false,
          "mutability": "readWrite",
          "returned": "default",
          "uniqueness": "none",
          "description": "A list of members of the Group.",
          "subAttributes": [
            {
              "name": "value",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readWrite",
              "returned": "default",
              "uniqueness": "none",
              "description": "Identifier of the member of this Group."
            },
            {
              "name": "$ref",
              "type": "reference",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The URI corresponding to a SCIM user that is a member of this Group.",
              "referenceTypes": [
                "User"
              ]
            },
            {
              "name": "display",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "A human-readable name for the group member."
            },
            {
              "name": "type",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": false,
              "mutability": "immutable",
              "returned": "default",
              "uniqueness": "none",
              "description": "A label indicating the type of resource, e.g., 'User' or 'Group'.",
              "canonicalValues": [
                "User"
              ]
            }
          ]
        },
        {
          "name": "meta",
          "type": "complex",
          "multiValued": false,
          "required": false,
          "mutability": "readOnly",
          "returned": "default",
          "uniqueness": "none",
          "description": "Resource metadata.",
          "subAttributes": [
            {
              "name": "resourceType",
              "type": "string",
              "multiValued": false,
              "required": false,
              "caseExact": true,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The name of the resource type of the resource."
            },
            {
              "name": "location",
              "type": "reference",
              "multiValued": false,
              "required": false,
              "caseExact": true,
              "mutability": "readOnly",
              "returned": "default",
              "uniqueness": "none",
              "description": "The URI of the resource being returned.",
              "referenceTypes": [
                "uri"
              ]
            }
          ]
        }
      ],
      "id": "urn:ietf:params:scim:schemas:core:2.0:Group",
      "name": "Group",
      "description": "Group Schema",
      "meta": {
        "location": "http://41e106ae.ngrok.io/v2/Schemas/urn:ietf:params:scim:schemas:core:2.0:Group",
        "resourceType": "Schema"
      }
    }
  ]
}
```

Error, status **400**
```json
{
    "schemas": [
        "urn:ietf:params:scim:api:messages:2.0:Error"
    ],
    "status": "400",
    "scimType": "app.unknown_schema",
    "detail": "Unknown schema"
}
```

## Attributes

For the next endpoints where users/groups can be requested or are returned in the response the client has mechanism to say which attributes of the resource are need to be added to the response and which - to be excluded from.
For this purpose the following (optional) query parameters are used:
* `attributes` - coma separated list of resource attributes that should be **included** to the response
* `excludedAttributes` - coma separated list of resource attributes that should be **excluded** to the response

For both parameters attribute complex path can be used where all parent attributes' names are concatenated to the necessary attribute name by period.

**Example**:

Let's assume we have the following user object:
```json
{
  "schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
  "externalId": "ef3b507c-d973-4d3f-821a-1fe5277f0af4",
  "name": {
    "formatted": "Anne X",
    "familyName": "X",
    "givenName": "Anne"
  },
  "profileUrl": "https://login.example.com/annex",
  "emails": [
    {
      "value": "anne@example.com",
      "type": "work",
      "primary": true
    },
    {
      "value": "anne@home.com",
      "type": "home"
    }
  ],
  "phoneNumbers": [
    {
      "value": "123-456-0789",
      "type": "work"
    }
  ],
  "photos": [
    {
      "value": "https://photos.example.com/profilephoto/72930000000Ccne/F",
      "type": "photo"
    },
    {
      "value": "https://photos.example.com/profilephoto/72930000000Ccne/T",
      "type": "thumbnail"
    }
  ],
  "title": "Founder",
  "locale": "en-US",
  "timezone": "Asia/Shanghai",
  "active":true,
  "meta": {
    "resourceType": "User",
    "location": "/v2/Users/ef3b507c-d973-4d3f-821a-1fe5277f0af4"
  }
}
```
then you may request only some of user's attributes by adding the following GET-parameter:
```
?attributes="externalId,name.formatted,emails.value,emails.work,photos"
```
In that case you'll received the following object:
```json
{
  "schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
  "externalId": "ef3b507c-d973-4d3f-821a-1fe5277f0af4",
  "name": {
    "formatted": "Anne X"
  },
  "emails": [
    {
      "value": "anne@example.com",
      "type": "work"
    },
    {
      "value": "anne@home.com",
      "type": "home"
    }
  ],
  "photos": [
    {
      "value": "https://photos.example.com/profilephoto/72930000000Ccne/F",
      "type": "photo"
    },
    {
      "value": "https://photos.example.com/profilephoto/72930000000Ccne/T",
      "type": "thumbnail"
    }
  ],
  "meta": {
    "resourceType": "User",
    "location": "/v2/Users/ef3b507c-d973-4d3f-821a-1fe5277f0af4"
  }
}
```
**Be aware!** Attributes `schemas` and `meta` are always included to the objected in response.

## Pagination

For the endpoints where users/groups can be requested as the list of entities the client has mechanism to paginate the result, that is to say which part of the list to return by using the following GET-parameters:
* `startIndex` - 1-based index that determines the start of the page returned in the result.
* `count` - count of entities in the page to return.

For example,

* `?startIndex=1&count=20` - asks the API to return first 20 entities
* `?startIndex=31&count=53` - asks the API to return 53 entities starting from 31st   

## Middleware

Different SCIM clients may make requests with bodies that differ (at least slightly).
For example, some of them may send us the user name in several fields separately while we accept it as the value of the one field.

To make this all work together we introduce the middleware - the processor that adapt input data to our standard structure and then passes it further to other processors.

There is a default middleware that is always applied. For the particular instance you can set the custom middleware by adding the setting `scim.middleware` to its configuration.
Then you need to add the middleware type to SCIM repository to `services/middleware` directory and change the factory code there so that it knows what middleware name corresponds to the new type.
 
 

## Patch operations

Not implemented yet

## Filters

Not implemented yet

## Users

### List of users

Summary: the non-blocked users that exist on our Service Provider side.

Endpoint: /users

Method: **GET**

Auth header: according to the chosen auth type 

Example: GET /users

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:api:messages:2.0:ListResponse"
    ],
    "totalResults": 2,
    "itemsPerPage": 250,
    "startIndex": 1,
    "resources": [
        {
            "schemas": [
                "urn:ietf:params:scim:schemas:core:2.0:User"
            ],
            "externalId": "ef3b507c-d973-4d3f-821a-1fe5277f0af4",
            "email": "anne@example.com",
            "name": "Anne Smith",
            "profileUrl": "https://login.example.com/annex",
            "jobTitle": "Founder",
            "locale": "en-US",
            "timezone": "Asia/Shanghai",
            "active": true,
            "phoneNumber": "123-456-0789",
            "photos": [
                {
                    "value": "https://photos.example.com/profilephoto/72930000000Ccne/F",
                    "type": "photo"
                },
                {
                    "value": "https://photos.example.com/profilephoto/72930000000Ccne/T",
                    "type": "thumbnail"
                }
            ],
            "addresses": [
                {
                    "streetAddress": "123 Home Street",
                    "locality": "Shanghai",
                    "postalCode": "200124",
                    "country": "CN"
                }
            ],
            "meta": {
                "resourceType": "User",
                "location": "/v2/users/ef3b507c-d973-4d3f-821a-1fe5277f0af4"
            }
        },
        {
            "schemas": [
                "urn:ietf:params:scim:schemas:core:2.0:User",
                "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
            ],
            "externalId": "51db7ce8-75d7-4c81-ad05-c98a9233811e",
            "email": "jack@example.com",
            "name": "Jack Green",
            "jobTitle": "Founder",
            "locale": "en-US",
            "timezone": "Asia/Shanghai",
            "active": true,
            "photos": [
                {
                    "value": "https://photos.example.com/profilephoto/72930000000Ccne/F",
                    "type": "photo"
                }
            ],
            "addresses": [
                {
                    "streetAddress": "157 Main Street",
                    "locality": "Beijing",
                    "postalCode": "12346",
                    "country": "CN"
                }
            ],
            "meta": {
                "resourceType": "User",
                "location": "/v2/users/51db7ce8-75d7-4c81-ad05-c98a9233811e"
            },
            "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": {
              "organization": "24s",
              "department": "Dev"
            }
        }
    ]
}
```

Error status codes: **500**, **400**

### User by ID

Summary: the non-blocked user that exists on our Service Provider side and has the given ID.

Endpoint: /users/`id`

Method: **GET**

Auth header: according to the chosen auth type 

Example: GET /users/51db7ce8-75d7-4c81-ad05-c98a9233811e

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:User",
        "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User"
    ],
    "externalId": "51db7ce8-75d7-4c81-ad05-c98a9233811e",
    "email": "jack@example.com",
    "name": "Jack Green",
    "jobTitle": "Founder",
    "locale": "en-US",
    "timezone": "Asia/Shanghai",
    "active": true,
    "photos": [
        {
            "value": "https://photos.example.com/profilephoto/72930000000Ccne/F",
            "type": "photo"
        }
    ],
    "addresses": [
        {
            "streetAddress": "157 Main Street",
            "locality": "Beijing",
            "postalCode": "12346",
            "country": "CN"
        }
    ],
    "meta": {
        "resourceType": "User",
        "location": "/v2/users/51db7ce8-75d7-4c81-ad05-c98a9233811e"
    },
   "urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": {
     "organization": "24s",
     "department": "EngineeringS"
   }
}
```

Error status codes: **500**, **404**

### Creation

Summary: create new user that has unique ID. 

Endpoint: /users

Method: **POST**

Auth header: according to the chosen auth type 

Example: POST /users

Request body:
```json
{
  "email": "willy@example.com",
  "name": "Willy Brown",
  "profileUrl": "https://login.example.com/willy",
  "active":true,
  "addresses": [
    {
      "streetAddress": "57 Wall Street",
      "locality": "New York",
      "country": "US"
    }
  ]
}
```

Response, status **201**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:User"
    ],
    "externalId": "c29ddf00-be78-4a00-8251-3737e7f87d96",
    "email": "willy@example.com",
    "name": "Willy Brown",
    "profileUrl": "https://login.example.com/willy",
    "active": true,
    "addresses": [
        {
            "streetAddress": "57 Wall Street",
            "locality": "New York",
            "country": "US"
        }
    ],
    "meta": {
        "resourceType": "User",
        "location": "/v2/users/c29ddf00-be78-4a00-8251-3737e7f87d96"
    }
}
```

Error status codes: **500**, **400**

### Modification

Summary: partially modify the non-blocked existing user that has the given ID. In other words, you may pass in request only those user's attributes and their values that you want to modify. Other attributes are left untouched.

Endpoint: /users/`id`

Method: **PATCH**

Auth header: according to the chosen auth type 

Example: PATCH /users/c29ddf00-be78-4a00-8251-3737e7f87d96

Request body:
```json
{
  "name": "Willy Wonka"
}
```

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:User"
    ],
    "externalId": "c29ddf00-be78-4a00-8251-3737e7f87d96",
    "email": "willy@example.com",
    "name": "Willy Wonka",
    "profileUrl": "https://login.example.com/willy",
    "active": true,
    "addresses": [
        {
            "streetAddress": "57 Wall Street",
            "locality": "New York",
            "country": "US"
        }
    ],
    "meta": {
        "resourceType": "User",
        "location": "/v2/users/c29ddf00-be78-4a00-8251-3737e7f87d96"
    }
}
```

Error status codes: **500**, **404**, **400**

### Replacement

Summary: fully replace attributes of the non-blocked existing user that has the given ID. In other words, even if you didn't add some attributes and their values to the request, they will be changed (set to their default values).

Endpoint: /users/`id`

Method: **PUT**

Auth header: according to the chosen auth type 

Example: PUT /users/c29ddf00-be78-4a00-8251-3737e7f87d96

Request body:
```json
{
  "externalId": "c29ddf00-be78-4a00-8251-3737e7f87d96",
  "email": "willy@example.com",
  "name": "Willy Brown",
  "profileUrl": "https://login.example.com/wily",
  "jobTitle": "CTO",
  "active":false,
  "addresses": [
    {
      "streetAddress": "123 Other str",
      "locality": "New York",
      "country": "US"
    }
  ]
}
```

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:User"
    ],
    "externalId": "c29ddf00-be78-4a00-8251-3737e7f87d96",
    "email": "willy@example.com",
    "name": "Willy Brown",
    "profileUrl": "https://login.example.com/wily",
    "jobTitle": "CTO",
    "addresses": [
        {
            "streetAddress": "123 Other str",
            "locality": "New York",
            "country": "US"
        }
    ],
    "meta": {
        "resourceType": "User",
        "location": "/v2/users/c29ddf00-be78-4a00-8251-3737e7f87d96"
    }
}
```

Error status codes: **500**, **404**, **400**

### Blocking

Summary: block user that exists on our Service Provider side and has the given ID.

Endpoint: /users

Method: **DELETE**

Auth header: according to the chosen auth type 

Example: DELETE /users/c29ddf00-be78-4a00-8251-3737e7f87d96

Response, status **204**

Error status codes: **500**, **404**

## Groups

### List of groups

Summary: the non-disabled groups that exist on our Service Provider side.

Endpoint: /groups

Method: **GET**

Auth header: according to the chosen auth type 

Example: GET /groups

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:api:messages:2.0:ListResponse"
    ],
    "totalResults": 2,
    "itemsPerPage": 250,
    "startIndex": 1,
    "resources": [
        {
            "schemas": [
                "urn:ietf:params:scim:schemas:core:2.0:Group"
            ],
            "externalId": "6de1cc84-d9d0-4fb4-a6a4-c85c52675887",
            "displayName": "Managers",
            "members": [
                {
                    "value": "anne@example.com",
                    "$ref": "/v2/users/ef3b507c-d973-4d3f-821a-1fe5277f0af4",
                    "type": "User"
                }
            ],
            "meta": {
                "resourceType": "Group",
                "location": "/v2/groups/38c5d26b-070a-42f5-89c8-882326f50b8b"
            }
        },
        {
            "schemas": [
                "urn:ietf:params:scim:schemas:core:2.0:Group"
            ],
            "externalId": "38c5d26b-070a-42f5-89c8-882326f50b8b",
            "displayName": "Operators",
            "members": [
                {
                    "value": "anne@example.com",
                    "$ref": "/v2/users/ef3b507c-d973-4d3f-821a-1fe5277f0af4",
                    "type": "User"
                },
                {
                    "value": "jack@example.com",
                    "$ref": "/v2/users/51db7ce8-75d7-4c81-ad05-c98a9233811e",
                    "type": "User"
                }
            ],
            "meta": {
                "resourceType": "Group",
                "location": "/v2/groups/6de1cc84-d9d0-4fb4-a6a4-c85c52675887"
            }
        }
    ]
}
```

Error status codes: **500**, **400**

### Group by ID

Summary: the non-disabled group that exist on our Service Provider side and has the given ID.

Endpoint: /groups

Method: **GET**

Auth header: according to the chosen auth type 

Example: GET /groups

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:Group"
    ],
    "externalId": "6de1cc84-d9d0-4fb4-a6a4-c85c52675887",
    "displayName": "Managers",
    "members": [
        {
            "value": "anne@example.com",
            "$ref": "/v2/users/ef3b507c-d973-4d3f-821a-1fe5277f0af4",
            "type": "User"
        }
    ],
    "meta": {
        "resourceType": "Group",
        "location": "/v2/groups/38c5d26b-070a-42f5-89c8-882326f50b8b"
    }
}
```

Error status codes: **500**, **404**

### Creation

Summary: create new group that has unique ID.

Endpoint: /groups

Method: **POST**

Auth header: according to the chosen auth type 

Example: POST /groups

Request body:

```json
{
  "displayName": "Developers",
  "members": [
    {
      "value": "peter@example.com",
      "$ref": "/v2/users/a549a5cf-63f8-4788-adab-ddc853a75fbf"
    },
    {
      "value": "susan@example.com",
      "$ref": "/v2/users/h478j87d-63f8-4788-5282-tdc853a75fbk"
    }
  ]
}
```

Response, status **201**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:Group"
    ],
    "externalId": "44c5d28b-f70a-42f5-kjhg-882326f50b83",
    "displayName": "Developers",
    "members": [
        {
            "value": "peter@example.com",
            "$ref": "/v2/users/a549a5cf-63f8-4788-adab-ddc853a75fbf"
        },
        {
            "value": "susan@example.com",
            "$ref": "/v2/users/h478j87d-63f8-4788-5282-tdc853a75fbk"
        }
    ],
    "meta": {
        "resourceType": "Group",
        "location": "/v2/groups/44c5d28b-f70a-42f5-kjhg-882326f50b83"
    }
}
```

Error status codes: **500**, **400**

### Modification

Summary: partially modify the non-disabled existing group that has the given ID. In other words, you may pass in request only those group's attributes and their values that you want to modify. Other attributes are left untouched.

Endpoint: /groups/`id`

Method: **PATCH**

Auth header: according to the chosen auth type 

Example: PATCH /groups/44c5d28b-f70a-42f5-kjhg-882326f50b83

Request body:

```json
{
  "members": [
    {
      "value": "simon@example.com",
      "$ref": "/v2/users/h478j87d-63f8-4788-5282-tdc853a75fbk"
    }
  ]
}
```

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:Group"
    ],
    "externalId": "44c5d28b-f70a-42f5-kjhg-882326f50b83",
    "displayName": "Developers",
    "members": [
        {
            "value": "simon@example.com",
            "$ref": "/v2/users/h478j87d-63f8-4788-5282-tdc853a75fbk"
        }
    ],
    "meta": {
        "resourceType": "Group",
        "location": "/v2/groups/44c5d28b-f70a-42f5-kjhg-882326f50b83"
    }
}
```

Error status codes: **500**, **404**, **400**

### Replacement

Summary: fully replace attributes of the non-blocked existing group that has the given ID. In other words, even if you didn't add some attributes and their values to the request, they will be changed (set to their default values).

Endpoint: /groups/`id`

Method: **PUT**

Auth header: according to the chosen auth type 

Example: PUT /groups/44c5d28b-f70a-42f5-kjhg-882326f50b83

Request body:

```json
{
  "displayName": "Developers"
}
```

Response, status **200**

```json
{
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:Group"
    ],
    "externalId": "44c5d28b-f70a-42f5-kjhg-882326f50b83",
    "displayName": "Developers",
    "meta": {
        "resourceType": "Group",
        "location": "/v2/groups/44c5d28b-f70a-42f5-kjhg-882326f50b83"
    }
}
```

Error status codes: **500**, **404**, **400**

### Disabling

Summary: disable group that exists on our Service Provider side and has the given ID.

Endpoint: /groups

Method: **DELETE**

Auth header: according to the chosen auth type 

Example: DELETE /groups/44c5d28b-f70a-42f5-kjhg-882326f50b83

Response, status **204**

Error status codes: **500**, **404**

## Bulk operations

Not implemented yet

## Searching with filter

Not implemented yet
