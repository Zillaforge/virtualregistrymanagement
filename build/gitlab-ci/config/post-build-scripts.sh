#!/bin/sh
if [ ${AUTOBUILD_PROD_MODE} ]
then
    echo "AUTOBUILD_PROD_MODE is ON"
    VERSION=`cat VERSION`
    AUTOBUILD_RELEASE_DIR="${AUTOBUILD_RELEASE_DIR:-/home/gitlab-runner/releases}"
    AUTOBUILD_RELEASE_DIR="${AUTOBUILD_RELEASE_DIR}/${CI_PROJECT_NAMESPACE}/VirtualRegistryManagement/${VERSION}"
    mkdir -p "${AUTOBUILD_RELEASE_DIR}"
    mv tmp/*/*-"${VERSION}"* "${AUTOBUILD_RELEASE_DIR}"
    echo "The AutoBuild Release Dir is: ${AUTOBUILD_RELEASE_DIR}"
else
    echo "AUTOBUILD_PROD_MODE is OFF"
fi