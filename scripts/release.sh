VERSION=`git describe --tags`
PREFIX=wikiloc-earth-layer_${VERSION}_

echo "Building assets for the relese ${VERSION}..."

# Generate assets
make gen-kml
make get-icons

# Cleanup relaese dir
rm -rf ./release

# Create dirs
mkdir ./release/ \
./release/${PREFIX}windows-amd64/ \
./release/${PREFIX}linux-amd64/ \
./release/${PREFIX}linux-arm64/

# Build for each platform
GOOS=windows GOARCH=amd64 go build -o ./release/${PREFIX}windows-amd64/wikiloc-earth-layer.exe ./cmd/server/wikiloc-earth-layer.go &
GOOS=linux GOARCH=amd64 go build -o ./release/${PREFIX}linux-amd64/wikiloc-earth-layer ./cmd/server/wikiloc-earth-layer.go
GOOS=linux GOARCH=arm64 go build -o ./release/${PREFIX}linux-arm64/wikiloc-earth-layer ./cmd/server/wikiloc-earth-layer.go

# Copy assets
cp -R ./web ./.env ./config.yml ./release/${PREFIX}windows-amd64/
cp -R ./web ./.env ./config.yml ./release/${PREFIX}linux-amd64/
cp -R ./web ./.env ./config.yml ./release/${PREFIX}linux-arm64/

# Zip folders
cd ./release/
zip -r ./${PREFIX}windows-amd64.zip ./${PREFIX}windows-amd64/
zip -r ./${PREFIX}linux-amd64.zip ./${PREFIX}linux-amd64/
zip -r ./${PREFIX}linux-arm64.zip ./${PREFIX}linux-arm64/

# Destroy release dirs
rm -rf ./${PREFIX}windows-amd64
rm -rf ./${PREFIX}linux-amd64
rm -rf ./${PREFIX}linux-arm64
