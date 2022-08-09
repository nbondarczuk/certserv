# This is a simple CA server implementation

## Description of the functionality

The CA server creates and returns CA Certificates. They are configured
with a json config file. The result Json structure provides the public
and private keys. Other return values are the creation time parameters
which are encoded in the certificate. 

The certificates are stored in a persistent storage which can be:

- in-memory cache
- file ssytem
- MongoDB
- ETCD

Another thing store in the persistent storage is the certificate seq no.
This value is incremented every time the server issues a new certificate.

## API description

The API access points are the following:

- POST /api/certificate/{address}

Creates new certificate for a given address.

- GET /api/certificate

Obtains all the create certificated from the server.

- GET /api/certificate/{address} 

Obtains existing certificate. the certificate must exist in the CA database.

- PATCH /api/certificate/{address}/{duration}/{unit:dwy}

Extends the validity period of the existing certificate by the number of units
starting with the certificate existing validity start date.

- DELETE /api/certificate/{address}

Removes the certificate.

