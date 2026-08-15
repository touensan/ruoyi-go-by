#!/usr/bin/env bash
set -euo pipefail

frontend_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
project_root="$(cd "$frontend_root/.." && pwd)"
variant="${APIAUTH_ADMIN_FRONTEND:-whiteyun-vue}"

case "$variant" in
  whiteyun-vue|ruoyi-vue)
    ;;
  *)
    echo "APIAUTH_ADMIN_FRONTEND must be whiteyun-vue or ruoyi-vue" >&2
    exit 2
    ;;
esac

(
  cd "$frontend_root/$variant"
  npm ci
  npm run build:prod
)

mkdir -p "$project_root/web/docs"
cp -R "$frontend_root/docs/." "$project_root/web/docs/"

printf 'built admin frontend: %s\n' "$variant"
