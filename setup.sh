#!/bin/bash
set -e
DOMAINS=(github.com api.github.com raw.github.com raw.githubusercontent.com avatars.githubusercontent.com avatars0.githubusercontent.com avatars1.githubusercontent.com avatars2.githubusercontent.com avatars3.githubusercontent.com github.githubassets.com objects.githubusercontent.com user-images.githubusercontent.com camo.githubusercontent.com cloud.githubusercontent.com gist.github.com github.io github.dev pages.github.com githubapp.com www.github.io)
echo "Adding GitHub domains to /etc/hosts..."
for d in "${DOMAINS[@]}"; do
    grep -q "127.0.0.1.*$d" /etc/hosts 2>/dev/null || echo "127.0.0.1 $d" | sudo tee -a /etc/hosts > /dev/null
done
echo "Done. ${#DOMAINS[@]} domains added."
