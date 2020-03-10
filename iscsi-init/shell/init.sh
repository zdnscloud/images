#!/bin/bash
if [ -z "${TARGET_HOST}" ];then
  echo "ERROR- You must provide target host address."
  exit 1
fi
if [ -z "${TARGET_PORT}" ];then
  echo "ERROR- You must provide target port."
  exit 1
fi
if [ -z "${TARGET_IQN}" ];then
  echo "ERROR- You must provide target iqn."
  exit 1
fi
if [ -z "${VOLUME_GROUP}" ];then
  echo "ERROR- You must provide volume group to init."
  exit 1
fi


chap() {
if [ -s "/root/secret/username" -a -s "/root/secret/password" ];then
    username=$(cat /root/secret/username)
    password=$(cat /root/secret/password)
    iscsiadm -m node -T ${TARGET_IQN} -o update --name node.session.auth.authmethod --value=CHAP
    iscsiadm -m node -T ${TARGET_IQN} -o update --name node.session.auth.username --value=${username}
    iscsiadm -m node -T ${TARGET_IQN} -o update --name node.session.auth.password --value=${password}
fi
}

discovery() {
    iscsiadm -m discovery -t sendtargets -p ${TARGET_HOST}:${TARGET_PORT}
}

login() {
    iscsiadm -m node -T ${TARGET_IQN} -p ${TARGET_HOST}:${TARGET_PORT} -l
}

logout() {
    iscsiadm -m node -T ${TARGET_IQN} -p ${TARGET_HOST}:${TARGET_PORT} -u
}

init() {
  Device=$(lsscsi -t -L|grep disk|grep ${TARGET_IQN} |awk '{print $NF}')
  if [[ "${Device}" == "/dev/*" ]];then
    echo "ERROR- Can not find iscsi disk"
    exit 1
  fi
  pvs|grep ${Device} -q
  if [ $? -ne 0 ];then
    pvcreate ${Device}
  fi
  vgs|grep  ${VOLUME_GROUP} -q
  if [ $? -ne 0 ];then
    vgcreate ${VOLUME_GROUP} ${Device}
  fi
}

discovery
chap
login
time=$(date +%s%N | md5sum | head -c 1)
sleep ${time}
init
if [ $? -ne 0 ];then
	exit 99
fi
tail -f /dev/null
