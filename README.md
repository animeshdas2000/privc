# Privc - Data Privacy Vault

Privc tokenizes sensitive values (card numbers, emails, IDs) and turns tokens back into values for authorized callers.

Tokens are AES-GCM ciphertexts with a random nonce, so the same input gives a different token each time and any tampered token is rejected.

## Getting Started

1. Clone the repository:

```
git clone https://github.com/animeshdas2000/privc.git
cd privc
```

2. Set configuration as environment variables (or in a `.env` file, which real environment variables override):

```
# AES key: 16, 24 or 32 bytes. 32 bytes gives AES-256.
ENCRYPTION_KEY=$(openssl rand -hex 16)
# Comma-separated API keys, each at least 32 characters.
API_KEYS=$(openssl rand -hex 32)
```

The server refuses to start if either is missing or invalid.

3. Build and run:

```
go build .
./privc
```

## API

All endpoints except `GET /` require `Authorization: Bearer <API key>`.

```
curl -X POST localhost:8080/tokenize \
  -H "Authorization: Bearer $API_KEY" \
  -d '{"data":{"field1":"4111111111111111"}}'

curl -X POST localhost:8080/detokenize \
  -H "Authorization: Bearer $API_KEY" \
  -d '{"data":{"field1":"<token from tokenize>"}}'
```

## Tests

```
go test ./...
```
