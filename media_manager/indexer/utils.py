import logging

from media_manager.config import AllEncompassingConfig
from media_manager.indexer.config import ScoringRuleSet
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.movies.schemas import Movie
from media_manager.tv.schemas import Show

log = logging.getLogger(__name__)


def evaluate_indexer_query_result(
    query_result: IndexerQueryResult, ruleset: ScoringRuleSet
) -> IndexerQueryResult | None:
    title_rules = AllEncompassingConfig().indexers.title_scoring_rules
    indexer_flag_rules = AllEncompassingConfig().indexers.indexer_flag_scoring_rules
    for rule_name in ruleset.rule_names:
        for rule in title_rules:
            if rule.name == rule_name:
                if (
                    any(
                        keyword.lower() in query_result.title.lower()
                        for keyword in rule.keywords
                    )
                    and not rule.negate
                ):
                    query_result.score += rule.score_modifier
        for rule in indexer_flag_rules:
            if rule.name == rule_name:
                if (
                    any(flag in query_result.flags for flag in rule.flags)
                    and not rule.negate
                ):
                    query_result.score += rule.score_modifier
    if query_result.score <= 0:
        return None

    return query_result


def evaluate_indexer_query_results(
    query_results: list[IndexerQueryResult], media: Show | Movie, is_tv: bool
) -> list[IndexerQueryResult]:
    scoring_rulesets: list[ScoringRuleSet] = (
        AllEncompassingConfig().indexers.scoring_rule_sets
    )
    for ruleset in scoring_rulesets:
        if (
            (media.library in ruleset.libraries)
            or ("ALL_TV" in ruleset.libraries and is_tv)
            or ("ALL_MOVIES" in ruleset.libraries and not is_tv)
        ):
            log.debug(
                f"Applying scoring ruleset {ruleset.name} for {media.name} ({media.year})"
            )
            for result in query_results:
                log.debug(
                    f"Applying scoring ruleset {ruleset.name} for IndexerQueryResult {result.title} for {media.name} ({media.year})"
                )
                result = evaluate_indexer_query_result(
                    query_result=result, ruleset=ruleset
                )
                if not result:
                    log.debug(
                        f"Indexer query result {result.title} did not pass scoring ruleset {ruleset.name} with score {result.score}, removing from results."
                    )
                    query_results.remove(result)
                else:
                    log.debug(
                        f"Indexer query result {result.title} passed scoring ruleset {ruleset.name} with score {result.score}."
                    )
    query_results.sort(reverse=True)
    return query_results
