#
# initialises Nexus Repository manager for use by an Artisan registry
#
CURRENTPASS=
echo "checking Nexus for temporary password file ..."
while [ -z "$CURRENTPASS" ]
do
  echo "password file not found - retrying in 5 seconds"
  sleep 5
  CURRENTPASS=$(docker exec nexus cat /nexus-data/admin.password)
done

echo wait for Nexus API
art curl -X GET \
  -a 25 \
  http://localhost:${bind=nexus-app:port}/service/rest/v1/status \
  -H 'accept: application/json'

echo "updating admin password from current temporary one"
art curl -X PUT \
  -u admin:${CURRENTPASS} \
  http://localhost:${bind=nexus-app:port}"/service/rest/v1/security/users/admin/change-password \
  -H 'accept: application/json','Content-Type: text/plain' \
  -d ${bind=nexus-app:var:NEXUS_ADMIN_PASSWORD}"

echo "creating new Artisan repository"
art curl -X POST \
  -u admin:${bind=nexus-app:var:NEXUS_ADMIN_PASSWORD} \
  http://localhost:${bind=nexus-app:port}/service/rest/v1/repositories/raw/hosted \
  -H 'accept: application/json','Content-Type: application/json' \
  -d '{
  "name": "artisan",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true,
    "writePolicy": "allow"
  },
  "cleanup": {
    "policyNames": [
      "string"
    ]
  },
  "component": {
    "proprietaryComponents": true
  },
  "raw": {
    "contentDisposition": "ATTACHMENT"
  }
}'

echo "disabling anonymous access"
art curl -X PUT \
  -u admin:${bind=nexus-app:var:NEXUS_ADMIN_PASSWORD} \
  http://localhost:${bind=nexus-app:port}/service/rest/v1/security/anonymous \
  -H 'accept: application/json','Content-Type: application/json' \
  -d '{"enabled": false}'
