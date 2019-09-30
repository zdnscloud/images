#!/bin/bash

username=$(whoami)
if [ $username != "root" ];then
	echo ""
	echo "you must use root to execute this script"
	echo ""
	exit 1
fi

if [ ! -r ./ips ];then
	echo ""
	echo "you must touch 'ips' file and write in the host ip list"
	echo ""
	exit 1
fi

read -p "Enter root password:" passwd
if [ -z $passwd ];then
	echo ""
	echo "password can not be null"
	echo ""
	exit 1
fi

ssh_args="-o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

sshkey() {
/usr/bin/expect <<EOF
set time 10
spawn ssh-keygen -t rsa
expect {
"*id_rsa):" {
send "\n";
exp_continue
}
"*(y/n)?" {
send "y\n"
exp_continue
}
"*passphrase):" {
send "\n"
exp_continue
}
"*again:" {
send "\n"
}
}
expect eof
EOF
}

ssh_copy_id() {
address=$1
/usr/bin/expect <<EOF
set time 30
spawn ssh-copy-id ${ssh_args} ${username}@${address}
expect {
"*yes/no" { send "yes\r"; exp_continue}
"*password:" {send "$passwd\r"}
}

expect eof
EOF
}

get_distribution() {
        lsb_dist=""
        if [ -r /etc/os-release ]; then
                lsb_dist="$(. /etc/os-release && echo "$ID")"
        fi
	lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"
	echo $lsb_dist
}


install_expect() {
	case "$lsb_dist" in
	centos)
		yum install -y -q expect >/dev/null
		;;
	ubuntu)
		apt-get install -y -qq expect >/dev/null
		;;
	esac
}

install() {
	lsb_dist=$( get_distribution )
	if [ ! -x /usr/bin/expect ];then
		install_expect
	fi
	if [ ! -f ~/.ssh/id_rsa ];then
		sshkey
	fi
	for ip in $(cat ./ips);do
		ssh_copy_id ${ip}
		scp ${ssh_args} ./docker_install.sh ${ip}:~/
		ssh ${ssh_args}  ${ip} "sh docker_install.sh"
	done
}

docker_user() {
	cp docker_user.sh  ips /home/docker/
	su - docker -c 'sh docker_user.sh'
}

install
docker_user
