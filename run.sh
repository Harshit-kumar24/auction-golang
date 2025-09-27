echo "building image..."
docker build -t eauction .

echo "running eauction container..."
docker-compose up -d