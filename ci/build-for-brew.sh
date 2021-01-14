#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
set -e

S3_BUCKET=im-saas-build-public-artifacts/tools
BUILD_DATE=$(date -u '+%Y-%m-%dT%H%M%S')

cwd=$(pwd)

function usage() {
    echo ""
    echo "Abstract : "
    echo "${0} [application-name]"
    echo ""
    echo "Example:"
    echo "${0} hello-world --master"
    echo "will use the commit id as version for brew otherwise the latest tag will be used"
    exit 0
}

master=false

POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    --master)
    master=true
    shift # past value
    ;;
    --cli)
    master=true
    shift # past value
    ;;
    --help)
      usage
    ;;
    *)    # unknown option
    POSITIONAL+=("$1") # save it in an array for later
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters

ApplicationName=${1}
RepoName=`echo $ApplicationName | sed 's/-//g'`

echo "build for brew master=${master}"
echo ""