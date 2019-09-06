#!/usr/bin/env bash

pushd `dirname $0` > /dev/null
SCRIPTPATH=`pwd -P`
popd > /dev/null
SCRIPTFILE=`basename $0`

${SCRIPTPATH}/stop.sh
sleep 5 # Waiting for closing
${SCRIPTPATH}/start.sh
