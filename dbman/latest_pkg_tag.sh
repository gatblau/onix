HASH=$(git rev-parse --short HEAD)
TIME=$(git show -s --format=%cI)

# shellcheck disable=SC2046
echo $TIME +%y%m%d


#2020-06-20T07:43:33+01:00
#v0.0.0-20200619171616-f89c32144be2