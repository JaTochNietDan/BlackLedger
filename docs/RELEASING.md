# Preparing public previews

Target repository: https://github.com/JaTochNietDan/black-ledger

The latest user scope prioritizes publishing source and automated distributions.
The license requirement is attribution plus no free commercial use: original code
uses PolyForm Noncommercial 1.0.0 with NOTICE. Describe it as source-available.
Do not claim OSI open-source status or that third-party media rights are settled.

## Automated builds

`.github/workflows/build.yml` runs on branch pushes, pull requests, manual dispatch
and version tags. A Linux verification job validates the workflow and runs Go vet/race tests, frontend tests
and the production build. Five native runner jobs then build, archive and smoke
test Windows x64, macOS x64/ARM64 and Linux x64/ARM64 distributions. Artifacts
expire after 14 days. Version tags beginning with `v` create a **draft prerelease**
with archives and SHA-256 files once all jobs pass. Review then publish the draft.
Non-tag builds never create releases. Only the release job has contents-write.

The packaging script can cross-compile locally with CGO disabled. Version labels
beginning with `v` require a clean committed checkout; use `dev-` labels for local QA. Native CI checks
are still necessary: successfully cross-compiling does not establish that a build
runs on another OS. Native smoke tests start the extracted executable from an
unrelated directory, fetch the UI/bundles and API, and create only a temporary save.
They also restore a stopped-game backup into a replacement game folder, verify
receipt persistence and check that the original save remains untouched. This
tests same-version folder replacement, not compatibility with all old saves.
They do not establish graphics performance, installer acceptance or a full campaign.

## First-publication gates

- Resolve every outstanding asset entry in ASSETS.md, including supplied WAVs.
- Review all Git history for private material and redistribution rights, not only
  the working tree. Consider a reviewed initial snapshot if history cannot be
  safely published; preserve the original private repository either way.
- Commit coherent gameplay changes currently owned by other ongoing work before
  selecting a release revision. Never release uncommitted local work by accident.
- Create the repository under the requested owner using an authenticated GitHub
  session; no remote or credential is assumed by these scripts.
- Run GitHub CI and inspect all five native artifacts. Local cross-build evidence
  alone does not satisfy this gate.
- Play a fresh 20–30-minute campaign, restart/reload it, verify consequences and
  save continuity, and inspect city/interior performance and controls in supported
  browsers. Existing docs/GOAL.md acceptance gaps remain visible.
- Decide whether an unsigned alpha archive is acceptable for the first preview.
  Signing/notarization need the owner's platform credentials; they are not faked
  or bypassed by this pipeline. Native installers remain future work.
- Review README, notices, release notes, known issues and checksums, then publish
  the draft prerelease. Publish a stable release only after its acceptance passes.

## Tagging an approved revision

After publishing the reviewed source and obtaining a green branch build, create
an annotated tag such as `v0.1.0-alpha.1` on the approved commit and push that tag.
Do not reuse an existing release tag or force-update published history. The
workflow creates a draft, so unfinished binaries are not presented as a stable
release automatically. Test manual runs without tags to verify the workflow first.

## Tooling references

Runner labels are taken from GitHub's [hosted runner reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners).
The workflow grants release writes only to the final job, following the
[workflow permissions reference](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax).
Actions are pinned to reviewed commit IDs; Dependabot proposes updates.
The original [PolyForm license source](https://github.com/polyformproject/polyform-licenses/blob/1.0.0/PolyForm-Noncommercial-1.0.0.md)
is included verbatim, with project-specific required attribution kept in NOTICE.
