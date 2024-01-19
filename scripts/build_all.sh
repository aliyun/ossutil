#!/bin/bash
rootdir=$(cd `dirname $0`; cd ..; pwd)
output=$rootdir/dist
version=`cat ./lib/const.go | grep -E "Version.+ string =" | cut -d"=" -f2 | xargs`

#output folder
mkdir -p $output

#ossutil-$version-mac-amd64.zip & ossutil-mac-amd64.zip, files: ossutil,ossutilmac64
echo "start build ossutil for darwin on amd64"
cd $rootdir
env GOOS=darwin GOARCH=amd64 go build -o $output/ossutilmac64

cd $output
mkdir -p ossutil-$version-mac-amd64
cp -f ossutilmac64 ossutil-$version-mac-amd64/ossutil
cp -f ossutilmac64 ossutil-$version-mac-amd64/ossutilmac64
zip -r ossutil-$version-mac-amd64.zip ossutil-$version-mac-amd64
mv -f ossutil-$version-mac-amd64 ossutil-mac-amd64
zip -r ossutil-mac-amd64.zip ossutil-mac-amd64
rm -rf ossutil-mac-amd64
echo "ossutil for darwin on amd64 built successfully"

#ossutil-$version-mac-arm64.zip & ossutil-mac-arm64.zip, files: ossutil,ossutilmac64
echo "start build ossutil for darwin on arm64"
cd $rootdir
env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $output/ossutilmacarm64

cd $output
mkdir -p ossutil-$version-mac-arm64
cp -f ossutilmacarm64 ossutil-$version-mac-arm64/ossutil
cp -f ossutilmacarm64 ossutil-$version-mac-arm64/ossutilmac64
zip -r ossutil-$version-mac-arm64.zip ossutil-$version-mac-arm64
mv -f ossutil-$version-mac-arm64 ossutil-mac-arm64
zip -r ossutil-mac-arm64.zip ossutil-mac-arm64
rm -rf ossutil-mac-arm64

echo "ossutil for darwin on arm64 built successfully"

#ossutil-$version-windows-386 & ossutil-windows-386, files: ossutil.bat,ossutil.exe,ossutil32.exe 
echo "start build ossutil for windows on 386"
cd $rootdir
env GOOS=windows GOARCH=386 go build -o $output/ossutil32.exe

cd $output
mkdir -p ossutil-$version-windows-386
cp -f ossutil32.exe ossutil-$version-windows-386/ossutil.exe
cp -f ossutil32.exe ossutil-$version-windows-386/ossutil32.exe
cp -f $rootdir/scripts/ossutil.bat ossutil-$version-windows-386/ossutil.bat
zip -r ossutil-$version-windows-386.zip ossutil-$version-windows-386
mv -f ossutil-$version-windows-386 ossutil32
rm -f ossutil32/ossutil.exe
zip -r ossutil32.zip ossutil32
rm -rf ossutil32

echo "ossutil for windows on 386 built successfully"


#ossutil-$version-windows-amd64 & ossutil-windows-amd64, files: ossutil.bat,ossutil.exe,ossutil64.exe 
echo "start build ossutil for windows on amd64"
cd $rootdir
env GOOS=windows GOARCH=amd64 go build -o $output/ossutil64.exe

cd $output
mkdir -p ossutil-$version-windows-amd64
cp -f ossutil64.exe ossutil-$version-windows-amd64/ossutil.exe
cp -f ossutil64.exe ossutil-$version-windows-amd64/ossutil64.exe
cp -f $rootdir/scripts/ossutil.bat ossutil-$version-windows-amd64/ossutil.bat
zip -r ossutil-$version-windows-amd64.zip ossutil-$version-windows-amd64
mv ossutil-$version-windows-amd64 ossutil64
rm -f ossutil64/ossutil.exe
zip -r ossutil64.zip ossutil64
rm -rf ossutil64

echo "ossutil for windows on amd64 built successfully"

#ossutil-$version-linux-386 & ossutil-linux-386, files: ossutil,ossutil32
echo "start build ossutil for linux on 386"
cd $rootdir
env CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o $output/ossutil32

cd $output
mkdir -p ossutil-$version-linux-386
cp -f ossutil32 ossutil-$version-linux-386/ossutil
cp -f ossutil32 ossutil-$version-linux-386/ossutil32
zip -r ossutil-$version-linux-386.zip ossutil-$version-linux-386
mv -f ossutil-$version-linux-386 ossutil-linux-386
zip -r ossutil-linux-386.zip ossutil-linux-386
rm -rf ossutil-linux-386

echo "ossutil for linux on 386 built successfully"

#ossutil-$version-linux-amd64 & ossutil-linux-amd64, files: ossutil,ossutil64
echo "start build ossutil for linux on amd64"
cd $rootdir
env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $output/ossutil64

cd $output
mkdir -p ossutil-$version-linux-amd64
cp -f ossutil64 ossutil-$version-linux-amd64/ossutil
cp -f ossutil64 ossutil-$version-linux-amd64/ossutil64
zip -r ossutil-$version-linux-amd64.zip ossutil-$version-linux-amd64
mv -f ossutil-$version-linux-amd64 ossutil-linux-amd64
zip -r ossutil-linux-amd64.zip ossutil-linux-amd64
rm -rf ossutil-linux-amd64

echo "ossutil for linux on amd64 built successfully"

#ossutil-$version-linux-arm & ossutil-linux-arm, files: ossutil,ossutil32
echo "start build ossutil for linux on arm"
cd $rootdir
env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o $output/ossutilarm32

cd $output
mkdir -p ossutil-$version-linux-arm
cp -f ossutilarm32 ossutil-$version-linux-arm/ossutil
cp -f ossutilarm32 ossutil-$version-linux-arm/ossutil32
zip -r ossutil-$version-linux-arm.zip ossutil-$version-linux-arm
mv -f ossutil-$version-linux-arm ossutil-linux-arm
zip -r ossutil-linux-arm.zip ossutil-linux-arm
rm -rf ossutil-linux-arm

echo "ossutil for linux on arm built successfully"

#ossutil-$version-linux-arm64 & ossutil-linux-arm64, files: ossutil,ossutil64
echo "start build ossutil for linux on arm64"
cd $rootdir
env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $output/ossutilarm64

cd $output
mkdir -p ossutil-$version-linux-arm64
cp -f ossutilarm64 ossutil-$version-linux-arm64/ossutil
cp -f ossutilarm64 ossutil-$version-linux-arm64/ossutil64
zip -r ossutil-$version-linux-arm64.zip ossutil-$version-linux-arm64
mv -f ossutil-$version-linux-arm64 ossutil-linux-arm64
zip -r ossutil-linux-arm64.zip ossutil-linux-arm64
rm -rf ossutil-linux-arm64

echo "ossutil for linux on arm64 built successfully"

#calc hash for zip files
cd $output
for file in $(ls *.zip); do
    sha256sum $file >> sha256sum.log
done