# CorePass and CORE ID access

CORE ID allowlists and a private gateway are not enabled. CORE ID metadata is rejected, and no authentication proof is accepted. A future integration must validate the actual CorePass OIDC tokens against issuer discovery and JWKS, bind the verified CORE ID to a protected object/share record, and keep private content off the public IPFS network. See [CorePass authentication documentation](https://docs.corepass.net/corepass-connector/authentication/).
