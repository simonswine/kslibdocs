#!/usr/bin/env bash
# updates the embedded site template

set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.."; pwd)
SRC=${1:-}
TEMPLATE_ROOT="${ROOT}/pkg/site/site_template"

if [ -z "${SRC}" ]; then
    echo "$0 <src>"
    exit 1
fi

if [ ! -e "${SRC}" ]; then
    echo "${SRC} does not exist"
    exit 1
fi

rm -rf ${TEMPLATE_ROOT}
cp -r ${SRC} ${TEMPLATE_ROOT}

# clean up
rm -rf ${TEMPLATE_ROOT}/public
rm -rf ${TEMPLATE_ROOT}/content
rm -rf ${TEMPLATE_ROOT}/node_modules