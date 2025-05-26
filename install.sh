#!/bin/bash
set -euo pipefail

# --- Pre-flight checks ---
for cmd in jq curl tar; do
  if ! command -v $cmd &>/dev/null; then
    echo "❌ ERROR: Required command '$cmd' is not installed."
    exit 1
  fi
done

# --- Config ---
REPO="nqngo/upbank-cli"

# --- Detect latest release tag ---
TAG=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | jq -r .tag_name)

if [[ -z "$TAG" || "$TAG" == "null" ]]; then
  echo "❌ ERROR: Failed to determine the latest release tag."
  exit 1
fi

# --- Detect OS and architecture ---
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  *) echo "❌ ERROR: Unsupported architecture: $ARCH"; exit 1 ;;
esac

FILENAME="upbank-cli_${TAG#v}_${OS}_${ARCH}.tar.gz"

# --- Get asset ID from GitHub API ---
ASSET_ID=$(curl -s "https://api.github.com/repos/$REPO/releases/tags/$TAG" |
  jq -r ".assets[] | select(.name == \"$FILENAME\") | .id")

if [[ -z "$ASSET_ID" ]]; then
  echo "❌ ERROR: Asset '$FILENAME' not found in release '$TAG'"
  exit 1
fi

# --- Download the asset ---
echo "📥 Downloading $FILENAME..."
curl -L -H "Accept: application/octet-stream" \
     "https://api.github.com/repos/$REPO/releases/assets/$ASSET_ID" \
     -o "$FILENAME"

# --- Extract and install ---
echo "📦 Extracting $FILENAME..."
tar -xzf "$FILENAME"
chmod +x upbank-cli

echo "🔧 Installing to /usr/local/bin (requires sudo)..."
sudo mv upbank-cli /usr/local/bin/

# Clean up downloaded file
echo "🧹 Cleaning up..."
rm -f "$FILENAME"

echo "✅ upbank-cli installed successfully!" 
