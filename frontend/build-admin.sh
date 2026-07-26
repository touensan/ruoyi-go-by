#!/usr/bin/env bash
set -euo pipefail

frontend_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
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

printf 'built admin frontend: %s\n' "$variant"
