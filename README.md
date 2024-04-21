# Privc - Data Privacy Vault

Privc is a secure and efficient solution for managing and storing sensitive data.

## Getting Started

To get started with Privacy Vault, follow these steps:

1. Installation: Clone the repository to your local machine and navigate to the project directory.

```
git clone https://github.com/your_username/privacy-vault.git
cd privacy-vault
```

2. You'll need to setup `.env` with the following keys:

```
ENCRYPTION_KEY = 'rand32digitEncryptionKey12345678'
IV = 'my16digitIvKey12'
```

3. Build and Run: Build the application and start the server.

```
go build .
./privc
```

4. Access the Application: Open your web browser and navigate to http://localhost:8080 to access the Privacy Vault application.
