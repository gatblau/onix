# fetch user home directory from /etc/passwd file
export USER_HOME=$(awk -F":" '{print $6}' /etc/passwd | grep -m1 `whoami`)
#mkdir -p ${USER_HOME}/keys
#cp -R /keys ${USER_HOME}
cp /keys/.pilot_sign.pgp /app/.
/app/pilotctl