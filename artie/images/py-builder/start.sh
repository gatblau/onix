artie build -t ${ARTEFACT_NAME} -p ${BUILD_PROFILE} .
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "artie build faild to execute..."
   exit -1
fi

artie push -u ${ARTEFACT_UNAME}:${ARTEFACT_PWD} ${ARTEFACT_NAME}
# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "artie push faild to execute..."
   exit -1
fi
