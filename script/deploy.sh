#!/bin/bash

echo "Deploying the application..."

# Example deployment using Docker
docker build -t mlvt:latest .
docker run -d -p 8080:8080 --name mlvt_container mlvt:latest

echo "Application deployed and running on port 8080."
