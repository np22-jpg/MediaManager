from os import environ

from fastapi import Depends, APIRouter
from fastapi.openapi.models import OAuthFlows, OAuthFlowAuthorizationCode
from fastapi.security import OpenIdConnect, OAuth2AuthorizationCodeBearer
from pydantic import BaseModel

from auth import router


#TODO: Implement OAuth2/Open ID Connect
class Settings(BaseModel):
    OAUTH2_AUTHORIZATION_URL: str
    OAUTH2_TOKEN_URL: str
    OAUTH2_SCOPE: str

    @property
    def oauth2_flows(self) -> OAuthFlows:
        return OAuthFlows(
            authorizationCode=OAuthFlowAuthorizationCode(
                authorizationUrl=self.OAUTH2_AUTHORIZATION_URL,
                tokenUrl=self.OAUTH2_TOKEN_URL,
                scopes={self.OAUTH2_SCOPE: "Access to this API"},
            ),
        )


oauth2 = OAuth2AuthorizationCodeBearer(
    authorizationUrl="/authorize",
    tokenUrl="/token",
)

@router.get("/foo")
async def bar(token = Depends()):
    return token


