if [[ -z "${SONAR_PROJECT_KEY+x}" ]]; then
    echo "SONAR_PROJECT_KEY must be provided"
    exit 1
fi

if [[ -z "${SONAR_URI+x}" ]]; then
    echo "SONAR_URI must be provided"
    exit 1
fi

if [[ -z "${SONAR_SOURCES+x}" ]]; then
    echo "SONAR_SOURCES must be provided"
    exit 1
fi

if [[ -z "${SONAR_BINARIES+x}" ]]; then
    echo "SONAR_BINARIES must be provided"
    exit 1
fi

if [[ -z "${SONAR_TOKEN+x}" ]]; then
    echo "SONAR_TOKEN must be provided"
    exit 1
fi

# merge sonar URI variable
printf "sonar.sourceEncoding=UTF-8\nsonar.host.url=%s\n" "${SONAR_URI}" >> /sonar/conf/sonar-scanner.properties

# initiates the scan process
sonar-scanner -Dsonar.projectKey="${SONAR_PROJECT_KEY}" -Dsonar.sources="${SONAR_SOURCES}" -Dsonar.java.binaries="${SONAR_BINARIES}" -Dsonar.login="${SONAR_TOKEN}"

# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
    printf "failed to run the Sonar scanner\n"
    exit 1
fi

printf "sonar server is %s\n" "${SONAR_URI}"
printf "project key is %s\n" "${SONAR_PROJECT_KEY}"

# retrieve analysis identifier from Sonar
analysis_id=$(curl -s "${SONAR_URI}/api/project_analyses/search?project=${SONAR_PROJECT_KEY}" | jq '.analyses[]' | jq '.key' | head -n 1)
analysis_id=$(echo "${analysis_id}" | tr -d "\"\`'")

printf "analysis is %s\n" "${analysis_id}"

# using the identifier retrieves the analysis status
analysis_status=$(curl -s "${SONAR_URI}/api/qualitygates/project_status?analysisId=${analysis_id}" | jq '.projectStatus.status')

printf "analysis status is %s\n" "${analysis_status}"

if [[ "${analysis_status}" = \""OK\"" ]]; then
  printf "Sonar Code Quality check => Pass\n"
  exit 0
fi

printf "Sonar Code Quality check => Failed\n"
exit 1
