# Git ID (gid) <!-- omit in toc -->

**Git ID** (`gid`) is an open, decentralized identity framework for humans and
autonomous agents. It binds cryptographic key material to stable, publicly
trusted identity namespaces — GitHub accounts and owned web domains — using
standard JWKS (JSON Web Key Sets), without requiring a central registry or
certificate authority.

## Table of Contents

- [Table of Contents](#table-of-contents)
- [Overview](#overview)
- [Design Principles](#design-principles)
- [Identity Providers](#identity-providers)
- [Key Publication Format: JWKS](#key-publication-format-jwks)
  - [`kid` Rules](#kid-rules)
  - [Key Expiry](#key-expiry)
  - [Minimal example `jwks.json`](#minimal-example-jwksjson)
- [Key Generation](#key-generation)
  - [Generate an Ed25519 Signing Key Pair](#generate-an-ed25519-signing-key-pair)
  - [Generate an X25519 Encryption Key Pair](#generate-an-x25519-encryption-key-pair)
- [GitHub-Based Identity](#github-based-identity)
  - [Setup](#setup)
  - [Repository Layout](#repository-layout)
- [Domain-Based Identity](#domain-based-identity)
  - [Single-Tenant Layout](#single-tenant-layout)
- [Multi-Tenant Domain Hosting](#multi-tenant-domain-hosting)
  - [Multi-Tenant Layout](#multi-tenant-layout)
- [Layout Discovery](#layout-discovery)
  - [Discovery Document](#discovery-document)
  - [Resolution Algorithm](#resolution-algorithm)
  - [Layout Selection Guide](#layout-selection-guide)
- [Agent Identity](#agent-identity)
  - [Agent Identifiers](#agent-identifiers)
  - [Registering an Agent (GitHub)](#registering-an-agent-github)
  - [Registering an Agent (Domain)](#registering-an-agent-domain)
  - [Owner Accountability](#owner-accountability)
- [Addressing and Resolution](#addressing-and-resolution)
  - [Address Formats](#address-formats)
  - [Resolution Algorithm](#resolution-algorithm-1)
- [Verification](#verification)
  - [Verifying an Owner](#verifying-an-owner)
  - [Verifying an Agent](#verifying-an-agent)
  - [Re-fetch Policy](#re-fetch-policy)
  - [Trust Chain Summary](#trust-chain-summary)
- [Messaging: Encryption and Authentication](#messaging-encryption-and-authentication)
  - [Signing (Authentication)](#signing-authentication)
  - [Encryption](#encryption)
- [Key Rotation, Expiry, and Revocation](#key-rotation-expiry-and-revocation)
  - [Rotating a Key](#rotating-a-key)
  - [Key Expiry](#key-expiry-1)
  - [Revoking an Agent](#revoking-an-agent)
  - [Rotation vs. Compromise](#rotation-vs-compromise)
- [Security Considerations](#security-considerations)
- [Privacy Considerations](#privacy-considerations)
- [Trust Model and Limitations](#trust-model-and-limitations)
- [Roadmap](#roadmap)
- [License](#license)
- [Contributing](#contributing)

## Overview

In the agentic era, identity is not only human. Software agents — AI assistants,
automated pipelines, bots, CI workers — increasingly act on behalf of people and
organisations: sending messages, signing artifacts, making decisions, calling
external services. These agents need stable, verifiable identities that are
anchored to a human or organisational owner.

`gid` provides that anchor. It defines:

- Where to publish cryptographic keys for a human or organisation (the
  **owner**)
- Where to publish cryptographic keys for autonomous **agents** acting on behalf
  of that owner
- How any third party can resolve, verify, and trust those keys
- How agents are addressed in a compact, human-readable form

The key insight is that **existing public identity providers already perform
identity verification** — GitHub, for example, requires email confirmation,
supports 2FA and passkeys, and is widely trusted in professional and developer
contexts. `gid` does not replace or duplicate that trust: it extends it by
attaching cryptographic key material to the identities those providers already
establish.

## Design Principles

- **Build on existing trust.** Use identity providers that already have public
  trust, rather than creating a new registry.
- **Standard formats.** Publish keys as JWKS (RFC 7517) for interoperability
  with existing OAuth, OIDC, and cryptographic tooling.
- **Minimal infrastructure.** No additional infrastructure. Git repository or a
  static web server is sufficient.
- **Explicit agent ownership.** Agents are registered under their owner's
  identity. The owner is accountable for the agents they register.

## Identity Providers

`gid` currently defines two identity provider (IDP) models. Both use JWKS as the
key publication format and share the same agent registration layout.

| Model          | Trust Anchor         | Owner controls                           |
| -------------- | -------------------- | ---------------------------------------- |
| **GitHub**     | GitHub account       | A public repository named `gid`          |
| **Web Domain** | DNS + TLS (CA trust) | `/.well-known/` paths on an owned domain |

Both models inherit the strengths and limitations of their underlying trust
anchor. A GitHub-based identity is as trustworthy as the GitHub account; a
domain-based identity is as trustworthy as control of the domain and its TLS
certificate.

Other Git hosting platforms (GitLab, Gitea, self-hosted) follow the same
repository layout as GitHub and may be added as formally defined IDPs in future
revisions.

## Key Publication Format: JWKS

All keys — for owners and agents — are published as **JWKS** (JSON Web Key Sets,
[RFC 7517](https://www.rfc-editor.org/rfc/rfc7517)).

A `jwks.json` file contains a `keys` array. Each key object carries:

- `kty`: key type (`OKP` for Ed25519/X25519, `EC` for ECDSA/ECDH, `RSA`)
- `use`: intended use — `sig` for signing/authentication, `enc` for encryption
- `alg`: specific algorithm (e.g. `EdDSA`, `ECDH-ES`)
- `kid`: a stable, unique key identifier used to select keys across messages and during rotation
- `crv`: curve name for `OKP` and `EC` keys (e.g. `Ed25519`, `X25519`)
- `x`: the public key material, encoded as base64url (URL-safe base64 without
  padding) of the raw public key bytes
- Public key material only — never private keys

Publishing separate keys with explicit `use` values is required. A signing key
(`Ed25519`, `use: sig`) and an encryption key (`X25519`, `use: enc`) are
distinct and must not be conflated. A single key must not be assigned `use: sig`
and `use: enc` simultaneously.

### `kid` Rules

- `kid` values must be unique within a `jwks.json` and must never be reused
  across rotations. A rotated key receives a new, distinct `kid`.
- During a rotation transition, both the old and new keys may coexist in
  `jwks.json` under their distinct `kid` values. The old key should be removed
  once in-flight messages using it are no longer expected.
- Verifiers that encounter a `kid` not present in the current `jwks.json` must
  treat the key as revoked and reject the message.

### Key Expiry

Each key object may carry an `exp` field containing a Unix timestamp (seconds
since epoch) after which the key must not be accepted by verifiers:

```json
"exp": 1893456000
```

Verifiers must reject keys whose `exp` has passed, regardless of their presence
in `jwks.json`. Owners should rotate keys before their expiry. Omitting `exp`
means the key has no declared expiry; verifiers in high-assurance contexts
should treat such keys with appropriate caution.

### Minimal example `jwks.json`

```json
{
  "keys": [
    {
      "kty": "OKP",
      "crv": "Ed25519",
      "use": "sig",
      "alg": "EdDSA",
      "kid": "owner-sig-1",
      "x": "<base64url-encoded-public-key>",
      "exp": 1893456000
    },
    {
      "kty": "OKP",
      "crv": "X25519",
      "use": "enc",
      "alg": "ECDH-ES",
      "kid": "owner-enc-1",
      "x": "<base64url-encoded-public-key>",
      "exp": 1893456000
    }
  ]
}
```

The `x` field contains the raw 32-byte public key, base64url-encoded without
padding.

## Key Generation

This section shows how to generate the required key pairs and format them for
`jwks.json`. A reference CLI (see `cli/`) automates this process; the manual
steps below are provided for auditability and environments where the CLI is
unavailable.

### Generate an Ed25519 Signing Key Pair

```bash
# Generate private key
openssl genpkey -algorithm ED25519 -out signing_private.pem

# Derive public key
openssl pkey -in signing_private.pem -pubout -out signing_public.pem

# Extract raw public key bytes as base64url (for the JWKS x field)
openssl pkey -in signing_public.pem -pubin -outform DER \
  | tail -c 32 \
  | base64 | tr '+/' '-_' | tr -d '='
```

### Generate an X25519 Encryption Key Pair

```bash
# Generate private key
openssl genpkey -algorithm X25519 -out encryption_private.pem

# Derive public key
openssl pkey -in encryption_private.pem -pubout -out encryption_public.pem

# Extract raw public key bytes as base64url (for the JWKS x field)
openssl pkey -in encryption_public.pem -pubin -outform DER \
  | tail -c 32 \
  | base64 | tr '+/' '-_' | tr -d '='
```

Store all private key files (`*_private.pem`) securely. They must never be
committed, shared, or published. For agent keys, the private key is held
exclusively by the agent's runtime environment.

## GitHub-Based Identity

### Setup

1. Create a public repository named exactly `gid` on GitHub.
2. Generate a signing key pair (Ed25519) and an encryption key pair (X25519) following the [Key Generation](#key-generation) section.
3. Publish both public keys as `jwks.json` at the root of the repository.
4. Commit and push.

The owner's JWKS is resolvable at:

```
https://raw.githubusercontent.com/<username>/gid/main/jwks.json
```

### Repository Layout

```
gid/                                        # GitHub repo: github.com/<username>/gid
├── jwks.json                               # Owner's public keyset
└── agents/
    ├── <agent-id>/
    │   └── jwks.json                       # Agent's public keyset
    └── <another-agent-id>/
        └── jwks.json
```

## Domain-Based Identity

An owner who controls a web domain publishes their JWKS at a standardised path
under `/.well-known/`. HTTPS is required; plain HTTP is not acceptable.

This follows the same pattern used by OAuth 2.0 and OIDC providers for their own
public keys, and requires only static file hosting with HTTPS.

### Single-Tenant Layout

Use this layout when the domain represents a single identity (personal domain, single-organisation domain).

```
/.well-known/
├── jwks.json                               # Owner's public keyset
└── agents/
    ├── <agent-id>/
    │   └── jwks.json
    └── <another-agent-id>/
        └── jwks.json
```

Owner JWKS URL: `https://<domain>/.well-known/jwks.json`  
Agent JWKS URL: `https://<domain>/.well-known/agents/<agent-id>/jwks.json`

## Multi-Tenant Domain Hosting

A single domain may host `gid` identities for multiple users. This is
appropriate for hosted platforms, developer tools managing keys on behalf of
users, or organisations with per-member identities.

Multi-tenant paths are prefixed with `/.well-known/gid/<username>/` to namespace
user identities without conflicting with other `/.well-known/` standards or with
each other.

### Multi-Tenant Layout

```
/.well-known/gid/
├── layout.json                             # Layout discovery document (see below)
├── <username>/
│   ├── jwks.json                           # User's public keyset
│   └── agents/
│       ├── <agent-id>/
│       │   └── jwks.json
│       └── <another-agent-id>/
│           └── jwks.json
└── <another-username>/
    ├── jwks.json
    └── agents/
        └── ...
```

Owner JWKS URL: `https://<domain>/.well-known/gid/<username>/jwks.json`  
Agent JWKS URL: `https://<domain>/.well-known/gid/<username>/agents/<agent-id>/jwks.json`

## Layout Discovery

A verifier given a domain address must determine whether the domain uses
single-tenant or multi-tenant layout before resolving a JWKS URL. This is done
via a **layout discovery document**.

### Discovery Document

A domain hosting `gid` identities must publish a discovery document at:

```
https://<domain>/.well-known/gid/layout.json
```

The document declares the layout type in use:

```json
{
  "version": "1",
  "layout": "single"
}
```

or

```json
{
  "version": "1",
  "layout": "multi"
}
```

### Resolution Algorithm

Given a domain-based address, a verifier:

1. Fetches `https://<domain>/.well-known/gid/layout.json` over HTTPS.
2. Reads the `layout` field.
3. If `layout` is `"single"`: resolves JWKS URLs using the single-tenant paths.
4. If `layout` is `"multi"`: resolves JWKS URLs using the multi-tenant paths.
5. If the discovery document is absent or malformed: the verifier must not attempt to guess the layout. It should report the identity as unresolvable.

### Layout Selection Guide

| Scenario                          | Layout                                   | Discovery document location                            |
| --------------------------------- | ---------------------------------------- | ------------------------------------------------------ |
| Personal domain (`alice.com`)     | `single`                                 | `https://alice.com/.well-known/gid/layout.json`        |
| Organisation domain, one identity | `single`                                 | `https://acme.org/.well-known/gid/layout.json`         |
| Platform hosting multiple users   | `multi`                                  | `https://platform.example/.well-known/gid/layout.json` |
| GitHub account                    | N/A — GitHub IDP uses its own resolution | N/A                                                    |

## Agent Identity

An **agent** is any autonomous process — an AI model instance, a software
assistant, a CI pipeline, a signing service — that acts on behalf of its owner.
Each agent has its own key pair, a stable identifier assigned by the owner, and
its public keys published under the owner's identity namespace.

### Agent Identifiers

An agent identifier (`agent-id`) is a lowercase slug: letters, digits, and hyphens. It must be:

- Unique within the owner's agent registry
- Stable across key rotations — the identifier refers to the agent role, not a specific key
- Chosen by the owner to reflect the agent's purpose (e.g. `assistant`, `ci-signer`, `research-agent-1`) or kept opaque (e.g. `agent-7f3a`) for privacy

### Registering an Agent (GitHub)

1. Generate a signing keypair (Ed25519) and optionally an encryption keypair
   (X25519) for the agent, following the [Key Generation](#key-generation)
   section.
2. Compose a `jwks.json` containing the agent's public keys.
3. Place the file at `agents/<agent-id>/jwks.json` in the owner's `gid` repository.
4. Commit and push.

```bash
mkdir -p agents/<agent-id>
cp agent_jwks.json agents/<agent-id>/jwks.json
git add agents/<agent-id>/jwks.json
git commit -m "Register agent: <agent-id>"
git push
```

The agent's private keys are held exclusively by the agent's runtime
environment. The owner does not need to retain them.

### Registering an Agent (Domain)

Place the agent's `jwks.json` at:

```
# Single-tenant
https://<domain>/.well-known/agents/<agent-id>/jwks.json

# Multi-tenant
https://<domain>/.well-known/gid/<username>/agents/<agent-id>/jwks.json
```

### Owner Accountability

Publishing an agent's keys under an owner's namespace is an explicit act of
endorsement. It asserts:

- The owner authorises this agent to act on their behalf.
- Messages signed by the agent's key are attributable to the owner's identity.
- The owner is responsible for revoking the agent's keys if compromised or if
  the agent is decommissioned.

## Addressing and Resolution

### Address Formats

Addresses are compact, human-readable identifiers. They are not URIs but map
deterministically to JWKS URLs via the resolution algorithm defined below.

**Owner addresses:**

| IDP                    | Address format        | Example                  |
| ---------------------- | --------------------- | ------------------------ |
| GitHub                 | `github:<username>`   | `github:alice`           |
| Domain (single-tenant) | `<domain>`            | `alice.com`              |
| Domain (multi-tenant)  | `<domain>/<username>` | `platform.example/alice` |

**Agent addresses:**

| IDP                    | Address format                   | Example                                   |
| ---------------------- | -------------------------------- | ----------------------------------------- |
| GitHub                 | `github:<username>/<agent-id>`   | `github:alice/research-agent-1`           |
| Domain (single-tenant) | `<domain>/<agent-id>`            | `alice.com/research-agent-1`              |
| Domain (multi-tenant)  | `<domain>/<username>/<agent-id>` | `platform.example/alice/research-agent-1` |

The `github:` prefix is a scheme token that identifies the GitHub IDP. For
domain-based addresses, the IDP is implied by the absence of a scheme token and
resolved via the layout discovery document.

### Resolution Algorithm

**GitHub IDP (`github:<username>`):**

| Target             | JWKS URL                                                                            |
| ------------------ | ----------------------------------------------------------------------------------- |
| Owner              | `https://raw.githubusercontent.com/<username>/gid/main/jwks.json`                   |
| Agent `<agent-id>` | `https://raw.githubusercontent.com/<username>/gid/main/agents/<agent-id>/jwks.json` |

No discovery step is required for GitHub-based addresses.

**Domain IDP:**

1. Parse the address to extract `<domain>`, and if present, `<username>` and `<agent-id>`.
2. Fetch the discovery document at `https://<domain>/.well-known/gid/layout.json`.
3. Read the `layout` field and apply the corresponding URL pattern:

| layout   | Target | JWKS URL                                                                  |
| -------- | ------ | ------------------------------------------------------------------------- |
| `single` | Owner  | `https://<domain>/.well-known/jwks.json`                                  |
| `single` | Agent  | `https://<domain>/.well-known/agents/<agent-id>/jwks.json`                |
| `multi`  | Owner  | `https://<domain>/.well-known/gid/<username>/jwks.json`                   |
| `multi`  | Agent  | `https://<domain>/.well-known/gid/<username>/agents/<agent-id>/jwks.json` |

4. Fetch the resolved JWKS URL over HTTPS.
5. If the discovery document is absent or malformed, treat the identity as unresolvable.

## Verification

### Verifying an Owner

1. Obtain the owner's address and parse the IDP type.
2. Resolve the JWKS URL using the [Resolution Algorithm](#resolution-algorithm).
3. Fetch `jwks.json` over HTTPS.
4. Select the key matching the `kid` referenced in the message or challenge.
5. Reject the key if its `exp` has passed or if the `kid` is not present in the current `jwks.json`.
6. Verify the signature (for `use: sig` keys) or decrypt (for `use: enc` keys).

### Verifying an Agent

1. Obtain the agent address and parse the owner and agent-id components.
2. Resolve the agent's JWKS URL using the [Resolution Algorithm](#resolution-algorithm).
3. Fetch the agent's `jwks.json` over HTTPS.
4. Select the key matching the `kid` in the message.
5. Reject if `exp` has passed or `kid` is absent from the current `jwks.json`.
6. Verify the signature or decrypt.

### Re-fetch Policy

Verifiers must not rely solely on cached JWKS. Recommended re-fetch intervals:

- **Default:** re-fetch every 24 hours, or at the start of each new session
- **High-assurance:** re-fetch on every verification
- **Short-lived sessions:** a single fetch at session start is acceptable for the session duration

### Trust Chain Summary

```
CA (TLS)
  └── HTTPS delivery of jwks.json
        └── IDP namespace (GitHub account / domain ownership)
              └── Owner's JWKS
                    └── Agent's JWKS (endorsed by owner via publication)
                          └── Signed or encrypted message
```

## Messaging: Encryption and Authentication

### Signing (Authentication)

Use the key with `"use": "sig"` (Ed25519 / `EdDSA`).

An agent signs outgoing messages with its private signing key. Recipients
resolve the agent address, fetch the JWKS, select the key by `kid`, and verify
the signature. No central broker is required.

### Encryption

Use the key with `"use": "enc"` (X25519 / `ECDH-ES`).

To send an encrypted message to an agent:

1. Fetch the agent's JWKS and select the key with `"use": "enc"`.
2. Perform ECDH-ES key agreement using the agent's X25519 public key and an ephemeral sender key pair.
3. Encrypt the message payload with the derived key (e.g. AES-256-GCM).
4. Only the agent's runtime, holding the corresponding X25519 private key, can decrypt it.

Signing and encryption keys are always separate. Ed25519 keys are for signatures
only; X25519 keys are for key agreement and encryption only. JWKS `use` fields
make this explicit and machine-readable.

## Key Rotation, Expiry, and Revocation

### Rotating a Key

1. Generate a new key pair (see [Key Generation](#key-generation)).
2. Add the new public key to `jwks.json` with a new, unique `kid`.
3. Optionally retain the old key temporarily under its original `kid` to allow
   in-flight messages to be processed.
4. Remove the old key once it is no longer needed.
5. Commit and push (GitHub) or redeploy (domain).

Verifiers that encounter a `kid` absent from the current `jwks.json` must treat
it as revoked.

### Key Expiry

Set an `exp` field on each key at publication time. Verifiers must reject keys
whose `exp` has passed. Owners should rotate keys before expiry to avoid a gap
in verifiable identity. A suggested default key lifetime is one year.

### Revoking an Agent

Remove the `agents/<agent-id>/` directory (GitHub) or the corresponding
`jwks.json` file (domain). Commit and push (GitHub) or redeploy (domain).

Verifiers must treat absence of an agent path as revocation.

### Rotation vs. Compromise

The framework provides no cryptographic mechanism to distinguish planned
rotation from a compromise event. Owners who need to signal compromise
explicitly should use an out-of-band channel (e.g. a signed notice committed to
the repository, a status file) until a formal compromise disclosure mechanism is
standardised (see Roadmap).

## Security Considerations

**IDP account security.** The security of all registered keys depends on the
security of the owner's IDP account. For GitHub: enable 2FA and use a hardware
security key or passkey. For domains: protect DNS registrar access and use
DNSSEC where available.

**HTTPS is required.** JWKS must always be fetched over HTTPS. Verifiers must
reject plain HTTP fetches entirely.

**Mutability.** Both GitHub repositories and web servers are mutable by their
operators. This framework does not provide cryptographic proof that a key has
not been silently replaced. Verifiers in high-assurance contexts should cache
key fingerprints at first use and alert on unexpected changes between
re-fetches.

**Agent key isolation.** Each agent holds only its own private keys. Keys must
not be shared between agents or exposed outside the agent's runtime environment.

**Stale cache risk.** Verifiers that cache JWKS without re-fetching may continue
to accept messages from revoked agents. Follow the [Re-fetch
Policy](#re-fetch-policy) appropriate to the application's security
requirements.

**`kid` collision.** `kid` values must be unique within a `jwks.json`. Verifiers
that encounter duplicate `kid` values within a single `jwks.json` should reject
the entire key set as malformed.

## Privacy Considerations

**Public key binding.** Publishing a `gid` identity publicly binds cryptographic
keys to a GitHub username or domain. This is an intentional, public act.

**Agent pseudonymity.** Owners may choose opaque agent identifiers (e.g.
`agent-7f3a`) rather than descriptive names to limit inference about agent
purpose or the services it interacts with.

**Metadata leakage.** Observers can monitor a `gid` repository or a domain's
`/.well-known/` paths to track agent registration, rotation, and decommissioning
events over time.

## Trust Model and Limitations

`gid` is a **pragmatic, low-ceremony identity layer** for developers and
autonomous agents. It is explicitly not a substitute for PKI in high-assurance
environments.

**What it provides:**

- A standardised, machine-readable way to find and use public keys for a GitHub
  user, domain owner, or their agents
- A chain of accountability from message to agent to human owner
- Interoperability with any tooling that understands JWKS and standard key
  algorithms

**What it does not provide:**

- Cryptographic proof of IDP account ownership beyond HTTPS delivery
- A tamper-proof audit log (GitHub commit history and web server logs are owner-controlled)
- Revocation push notifications or real-time revocation checks
- A cryptographic distinction between routine key rotation and compromise
- Protection against a compromised or compelled IDP

**Appropriate use cases:** developer tooling, agent-to-agent communication, API
authentication, open-source project identity, automated signing pipelines.

**Not appropriate for:** legally binding identity, high-value financial
transactions, contexts where IDP unavailability or compromise would be
catastrophic.

## Roadmap

- **Agent Registry Manifest.** Standardise an optional `agents/index.json`
  listing all registered agents with metadata (description, scope, expiry) to
  enable agent discovery without prior knowledge of agent identifiers.
- **Agent Metadata.** Structured per-agent metadata: authorization scope,
  runtime identifier, model version, capability declarations.
- **Compromise Disclosure.** A standard mechanism for owners to publicly signal
  key compromise, distinct from routine rotation.
- **Historical Key Resolution.** Allow retrieval of past public keys via Git
  history (GitHub model) or versioned archives (domain model).
- **Cross-Owner Delegation.** Define a mechanism for one owner to delegate a
  registered agent to act on behalf of another owner, with scoped permissions.
- **Organisational Identity.** Multi-layer org identity with team-level agent
  registries and delegated sub-registries.
- **Additional IDP Models.** GitLab, Gitea, self-hosted Git servers.
- **Post-Quantum Keys.** Add support for PQC algorithms (e.g. ML-DSA /
  Dilithium) as JWKS key types once tooling matures.

## License

This specification is dedicated to the public domain under the [CC0 1.0
Universal](https://creativecommons.org/publicdomain/zero/1.0/) license. You are
free to copy, modify, distribute, and use it without restriction.

Reference implementation code is licensed under the [MIT License](cli/LICENSE).

## Contributing

Contributions are welcome. Please open an issue before submitting a pull
request. Each pull request must resolve a single open issue.

**Project principles:**

- **Simplicity first.** Keep the spec minimal and implementable.
- **Standard formats.** Prefer existing standards (JWKS, ECDH-ES, EdDSA) over bespoke formats.
- **Honest trust model.** Do not overstate the security properties of any IDP or the framework itself.
- **Backward compatibility.** Maintain compatibility when possible.

By contributing, you agree that your submissions are licensed under CC0 1.0 Universal.
