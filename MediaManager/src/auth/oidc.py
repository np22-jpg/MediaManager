from fastapi import Depends, APIRouter
from fastapi.security import OpenIdConnect

oidc = OpenIdConnect(openIdConnectUrl="http://localhost:8080/realms/tools/.well-known/openid-configuration")
app = APIRouter()

@app.get("/foo")
async def bar(token = Depends(oidc)):
    return token