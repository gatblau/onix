curl -u admin:admin123 -X POST "http://localhost:8081/service/rest/v1/repositories/raw/hosted" \
  -H "accept: application/json" \
  -H "Content-Type: application/json" \
  -d "{ \"name\": \"gatblau\", \"online\": true, \"storage\": { \"blobStoreName\": \"default\", \"strictContentTypeValidation\": true, \"writePolicy\": \"allow_once\" }, \"cleanup\": { \"policyNames\": [ \"string\" ] }, \"raw\": { \"contentDisposition\": \"ATTACHMENT\" }}"