# Environment Secrets with SOPS and age

Notegic uses SOPS and age to maintain encrypted environment files locally and
on deployment hosts. No secrets server or repository-stored environment
artifact is required for this workflow.

## Repository files

```text
.sops.yaml                 # local/deployment-only SOPS rules
.env.example                         # local-only variable reference, never tracked
.env                                 # development plaintext file
secrets/envs/.env.<environment>      # production/test/staging plaintext files
.env.enc                             # development encrypted file
secrets/envs/.env.<environment>.enc  # production/test/staging encrypted files
```

Encrypted files are transferred directly through an approved private channel
and are never committed to GitHub. Their ciphertext is not a secret, but
plaintext values and age identities are. SOPS may leave dotenv key names and
structure visible, so sensitive information must not be encoded in variable
names.

## New developer onboarding

1. Install SOPS and age.
   On macOS with Homebrew:

   ```sh
   brew install sops age
   ```

2. Generate an age identity locally with `age-keygen`.
3. Send only the generated public recipient (`age1...`) to the repository
   maintainer.
4. The maintainer adds the recipient to the local/deployment `.sops.yaml` and runs:

   ```sh
   make env-updatekeys -e development
   # transfer the encrypted artifact through the private deployment channel
   ```

5. Receive the updated encrypted file through the approved private channel and
   run:

   ```sh
   make env-decrypt -e development
   ```

Private identities stay on the owner’s machine. They must not be committed,
emailed together with the encrypted file, copied into Docker images, or printed
in logs.

## Local commands

```sh
make env-encrypt -e development
make env-decrypt -e development
make env-edit -e development
make env-updatekeys -e development
make env-rotate -e development

# Encrypt all supported environments in one command.
make env-encrypt-all
```

The first encryption requires at least one real age recipient in the local
`.sops.yaml`. The Makefile writes encryption and decryption output to a
temporary file and replaces the destination only after SOPS succeeds, so a
failed operation does not destroy a previously valid artifact.
The plaintext environment files are ignored by Git and are written with mode
`0600` by `env-decrypt`. Development is a special case: its plaintext is `.env`
and its encrypted artifact is `.env.enc`; production, test, and
staging use `secrets/envs/.env.<environment>` and
`secrets/envs/.env.<environment>.enc`. GNU Make reserves `-a`, so use
`env-encrypt-all` instead of `make env-encrypt -a`.

## CI/CD and staging

CI and Jenkins use a dedicated age identity stored in their credential store.
Production and staging identities are separate from developer identities. A
staging host may receive only the encrypted file and keep the private identity
in its credential store; the deployment script decrypts into a temporary file
for Docker Compose and removes it when the command exits:

```sh
IMAGE_REGISTRY=ghcr.io/ORG/REPO \
IMAGE_TAG=TAG \
COMPOSE_ENCRYPTED_ENV_FILE=/etc/notegic/staging.env.enc \
SOPS_CONFIG_FILE=/workspace/notegic-backend/.sops.yaml \
SOPS_AGE_KEY_FILE=/etc/notegic/sops/age/keys-staging.txt \
make staging-deploy
```

When `COMPOSE_ENCRYPTED_ENV_FILE` is used, `SOPS_AGE_KEY_FILE` is required.
The deployment and smoke scripts verify that the file exists before asking
SOPS to decrypt. The private key file is supplied by the deployment host or
credential store and is never copied into the repository or image.

The existing `COMPOSE_ENV_FILE=/etc/notegic/staging.env` path remains supported
for compatibility while encrypted deployment is introduced.

### Deployment host layout

Place each environment's encrypted artifact, SOPS rules, and private age
identity on the corresponding deployment host (or mount them from its secret
store) with restrictive permissions:

```text
/etc/notegic/sops/.sops.yaml
/etc/notegic/sops/age/keys-staging.txt
/etc/notegic/secrets/envs/.env.staging.enc

# Production uses a separate identity and artifact:
/etc/notegic/sops/age/keys-production.txt
/etc/notegic/secrets/envs/.env.production.enc
```

The staging command should point at the first three files:

```sh
COMPOSE_ENCRYPTED_ENV_FILE=/etc/notegic/secrets/envs/.env.staging.enc \
SOPS_CONFIG_FILE=/etc/notegic/sops/.sops.yaml \
SOPS_AGE_KEY_FILE=/etc/notegic/sops/age/keys-staging.txt \
make staging-deploy
```

The private key files should be mode `0600` and readable only by the deploy
user. The `.sops.yaml` file contains public recipients and is not a private
key, but it should still be managed as deployment configuration. Never place
`keys-*.txt` under the repository, `secrets/`, a container image, or a CI log.

## Rotation and removal

When a member or host is removed:

1. Remove its public recipient from `.sops.yaml`.
2. Run `make env-updatekeys -e <environment>`.
3. Run `make env-rotate -e <environment>`.
4. Rotate the actual passwords, tokens, and API credentials if the identity may
   have been compromised.

Git history is an audit and version source, not a revocation mechanism. A
previously authorized identity may still decrypt old commits, so compromised
credentials must be rotated and old access must not be treated as revoked by
history rewriting alone.
