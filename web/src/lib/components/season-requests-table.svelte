<script lang="ts">
    import {
        convertTorrentSeasonRangeToIntegerRange, getFullyQualifiedShowName,
        getTorrentQualityString,
        getTorrentStatusString
    } from "$lib/utils.js";
    import type {SeasonRequest} from "$lib/types.js";
    import CheckmarkX from "$lib/components/checkmark-x.svelte";
    import * as Table from "$lib/components/ui/table/index.js";

    let {
        requests, filter = () => {
            return true
        }
    }: { requests: SeasonRequest[], filter: (SeasonRequest) => boolean } = $props();

</script>

<Table.Root>
    <Table.Caption>A list of all requests.</Table.Caption>
    <Table.Header>
        <Table.Row>
            <Table.Head>Show</Table.Head>
            <Table.Head>Season</Table.Head>
            <Table.Head>Minimum Quality</Table.Head>
            <Table.Head>Wanted Quality</Table.Head>
            <Table.Head>Requested by</Table.Head>
            <Table.Head>Approved</Table.Head>
            <Table.Head>Approved by</Table.Head>
            <Table.Head>Actions</Table.Head>
        </Table.Row>
    </Table.Header>
    <Table.Body>
        {#each requests as request (request.id)}
            {#if filter(request)}
                <Table.Row>
                    <Table.Cell class="font-medium">
                        {getFullyQualifiedShowName(request.show)}
                    </Table.Cell>
                    <Table.Cell>
                        {request.season.number}
                    </Table.Cell>
                    <Table.Cell class="font-medium">
                        {getTorrentQualityString(request.min_quality)}
                    </Table.Cell>
                    <Table.Cell class="font-medium">
                        {getTorrentQualityString(request.wanted_quality)}
                    </Table.Cell>
                    <Table.Cell>
                        {request.requested_by?.email}
                    </Table.Cell>
                    <Table.Cell>
                        <CheckmarkX state={request.authorized}/>
                    </Table.Cell>
                    <Table.Cell>
                        {request.authorized_by?.email}
                    </Table.Cell>
                </Table.Row>
            {/if}
        {/each}
    </Table.Body>
</Table.Root>