#!/usr/bin/env bash
# launch minio
MINIO_ROOT_USER=admin MINIO_ROOT_PASSWORD=password ~/minio/minio server ~/minio/data --console-address ":9001" & MINIO_PID=$!
printf "MinIO process ID: %s\n" $MINIO_PID
echo "kill $MINIO_PID" > minio_down.sh
wait $MINIO_PID