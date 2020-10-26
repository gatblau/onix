curl -X POST "http://localhost:8081/service/rest/v1/components?repository=gatblau" \
  -H "accept: application/json" \
  -H "Content-Type: multipart/form-data" \
  -F "raw.directory=boot-snapshot" \
  -F "raw.asset1=@261020123238270-5113f5c110.json;type=application/json" \
  -F "raw.asset1.filename=261020123238270-5113f5c110.json" \
  -F "raw.asset2=@261020123238270-5113f5c110.zip;type=application/zip" \
  -F "raw.asset2.filename=261020123238270-5113f5c110.zip"