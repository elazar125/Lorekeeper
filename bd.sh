#! /usr/bin/bash

# [b]uild & [d]eploy
env CGO_ENABLED=0 GOOS=linux go build -o dist/lorekeeper.new
scp dist/lorekeeper.new spellslingerer:~/lorekeeper
ssh spellslingerer -t << 'EOF'
systemctl stop lorekeeper
mv ~/lorekeeper/lorekeeper.new ~/lorekeeper/lorekeeper
chmod +x ~/lorekeeper/lorekeeper
systemctl start lorekeeper
EOF
