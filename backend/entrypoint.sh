#!/bin/sh

set -o errexit
set -o nounset
set -o pipefail

ATLAS_ENV=${ATLAS_ENV:-"prod"}
ATLAS_URL=${ATLAS_URL:-${DATABASE_CONNECTION_STRING:-""}}

atlas migrate apply \
	--env "${ATLAS_ENV}" \
	--url "${ATLAS_URL}"

./instancer
