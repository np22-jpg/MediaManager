from pydantic_settings import BaseSettings


class ProwlarrConfig(BaseSettings):
    enabled: bool = False
    api_key: str = ""
    url: str = "http://localhost:9696"


class JackettConfig(BaseSettings):
    enabled: bool = False
    api_key: str = ""
    url: str = "http://localhost:9696"
    indexers: list[str] = ["all"]

class ScoringRule(BaseSettings):
    name: str
    score_modifier: int = 0
    negate: bool = False

class TitleScoringRule(ScoringRule):
    keywords: list[str]

class IndexerFlagScoringRule(ScoringRule):
    flags: list[str]

class ScoringRuleSet(BaseSettings):
    name: str
    tags: list[str] = []
    rule_names: list[str] = []

class IndexerConfig(BaseSettings):
    prowlarr: ProwlarrConfig
    jackett: JackettConfig
    title_scoring_rules: list[TitleScoringRule] = []
    indexer_flag_scoring_rules: list[IndexerFlagScoringRule] = []