#!/bin/bash

set -e

rm -rf /tmp/vault

DIRECTORY=`dirname $0`

vault server -config ${DIRECTORY}/vault.hcl