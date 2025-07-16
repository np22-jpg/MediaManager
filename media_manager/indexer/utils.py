from media_manager.config import AllEncompassingConfig
from media_manager.indexer.config import ScoringRuleSet
from media_manager.indexer.schemas import IndexerQueryResult


def evaluate_indexer_query_result(query_result: IndexerQueryResult, ruleset: ScoringRuleSet) -> IndexerQueryResult|None:
    title_rules = AllEncompassingConfig().indexers.title_scoring_rules
    indexer_flag_rules = AllEncompassingConfig().indexers.indexer_flag_scoring_rules
    for rule_name in ruleset.rule_names:
        for rule in title_rules:
            if rule.name == rule_name:
                if any(keyword.lower() in query_result.title.lower() for keyword in rule.keywords) and not rule.negate:
                    query_result.score += rule.score_modifier
        for rule in indexer_flag_rules:
            if rule.name == rule_name:
                if any(flag in query_result.flags for flag in rule.flags) and not rule.negate:
                    query_result.score += rule.score_modifier
    if query_result.score <= 0:
        return None

    return query_result