#!/usr/bin/env bash
set -euo pipefail
root=$(cd "$(dirname "$0")/.." && pwd)
workflow="$root/.github/workflows/m7-pretag-approval.yml"
[[ "$(grep -Ec '^  create_tag:$' "$workflow")" -eq 1 ]]
create_tag_job=$(awk '
  /^  create_tag:$/ { in_job = 1 }
  in_job && /^  [[:alnum:]_-]+:$/ && $0 != "  create_tag:" { exit }
  in_job { print }
' "$workflow")
grep -Fq '    needs: quality' <<<"$create_tag_job"
grep -Fq '      contents: write' <<<"$create_tag_job"
grep -Fq '[[ "$GITHUB_REF" == "refs/heads/main" ]]' <<<"$create_tag_job"
grep -Fq '[[ "$GITHUB_WORKFLOW_SHA" == "$GITHUB_SHA" ]]' <<<"$create_tag_job"
grep -Fq '[[ "$VERSION" =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]' <<<"$create_tag_job"
grep -Fq 'ref=refs/tags/$tag' <<<"$create_tag_job"
grep -Fq 'object=$GITHUB_SHA' <<<"$create_tag_job"
grep -Fq 'sha=$tag_object' <<<"$create_tag_job"
grep -Fq 'GH_TOKEN: ${{ github.token }}' <<<"$create_tag_job"
! grep -Fq 'actions/checkout' <<<"$create_tag_job"
! grep -Fq 'scripts/' <<<"$create_tag_job"
[[ "$(grep -Fc 'contents: write' "$workflow")" -eq 1 ]]
if grep -Eq -- '--force|--force-with-lease' <<<"$create_tag_job"; then
  printf 'release tag workflow must never force-update refs\n' >&2
  exit 1
fi
if grep -Eq '^  dispatch_release:' "$workflow"; then
  printf 'core module must not have a release dispatch job\n' >&2
  exit 1
fi
if grep -Fq 'actions: write' "$workflow"; then
  printf 'core module must not grant Actions write permission\n' >&2
  exit 1
fi
if [[ -f "$root/.github/workflows/release.yml" ]]; then
  printf 'core module has an unexpected release workflow\n' >&2
  exit 1
fi
printf 'PASS: immutable tag creation is isolated and no dispatch permission is granted\n'
