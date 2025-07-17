import logging

from media_manager.config import AllEncompassingConfig
from media_manager.indexer.config import ScoringRuleSet
from media_manager.indexer.schemas import IndexerQueryResult
from media_manager.movies.schemas import Movie
from media_manager.tv.schemas import Show

log = logging.getLogger(__name__)


def evaluate_indexer_query_result(
    query_result: IndexerQueryResult, ruleset: ScoringRuleSet
) -> (IndexerQueryResult, bool):
    title_rules = AllEncompassingConfig().indexers.title_scoring_rules
    indexer_flag_rules = AllEncompassingConfig().indexers.indexer_flag_scoring_rules
    for rule_name in ruleset.rule_names:
        for rule in title_rules:
            if rule.name == rule_name:
                log.debug(f"Applying rule {rule.name} to {query_result.title}")
                if (
                    any(
                        keyword.lower() in query_result.title.lower()
                        for keyword in rule.keywords
                    )
                    and not rule.negate
                ):
                    log.debug(
                        f"Rule {rule.name} with keywords {rule.keywords} matched for {query_result.title}"
                    )
                    query_result.score += rule.score_modifier
                elif (
                    not any(
                        keyword.lower() in query_result.title.lower()
                        for keyword in rule.keywords
                    )
                    and rule.negate
                ):
                    log.debug(
                        f"Negated rule {rule.name} with keywords {rule.keywords} matched for {query_result.title}"
                    )
                    query_result.score += rule.score_modifier
                else:
                    log.debug(
                        f"Rule {rule.name} with keywords {rule.keywords} did not match for {query_result.title}"
                    )
        for rule in indexer_flag_rules:
            if rule.name == rule_name:
                log.debug(f"Applying rule {rule.name} to {query_result.title}")
                if (
                    any(flag in query_result.flags for flag in rule.flags)
                    and not rule.negate
                ):
                    log.debug(
                        f"Rule {rule.name} with flags {rule.flags} matched for {query_result.title} with flags {query_result.flags}"
                    )
                    query_result.score += rule.score_modifier
                elif (
                    not any(flag in query_result.flags for flag in rule.flags)
                    and rule.negate
                ):
                    log.debug(
                        f"Negated rule {rule.name} with flags {rule.flags} matched for {query_result.title} with flags {query_result.flags}"
                    )
                    query_result.score += rule.score_modifier
                else:
                    log.debug(
                        f"Rule {rule.name} with flags {rule.flags} did not match for {query_result.title} with flags {query_result.flags}"
                    )
    if query_result.score <= 0:
        return query_result, False

    return query_result, True


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
                result, passed = evaluate_indexer_query_result(
                    query_result=result, ruleset=ruleset
                )
                if not passed:
                    log.debug(
                        f"Indexer query result {result.title} did not pass scoring ruleset {ruleset.name} with score {result.score}, removing from results."
                    )
                else:
                    log.debug(
                        f"Indexer query result {result.title} passed scoring ruleset {ruleset.name} with score {result.score}."
                    )

    query_results = [result for result in query_results if result.score >= 0]
    query_results.sort(reverse=True)
    return query_results
