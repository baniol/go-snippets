# Go jwt

Asymmetric encryption, which means that because a token is signed, the receiver only needs the public key of the signer to validate that the token has indeed come from a trusted source and this allows us to lock down access to the private keys to an authorization server.

base64

headers:

```json
{ 
  "alg": "RS256", 
  "typ": "JWT" 
}
```

The second object payload which contains the details of the claims related to the token:

```json
{ 
  "userID": "abcsd232fjfj", 
  "accessLevel": "user" 
}
```

And finally, the third part is the signature, which is an optional element

Base64URL(header).Base64URL(payload)

The format of the signature can either be symmetrical (HS256) using a shared secret or asymmetrical (RS256), which uses public and private keys. For JWTs, the best option is the asymmetrical option as for a service which needs to authenticate the JWT, it only requires the public part of the key.

We can validate our JWT only using the command line. First, we need to convert our base64URL-encoded signature into standard base64 encoding by replacing _ with / and - with +. We can then pipe that into the base64 command-line application and pass in the -D flag to decode the input; we then output this into a file:

```
cat signature.txt | sed -e 's/_/\//g' -e 's/-/+/g' | base64 -D > signature.sha256 
```