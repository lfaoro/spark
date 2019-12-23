# 🔥 Fireblaze Vault

> Fireblaze Vault is a tokenization service, aims to be an open platform designed to protect your sensitive data and inherit best-in-class security posture in order to fast-track certifications like PCI DSS, SOC2, HIPAA and others.
>
>Fireblaze Vault helps with tokenization and secure storage of sensitive data, and digital assets like [PII](https://en.wikipedia.org/wiki/Personal_data), [Credit Cards](https://en.wikipedia.org/wiki/Credit_card), Passports/IDs, Credentials, and more.

[![BSD License](https://img.shields.io/badge/license-BSD-blue.svg?style=flat)](LICENSE)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Flfaoro%2Fflares.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Flfaoro%2Fspark?ref=badge_shield)
[![Go Report Card](https://goreportcard.com/badge/github.com/lfaoro/spark)](https://goreportcard.com/report/github.com/lfaoro/spark)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](code-of-conduct.md)

## Insights

- vaulting of payment card data (avoid liability and being locked-in to a payment provider)
- analytics on card scheme, brand, type, currency, banks
- risk assessment based on geolocation, ip address, black lists
- 1-click payment solution, driving impulsive sales up by 55%, removing the barrier of card details re-entry
- automated AML checks on passports/IDs
- GDPR compliant personal identifiable information (PII) storage

### Store a payment card
#### Request
```shell script
curl -X POST \
  http://localhost:3000/v1/card \
  -H 'Content-Type: application/json' \
  -d '{
    "holder": "leonardo", # Cardholder name
    "number": "4415281263901560", # Payment card number
    "exp_month": 1, # Expiry month
    "exp_year": 2022, # Expiry year
    "cvc": 123, # MC(Card Verification Code), VISA(Card Verification Value)
    "auto_delete": "THREE_MONTHS" # Delete this data in 3 months
}'
```
#### Response
```json
{
  "auto_delete_on": "2020-06-27T07:08:31.500606Z",
  "expires_on": "2022-02-01T00:00:00.000000001Z",
  "first_six": 466945,
  "hash": "ZmJpZC0xNDQzNjM1MzE3MzMxNzc2MTQ4V06Nh[...]",
  "last_four": 8424,
  "metadata": {
    "currency": "USD",
    "issuer": {
      "country": "United States of America",
      "country_code": "US",
      "latitude": 38,
      "longitude": -97,
      "map": "https://www.google.com/maps/search/?api=1&query=38,-97"
    },
    "scheme": "visa"
  },
  "mpi": {
    "acs": "https://secure5.arcot.com/acspage/cap?RID=35325&VAA=B",
    "eci": 2,
    "enrolled": true,
    "par": "eNpdU8tymzAU3ecrvMumYz1AgD2yZnDsTpMZ[...]"
  },
  "request_ip": "127.0.0.1",
  "risk": {
    "score": 30
  },
  "token": "tok_e4912b25-b8ef-4cf8-bb0d-449bcaf58e08",
  "user_agent": "grpc-go/1.25.1"
}

```

## Tech stack

We use [protobuf]() to serialize the data and [gRPC](https://grpc.io) to transport it, for compatibility we also support JSON serilization over HTTP transport via reverse-proxy, auto-generated thanks to [grpc-gateway](), which also generates the [Swagger]() documentation, available at https://doc.fireblaze.io/card.

Sensitive data is encrypted at rest using [AES-GCM](https://eprint.iacr.org/2017/168.pdf) and an [HSM](https://en.wikipedia.org/wiki/Hardware_security_module) module to generate entropy for the encryption keys which must be [FIPS 140-2 Level 3](https://nvlpubs.nist.gov/nistpubs/FIPS/NIST.FIPS.140-2.pdf) certified to meet compliance. Check the [kms](pkg/kms) service for the [GCP CloudKMS](https://cloud.google.com/kms/) implementation. Feel free to extend the interface with other implementations e.g. [AWS CloudHSM](https://aws.amazon.com/cloudhsm)

We like to think of data in graphs, leveraging [ent](https://entgo.io/) as our entity framework, which supports PostgreSQL, MySQL, SQLite, Gremlin.

The infrastructure is designed around [Kubernetes](https://kubernetes.io/) with the goal of passing [PCI-DSS](https://www.pcisecuritystandards.org/documents/PCI_DSS_v3-2-1.pdf?agreement=true&time=1573855946115) Level 1 compliance.

The pipelines run on our self-hosted [Gitlab](https://code.fireblaze.io/users/sign_in), feel free to request access, you can sign-in with your Github account.

Fireblaze Vault is currently in [MVP](https://en.wikipedia.org/wiki/Minimum_viable_product) status, we're proud to solve this challenge and excited to share it with the community.

### Technical features

- compliant tokenization of digital assets
- payment card validation w/ regex & luhn check
- payment card metadata retrieval
- payment card risk probability
