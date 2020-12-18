# check for input variables
if [[ -z "${ARTEFACT_NAME+x}" ]]; then
    echo "ARTEFACT_NAME must be provided"
    exit 1
fi
if [[ -z "${ARTEFACT_REG_USER+x}" ]]; then
    echo "ARTEFACT_REG_USER must be provided"
    exit 1
fi
if [[ -z "${ARTEFACT_REG_PWD+x}" ]]; then
    echo "ARTEFACT_REG_PWD must be provided"
    exit 1
fi
if [[ -z "${BUILD_PROFILE+x}" ]]; then
    echo "BUILD_PROFILE must be provided"
    exit 1
fi
# import the signing key from the mount point
art key import -k /keys/root_rsa_key.pgp
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to import signing key"
   exit 1
fi

# import the verification key from the mount point
art key import /keys/root_rsa_pub.pgp
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to import verification key"
   exit 1
fi

# build the artefact
art build -t "${ARTEFACT_NAME}" -p "${BUILD_PROFILE}" .
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to build artefact"
   exit 1
fi

# push the built artefact
art push -u "${ARTEFACT_REG_USER}:${ARTEFACT_REG_PWD}" "${ARTEFACT_NAME}"
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to push artefact"
   exit 1
fi