#!/usr/bin/env sh
set -eu

project_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
output_dir=${1:?output directory required}
binary=${2:?Go binary required}
frontend_dist=${3:?frontend dist required}
build_info=${4:?BUILD_INFO.json required}
version=${5:-1.1.3}

test -x "$binary"
test -f "$frontend_dist/index.html"
test -f "$build_info"
mkdir -p "$output_dir"
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT HUP INT TERM
mkdir -p "$stage/web/admin" "$stage/runtime/uploads"

cp "$binary" "$stage/ruoyi-go-by"
cp -a "$frontend_dist/." "$stage/web/admin/"
cp "$project_dir/frontend/RuoYi-Vue3-ts/LICENSE" "$stage/web/admin/LICENSE"
for file in application-example.yaml ruoyi.sql service.sh service.bat license README.md SECURITY.md; do
    cp "$project_dir/$file" "$stage/"
done
cp -a "$project_dir/docs" "$stage/"
cp "$build_info" "$stage/BUILD_INFO.json"

archive="$output_dir/ruoyi-go-by-v${version}-linux-x86_64.tar.gz"
tar -C "$stage" --sort=name --mtime='@0' --owner=0 --group=0 --numeric-owner -czf "$archive" .
(cd "$output_dir" && sha256sum "$(basename "$archive")" > SHA256SUMS)
printf 'Package: %s\n' "$archive"
