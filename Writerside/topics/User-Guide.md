# Usage

If you are coming from Radarr or Sonarr you will find that MediaManager does things a bit differently.
Instead of completely automatically downloading and managing your media, MediaManager focuses on providing an
easy-to-use interface to guide you through the process of finding and downloading media. Advanced features like multiple
qualities of a show/movie necessitate such a paradigm shift. __So here is a quick step-by-step guide to get you started:
__

<tabs>
    <tab id="as-a-user" title="as a user">
        <procedure title="Downloading/Requesting a show" id="request-show-user">
           <step>Add a show on the "Add Show" page</step>
           <step>After adding the show you will be redirected to the show's page.</step>
           <step>There you can click the "Request Season" button.</step>
           <step>Select one or more seasons that you want to download</step>
           <step>Then select the "Min Quality", this will be the minimum resolution of the content to download.</step>
           <step>Then select the "Wanted Quality", this will be the <strong>maximum</strong> resolution of the content to download.</step>
           <step>Finally click Submit request, though this is not the last step!</step>
           <step>An administrator first has to approve your request for download, only then will the requested content be downloaded.</step>
           <p>Congratulation! You've downloaded a show.</p>
        </procedure>
    </tab>
    <tab id="as-an-admin" title="as an admin">
        <procedure title="Requesting a show" id="request-show-admin">
           <step>Add a show on the "Add Show" page</step>
           <step>After adding the show you will be redirected to the show's page.</step>
           <step>There you can click the "Request Season" button.</step>
           <step>Select one or more seasons that you want to download</step>
           <step>Then select the "Min Quality", this will be the minimum resolution of the content to download.</step>
           <step>Then select the "Wanted Quality", this will be the <strong>maximum</strong> resolution of the content to download.</step>
           <step>Finally click Submit request, as you are an admin, your request will be automatically approved.</step>
           <p>Congratulation! You've downloaded a show.</p>
        </procedure>
        <procedure title="Downloading a show" id="download-show-admin">
            <p>You can only directly download a show if you are an admin!</p>
           <step>Go to a show's page.</step>
           <step>There you can click the "Download Season" button.</step>
           <step>Enter the season's number that you want to download</step>
           <step>Then optionally select the "File Path Suffix", <strong>it needs to be unique per season per show!</strong> </step>
           <step>Then click "Download" on a torrent that you want to download.</step>
           <p>Congratulation! You've downloaded a show.</p>
        </procedure>
        <procedure title="Managing requests" id="approving-request-admin">
           <p>Users need their requests to be approved by an admin, to do this follow these steps:</p>
           <step>Go to the "Requests" page.</step>
           <step>There you can approve, delete or modify a user's request.</step>
        </procedure>
    </tab>
</tabs>


