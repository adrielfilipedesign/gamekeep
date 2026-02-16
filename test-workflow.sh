#!/bin/bash

# SaveSync Test Script
# This script demonstrates the basic workflow of SaveSync

set -e

echo "==================================="
echo "SaveSync - Example Usage Script"
echo "==================================="
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BINARY="./savesync"

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo "Error: savesync binary not found!"
    echo "Please build first: go build -o savesync cmd/main.go"
    exit 1
fi

echo -e "${BLUE}Step 1: Adding a test game${NC}"
$BINARY add-game --name "Test Game" --path "/tmp/savesync-test/saves"
echo ""

echo -e "${BLUE}Step 2: Listing registered games${NC}"
$BINARY list-games
echo ""

echo -e "${BLUE}Step 3: Creating some test save files${NC}"
mkdir -p /tmp/savesync-test/saves
echo "Save data 1" > /tmp/savesync-test/saves/save1.dat
echo "Save data 2" > /tmp/savesync-test/saves/save2.dat
echo "Player progress: Level 10" > /tmp/savesync-test/saves/progress.txt
echo -e "${GREEN}✓ Test save files created${NC}"
echo ""

echo -e "${BLUE}Step 4: Creating first checkpoint${NC}"
$BINARY checkpoint --game "test_game" --name "Initial Save" --note "Starting point"
echo ""

echo -e "${BLUE}Step 5: Modifying save files${NC}"
echo "Save data 1 - MODIFIED" > /tmp/savesync-test/saves/save1.dat
echo "Player progress: Level 25" > /tmp/savesync-test/saves/progress.txt
echo -e "${GREEN}✓ Save files modified${NC}"
echo ""

echo -e "${BLUE}Step 6: Creating second checkpoint${NC}"
$BINARY checkpoint --game "test_game" --name "Mid Game" --note "Level 25 checkpoint"
echo ""

echo -e "${BLUE}Step 7: Listing all checkpoints${NC}"
$BINARY list --game "test_game"
echo ""

echo -e "${BLUE}Step 8: Checking current save content${NC}"
echo "Current progress.txt content:"
cat /tmp/savesync-test/saves/progress.txt
echo ""

echo -e "${BLUE}Step 9: Restoring first checkpoint${NC}"
# Get the first checkpoint ID (this is simplified, in real usage you'd parse the output)
CHECKPOINT_ID=$($BINARY list --game "test_game" | grep -A 1 "Initial Save" | tail -1 | awk '{print $1}')
if [ ! -z "$CHECKPOINT_ID" ]; then
    $BINARY restore --checkpoint "$CHECKPOINT_ID"
    echo ""
    
    echo -e "${BLUE}Step 10: Verifying restored content${NC}"
    echo "Restored progress.txt content:"
    cat /tmp/savesync-test/saves/progress.txt
    echo ""
fi

echo -e "${GREEN}==================================="
echo "✓ Test workflow completed successfully!"
echo "===================================${NC}"
echo ""
echo "Cleanup: To remove test data, run:"
echo "  rm -rf /tmp/savesync-test"
echo "  rm -rf ~/.savesync"
