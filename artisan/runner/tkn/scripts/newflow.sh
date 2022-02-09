# use this sript to post an Artisan Flow specification to Onix

# define the location of the Onix service
ONIX_URI=http://localhost:8080

# ensure the credentials are correct to authenticate to Onix
USER=admin
PWD=ybGkxQ7GTVu3ek9eVQQM

# define the flow Key in Onix
FLOW_KEY=testflow

# post the flow defined in flowspec.json to Onix
# update flowspec.json using your specific flow spec within the "meta" attribute in JSON format
curl -X PUT "${ONIX_URI}/item/${FLOW_KEY}" \
     -u "${USER}:${PWD}" \
     -H  "accept: application/json" \
     -H  "Content-Type: application/json" \
     -d "@flowspec.json"

# to get the stored flow you can execute the following command
#  curl ${ONIX_URI}/item/${FLOW_KEY} \
#       -u "${USER}:${PWD}" \
#       -H  "accept: application/json"
