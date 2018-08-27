#!/usr/bin/env bash
# creates a copy of the site_template to hack on

set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.."; pwd)
DEST=${1:-}

if [ -z "${DEST}" ]; then
    echo "$0 <dest>"
    exit 1
fi

if [ -e "${DEST}" ]; then
    echo "${DEST} exists"
    exit 1
fi

cp -r pkg/site/site_template ${DEST}