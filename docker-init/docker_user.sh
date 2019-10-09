#!/bin/bash
if [ ! -r ./ips ];then
	echo ""
	echo "you must touch 'ips' file and write in the host ip list"
	echo ""
	exit 1
fi

username="docker"
passwd="Zcloud!@#456"
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

install() {
	if [ ! -f ~/.ssh/id_rsa ];then
		sshkey
	fi
	for ip in $(cat ./ips);do
		ssh_copy_id ${ip}
	done
}

install
