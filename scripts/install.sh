#!/bin/sh
# This script provides easy installation of the lastest version 
# of Artisan (if is needed) into /usr/loca/bin directory.
# It's suitable for any CI/CD tool. You can invoke this script 
# with command:
#
#      $ curl -sfL https://artisan.unravela.com/install.sh | sh
#

# Ensure the temp and isntall directory
echo "Initialize..."
TMPDIR=$(mktemp -d)
INSTALLDIR=/usr/local/bin

# Download and untar the archive to temp. directory
echo "Downloading latest..."
URL=`curl -s https://api.github.com/repos/unravela/artisan/releases/latest | grep "browser_download_url.*Linux.*tar.gz" | cut -d : -f 2,3 | tr -d \"`
echo $URL
curl -s -L -o $TMPDIR/artisan.tar.gz $URL

# extract and install to ./bin/ folder
echo "Extracting..."
tar -xzf $TMPDIR/artisan.tar.gz -C $TMPDIR
mkdir -p $INSTALLDIR
cp $TMPDIR/artisan $INSTALLDIR/artisan

# cleanup
echo "Cleanup..."
rm -rf $TMPDIR

