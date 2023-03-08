#!/bin/bash

set -eu
set -o pipefail

# deploys all the example contracts

for initsh in $(ls **/contracts/init.sh); do
	exampledir=$(dirname $initsh)
	echo ""
	echo "+---------------------------------"
	echo "| Deploying ${exampledir}..."
	echo "+---------------------------------"
	( cd "${exampledir}" && ./init.sh )
done
