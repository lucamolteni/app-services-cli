#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

PACKAGE_NAME="linux-x64"
if [[ $OSTYPE == 'darwin'* ]]; then
  PACKAGE_NAME="osx-x64"
fi
URL="https://github.com/microsoft/kiota/releases/download/v1.0.1/${PACKAGE_NAME}.zip"

COMMAND="kiota"
if ! command -v $COMMAND &> /dev/null
then
  echo "System wide kiota could not be found, using local version"
  if [[ ! -f $SCRIPT_DIR/kiota ]]
  then
    echo "Local kiota could not be found, downloading"
    rm -rf $SCRIPT_DIR/tmp-kiota
    mkdir -p $SCRIPT_DIR/tmp-kiota
    curl -sL $URL > $SCRIPT_DIR/tmp-kiota/kiota.zip
    unzip $SCRIPT_DIR/tmp-kiota/kiota.zip -d $SCRIPT_DIR/tmp-kiota

    mkdir -p $SCRIPT_DIR/tmp-kiota/bin
    cp $SCRIPT_DIR/tmp-kiota/*/kiota $SCRIPT_DIR/kiota
    chmod a+x $SCRIPT_DIR/kiota
    rm -rf $SCRIPT_DIR/tmp-kiota
  fi
  COMMAND="$SCRIPT_DIR/kiota"
fi


$COMMAND generate \
  --language go \
  --openapi https://raw.githubusercontent.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/main/openapi/kas-fleet-manager.yaml \
  --clean-output \
  -o $SCRIPT_DIR/../pkg/apisdk/kafkamgmt \
  --namespace-name github.com/redhat-developer/app-services-cli/pkg/apisdk/kafkamgmt

# TODO this is not the right openapi file
$COMMAND generate \
  --language go \
  --openapi https://raw.githubusercontent.com/redhat-developer/app-services-sdk-go/e660b42dabba43265b622693c94c748fb78ca62c/.openapi/service-accounts.yaml \
  --clean-output \
  -o $SCRIPT_DIR/../pkg/apisdk/svcacctmgmt \
  --namespace-name github.com/redhat-developer/app-services-cli/pkg/apisdk/svcacctmgmt

