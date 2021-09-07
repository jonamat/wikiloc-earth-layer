# Generate assets
make gen-kml
make get-icons

# Create dirs
mkdir ./release/ 
mkdir ./release/windows-amd64/
mkdir ./release/linux-amd64/
mkdir ./release/linux-arm64/

# Build for each platform
GOOS=windows
GOARCH=amd64
go build -v -x -o ./release/windows-amd64/wikiloc-earth-layer.exe ./cmd/server/wikiloc-earth-layer.go &

GOOS=linux
GOARCH=amd64
go build -v -x -o ./release/linux-amd64/wikiloc-earth-layer ./cmd/server/wikiloc-earth-layer.go

GOOS=linux
GOARCH=arm64
go build -v -x -o ./release/linux-arm64/wikiloc-earth-layer ./cmd/server/wikiloc-earth-layer.go

# Copy web assets
cp -R ./web ./release/windows-amd64/web/
cp -R ./web ./release/linux-amd64/web/
cp -R ./web ./release/linux-arm64/web/

# Copy dotfile
cp -R ./.env ./release/windows-amd64/.env
cp -R ./.env ./release/linux-amd64/.env
cp -R ./.env ./release/linux-arm64/.env

# Zip folders
cd ./release/
zip -r ./windows-amd64.zip ./windows-amd64/
zip -r ./linux-amd64.zip ./linux-amd64/
zip -r ./linux-arm64.zip ./linux-arm64/

# Destroy release dirs
rm -rf ./windows-amd64
rm -rf ./linux-amd64
rm -rf ./linux-arm64