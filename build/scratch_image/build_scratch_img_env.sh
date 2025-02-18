TMP=build/scratch_image/tmp/
PACKAGE=VirtualRegistryManagement
EXEC_FILE=$TMP/$PACKAGE
TMP_LINK_FILES=$TMP/link_files

echo $(pwd)
mkdir -p $TMP_LINK_FILES
linkPkgs=$(ldd $EXEC_FILE | egrep -o '\/[\/a-z0-9.+\_-]+')

for pkg in $linkPkgs; 
do 
    lPath=$TMP_LINK_FILES/$(dirname "$pkg" | cut -c2-)
    mkdir -p $lPath
    cp $pkg $lPath
done

cp /lib/x86_64-linux-gnu/libresolv.so.2 $TMP_LINK_FILES/lib/x86_64-linux-gnu/libresolv.so.2
cp /lib/x86_64-linux-gnu/libnss_dns.so.2 $TMP_LINK_FILES/lib/x86_64-linux-gnu/libnss_dns.so.2
cp /usr/local/go/lib/time/zoneinfo.zip $TMP_LINK_FILES/zoneinfo.zip
# sed -i '' 's/$(PREVERSION)/$(VERSION)/g' $(PWD)/build/AETriton.spec

