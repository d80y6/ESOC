from fastapi import Request, HTTPException
from jose import jwt, JWTError
import requests
import os

class AuthMiddleware:
    def __init__(self, issuer_url: str, audience: str):
        self.issuer_url = issuer_url
        self.audience = audience
        self.jwks = self._get_jwks()

    def _get_jwks(self):
        try:
            resp = requests.get(f"{self.issuer_url}/protocol/openid-connect/certs")
            return resp.json()
        except Exception:
            return None

    async def __call__(self, request: Request, call_next):
        if request.url.path in ["/health", "/docs", "/openapi.json"]:
            return await call_next(request)

        auth_header = request.headers.get("Authorization")
        if not auth_header or not auth_header.startswith("Bearer "):
            raise HTTPException(status_code=401, detail="Missing or invalid token")

        token = auth_header.split(" ")[1]
        try:
            # In a real scenario, we should validate against JWKS.
            # For this remediation, we'll implement a robust check.
            payload = jwt.decode(token, self.jwks, algorithms=["RS256"], audience=self.audience, issuer=self.issuer_url)

            tenant_id = payload.get("tenant_id")
            if not tenant_id:
                raise HTTPException(status_code=401, detail="Missing tenant_id in token")

            request.state.tenant_id = tenant_id
            request.state.user = payload
        except JWTError as e:
            raise HTTPException(status_code=401, detail=f"Invalid token: {str(e)}")

        return await call_next(request)
