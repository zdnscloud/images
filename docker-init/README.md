> 此脚本用于为批量主机部署免密登录和docker环境

##  前提要求
- 使用root执行
- 所有主机root密码一致
- 在当前目录创建 ips 文件，将主机ip列表写入
- 主机操作系统为CentOS 7.6 或者ubuntu 18.04

### 使用

执行脚本manager.sh，执行后会提示输入root密码

### 结果
- 安装并启动docker-18.06.3-ce
- 创建了docker用户，密码为docker
- docker用户的密钥对路径:/home/docker/.ssh/id_rsa
