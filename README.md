# Git ID (gid) <!-- omit in toc -->

**Git ID** or **gid** is a simple, decentralized digital identity framework that enables users to publicly link cryptographic keys to their GitHub profiles using standard Git repositories.

## Table of Contents <!-- omit in toc -->

- [Overview](#overview)
- [Create Your Git ID](#create-your-git-id)
- [Update Your Git ID](#update-your-git-id)
- [Verify a Git ID](#verify-a-git-id)
- [Security Considerations](#security-considerations)
- [Privacy Considerations](#privacy-considerations)
- [Dependency Considerations](#dependency-considerations)
- [Roadmap](#roadmap)
- [License](#license)
- [Contributing](#contributing)
- [Notes](#notes)

## Overview

Git ID binds public keys to GitHub accounts, allowing third parties to verify control over a GitHub identity cryptographically.

## Create Your Git ID

To create your Git ID:

1. **Create a Public Repository**  
   - Name your repository exactly `gid` on GitHub.

2. **Generate an Ed25519 Private key**  
   Use a tool such as `openssl`:

   ```bash
   openssl genpkey -algorithm ED25519 -out gid_private.pem
   ```

   This command creates a private key and saves it to `gid_private.pem`. Store this file securely — it is your private key and must be kept confidential. You will need this key later to sign data and generate cryptographic proofs.

3. **Export the Public Key in PEM Format**  

   Derive the public key from the private key and save it in PEM format:

   ```bash
   openssl pkey -in gid_private.pem -pubout -out gid.pem
   ```
   This command reads your private key from gid_private.pem and writes the corresponding public key to gid.pem.


4. **Commit and Push Your Public Key**
   - Place `gid.pem` at the root of your `gid` repository.
   - Commit and push the file:

   ```bash
   git add gid.pem
   git commit -m "Add initial Git ID public key"
   git push
   ```

## Update Your Git ID

To update your Git ID:

- Generate a new Ed25519 keypair if needed.
- Replace `gid.pem` in your repository with the new public key.
- Commit and push the changes.

## Verify a Git ID

To verify a Git ID:

1. **Receive GitHub handler and a signature**
   - Receive a signed or encrypted file or a response to a cryptographic challenge along with user's GitHub handler.

2. **Fetch the Public Key**
   - Retrieve `gid.pem` from the user's `gid` repository

3. **Perform Verification**
   - Verify that the public key verifies the cryptographic proof

Successful verification confirms that the user controlled the GitHub account when the public key was published and had not revoked or replaced either the public key or its associated private key at the time the cryptographic proof was generated.

## Security Considerations

- **GitHub Account Security:**  
  Enable strong passwords and two-factor authentication (2FA).
- **Key Compromise:**  
  If your GitHub account or private key is compromised, immediately replace `gid.pem` with a new key and notify users if possible.
- **Tampering Risk:**  
  Since repositories are mutable, verifiers should optionally cache public key fingerprints and monitor for unauthorized changes.

## Privacy Considerations

- **Public Identity Binding:**  
  Git ID publicly binds a cryptographic public key to your GitHub username. This link may will expose your identity across platforms that use your key.
- **Metadata Leakage:**  
  Verifiers or observers can correlate activities tied to your Git ID key to analyze your network, behaviour, or affiliations.

## Dependency Considerations

- Git ID relies on the **availability and security of GitHub** for hosting the public key.
- While Git ID uses Git for distribution, it depends on GitHub's infrastructure and identification and authentication processes.
- Future extensions may support other Git hosting providers (e.g., GitLab, Bitbucket) and self-hosted or decentralized Git networks.

## Roadmap

- **Multi-Key Support:** Allow users to publish multiple keys (e.g., for different purposes).
- **Multi-Key-Type Support:** Add support for RSA, ECC, and Post-Quantum Cryptography (PQC) keys.
- **JWK Format Investigation:** Evaluate using JSON Web Keys (JWK) for enhanced interoperability.
- **Key Purpose Restrictions:** Enable users to specify intended key usage (signing, encryption, authentication).
- **Historical Key Resolution:** Allow retrieval of old public keys and support for key rotation history.
- **Compromise Disclosure Mechanism:** Establish standards for publicly disclosing key compromises.
- **Organisational Identity and Public Key** Define the rules for organisational keys with multi-layered structure.

## License

This project specification is dedicated to the public domain under the [CC0 1.0 Universal](https://creativecommons.org/publicdomain/zero/1.0/) license.  
You are free to copy, modify, distribute, and use it without restriction.

This project code is licensed under the [MIT License](cli/LICENSE).

## Contributing

Contributions to **Git ID** are welcome!

Ways to contribute:

1. **Propose Ideas:**  
   - Open an issue for improvements, clarifications, or new features.

2. **Report Issues:**  
   - Report errors, ambiguities, or security concerns via GitHub Issues.

3. **Submit Pull Requests:**  
   - Before you start working on a pull request, make sure that the capability or improvement has been agreed upon in the GitHub issue.
   - Fork the repository, create a new branch, make clean commits, and open a pull request with a clear description.
   - A single pull request MUST resolve a single open issue.

4. **Follow Project Principles:**  
   - **Simplicity First:** Keep implementations minimal and understandable.
   - **Security and Privacy:** Prioritize the safety and discretion of users.
   - **Backward Compatibility:** Maintain compatibility when possible.

By contributing, you agree that your submissions are licensed under CC0 1.0 Universal.

## Notes

- Git ID currently supports only **GitHub**.
- The model is **Git host-agnostic** and can be extended to other platforms (e.g., GitLab, Bitbucket, self-hosted Git servers, decentralized Git systems).
