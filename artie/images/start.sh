# import the signing key from the mount point
artie key import -k /keys/root_rsa_key.pgp
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to import signing key"
   exit 1
fi

# import the verification key from the mount point
artie key import /keys/root_rsa_pub.pgp
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to import verification key"
   exit 1
fi

# build the artefact
artie build -t "${ARTEFACT_NAME}" -p "${BUILD_PROFILE}" .
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to build artefact"
   exit 1
fi

# push the built artefact
artie push -u "${ARTEFACT_UNAME}:${ARTEFACT_PWD}" "${ARTEFACT_NAME}"
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to push artefact"
   exit 1
fi