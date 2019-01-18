# gobundledhttp

Provides convenience functions to generate an *http.Client with bundled CA certificates.

Available Methods:

```go
// Normal TLS-capable client
tlsClient := gobundledhttp.NewClient()

// Maximum security client (TLS 1.2 only, restricted ciphersuite)
maxTLSClient := gobundledhttp.MaxClient()

// Disables SSL certificate checking
noTLSClient := gobundledhttp.InsecureClient()

// Get the default cert pool to build your own client
certpool := gobundledhttp.GetPool()

// Build an OAuth2 context with trusted pool
ctx := context.WithValue(
        context.Background(),
        oauth2.HTTPClient,
        gobundledhttp.NewClient(),
      )

// Add a PEM byte slice (b) to the pool
success = gonordhttp.AddPEM(b)
```

To update certificates (needed only if the cacerts source changes):

```bash
go generate
```
