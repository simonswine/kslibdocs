const versionLinks =document.querySelectorAll('a.kind-version');

for (const el of versionLinks) {
    el.addEventListener('mouseover', (event) => {
        const parts = event.target.id.split("-");
        const summaryID = `summary-group-kind-${parts[2]}-${parts[3]}`;
        const version = event.target.dataset.version;

        const summaries = document.querySelectorAll(`#${summaryID} .kind-summary`);
        for (const summary of summaries) {
            summary.style.display = "none";
        }

        const toShowVer = `${parts[2]}-${version}-${parts[3]}`;
        const toShow = document.getElementById(toShowVer);
        toShow.style.display = "block";
    });
}
