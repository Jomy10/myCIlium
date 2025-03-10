set -x
set -e

pkg install -y \
  curl

# Download GO
cd /

curl -L -O https://go.dev/dl/go1.24.1.freebsd-amd64.tar.gz

# Check Checksum

checksum="47d7de8bb64d5c3ee7b6723aa62d5ecb11e3568ef2249bbe1d4bbd432d37c00c"
file_checksum="$(sha256sum go1.24.1.freebsd-amd64.tar.gz | awk '{print $1}')"

if [ $checksum != $file_checksum ]; then
  echo "Invalid checksum!!"
  exit 1
fi

# Install go

tar xvf go1.24.1.freebsd-amd64.tar.gz

# Build server

cd /server

/go/bin/go build -o mycilium-orchestrator -ldflags "-s -w"

mv mycilium-orchestrator /mycilium-orchestrator

# Remove unneeded files

rm -rf /server
rm -rf /go
rm -rf /root/go
