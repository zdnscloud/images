# hub-helper
## hub-helper是什么
hub-helper是一个简单的dockerhub cli工具，用于在命令行下查找、删除自己dockerhub账户下镜像的一个工具
## 使用方法
使用工具之前需要先执行`hub-helper config`命令配置dockerhub用户名及密码，配置文件会保存在用户家目录下的`.hub-helper.json`中
> config命令user参数默认为zdnscloud，password参数为必填项
命令示例如下：`hub-helper config --user hiwyw --password Test@123`
* repo
    * list：列出账户下所有的仓库
    > 简写一级命令repos
    * info：打印某仓库详情，需在info命令后输入repo名称，如`hub-helper repo info singlecloud`
    * search：搜索账户下所有仓库，打印包含关键字的仓库名称，如`hub-helper repo search single`
* tag
    * list：列出某镜像的所有tag，如`hub-helper tag singlecloud`
    > 简写一级命令tags
    * info：打印tag详情，如`hub-helper tag info singlecloud master`
    * delete：删除镜像的指定tag，如`hub-helper tag delete singlecloud ws-test`
    > 简写一级命令delete
    * search：搜索镜像所有tag，打印包含关键字的tag名称，如`hub-helper tag search singlecloud tap`

