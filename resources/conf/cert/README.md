**only use dev mode**
```
openssl x509 -req -days 3650 -in certificate_request.csr -signkey private_key.pem -out certificate.crt
```
**in production, use yourself security certificate file to protect your server.**