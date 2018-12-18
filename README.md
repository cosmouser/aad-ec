# aad-ec
## version 0.1

AzureAD Entitlement Checker connects to AzureAD and provides a request endpoint and a graphical web interface for interacting with the directory. The tool acquires authorization tokens in order to send authorized requests to Azure.

## usage
```
Usage of ./aad-ec:
  -config string
        path to config.json (default "./config.json")
  -port int
        port to listen on (default 8080)
```

## routes
```
/               # Index
/callback       # RedirectURI for AAD
/login          # Provides URL for delegated permissions sign in
/ece/getPlans   # API endpoint for user lookups
```

## example request URI
```
/ece/getPlans?uid=principal@uni.edu&version=0.1
```