#!/usr/bin/env bash
PACKAGE="VirtualRegistryManagement"
TMP=tmp/rpm
VERSION=$(cat VERSION)
EXEC=$PACKAGE\_$VERSION
SOURCES=$TMP/SOURCES
SOURCES_PACKAGE=$SOURCES/$PACKAGE\-$VERSION
SPECS=$TMP/SPECS
RELEASE=1
TOPDIR="_topdir $(pwd)/$TMP"

echo "Make Temp Folder"
mkdir -p $TMP

echo "Make SPECS Folder"
mkdir -p $SPECS
cp build/$PACKAGE.spec $SPECS

echo "Make SOURCE Folder"
mkdir -p $SOURCES

echo "Make PACKAGE Folder"
mkdir -p  $SOURCES_PACKAGE
cp systemd/$PACKAGE.service $SOURCES_PACKAGE
cp tmp/$EXEC $SOURCES_PACKAGE/$PACKAGE
tar -zcvf $SOURCES/$PACKAGE\_$VERSION.tar.gz --directory=$SOURCES/ $PACKAGE\-$VERSION

rpmbuild --define "$TOPDIR" -bb $(pwd)/$SPECS/$PACKAGE.spec

mv $TMP/RPMS/x86_64/$PACKAGE-$VERSION-1.x86_64.rpm tmp/
rm -rf $TMP