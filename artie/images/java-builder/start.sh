artie build -t "${ARTEFACT_NAME}" -p "${BUILD_PROFILE}" .
artie push -u "${ARTEFACT_UNAME}:${ARTEFACT_PWD}" "${ARTEFACT_NAME}"