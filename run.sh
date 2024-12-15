#!/bin/bash

#set -e  # Exit immediately if a command exits with a non-zero status
#set -x  # Print commands and their arguments as they are executed

ports=(3000 8083)
for port in "${ports[@]}"; do
  echo "Checking port $port"
  pids=$(lsof -ti tcp:$port)
  if [ -n "$pids" ]; then
    echo "Killing processes on port $port: $pids"
    echo "$pids" | xargs kill -9
  else
    echo "No processes found on port $port"
  fi
done

#codesign --sign "Jonathan Leahy" --force ./public/argocd

ttab "cd server/src; go build -o ../server; cd ../; ./server --webserver"
ttab "cd dashboard/; nvm use 20.11.1; npm install; npm run dev;"
