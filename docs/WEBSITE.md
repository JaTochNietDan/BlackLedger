# Game website

The public landing page is built from `website/` and published to
https://blackledger.oleary.club/ by `.github/workflows/pages.yml`.
GitHub repository Settings → Pages must use **GitHub Actions** as its source.
Changes to the website, screenshots, recorded encounter or build script on
`main` publish automatically. Pull requests build without deploying.

## Local preview

```sh
python3 scripts/build-site.py
python3 -m http.server 8876 --bind 127.0.0.1 --directory .runtime/site
```

Open http://127.0.0.1:8876/. The build uses Python's standard library and copies
only the landing page, three selected screenshots and its own assets. It embeds
the recorded encounter from `docs/examples/director-encounter.json`; the choices
compare that encounter's stakes and do not run a live AI model. Game saves,
services and model files are never published by this workflow.

The page introduces the AI director experiment, shows actual development
screenshots and links to desktop builds. GitHub Pages hosts the marketing site;
the game and local AI run on the player's computer.

## Verification

The initial implementation was checked in desktop and narrow mobile browser
layouts, with all four story approaches, keyboard dismissal and focus return
for the screenshot dialog, and browser console checks. The workflow passed
`actionlint`. Download links describe development artifacts separately from
published releases; platform acceptance and asset provenance work remain tracked
in `docs/RELEASING.md` and `ASSETS.md`.

The first deployment was verified live on GitHub Pages, including its project
subpath assets. Pages uses GitHub Actions as the source with HTTPS enforced.

## Custom domain

`blackledger.oleary.club` is configured in the repository's Pages settings.
Its DNS-only CNAME points to `jatochnietdan.github.io` (without a repository
path). Terraform owns this record in the `ops` repository at
`terraform/zones/oleary.club/cloudflare_record.tf`.

GitHub Actions publishing uses the Pages custom-domain setting; it does not
require a source `CNAME` file. Keep the canonical and social preview URLs in
`website/index.html` and the README website link aligned with that setting.
