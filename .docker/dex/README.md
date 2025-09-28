# Dex

Dex is OIDC provider designed for testing/integration.

## Configuration

The app should be configured with the following parameters:
- Client ID: dex
- Client Secret: dex-secret
- Issuer URL: http://localhost:5556/dex
- The OAuth2 redirect URL should be set to: http://localhost:8080/callback

The OIDC discovery endpoint is:

http://localhost:5556/dex/.well-known/openid-configuration

## Static Users

- User 1:
  - Username: user
  - Password: password
  - Email: user@example.com
  - UserID: 1234

## Generation new passwords

Dex requires passwords to be hashed with bcrypt, not plain text.

You can generate your via cli or using online tools like https://bcrypt-generator.com/

```bash
htpasswd -bnBC 10 "" password | tr -d ':\n'
```
