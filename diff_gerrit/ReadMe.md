## 1.安装 ssh 客户端

https://blog.csdn.net/wm609972715/article/details/83759114


## 2.配置 ssh 公钥和密钥，注意公钥和密钥千万不要填写错误

在用户目录下新建.ssh目录，把ssh的公钥（id_rsa_pub）和秘钥文件(id_rsa)放进去，可以实现免密登录
https://blog.csdn.net/qq_43901693/article/details/103700272​​​​​​​


## 3.测试ssh在windows系统上运行是否正常


修改下面的命令，把steve.jobs换成你的名字，然后在ssh客户端运行一下
```c
ssh -p 29418 steve.jobs@hmdgerritserver.hmdglobal.com gerrit query --format=json '(owner:'church.zhong@hmdglobal.com' AND status:merged AND branch:sm4350_TF_AKT_DEV) AND after:'2021-03-01'' > church.json
```


文件church.json内容可能是：
```json
{"project":"HMD/platform/vendor/opensource/audio-kernel","branch":"sm4350_TF_AKT_DEV","id":"I3dcca06f99ce25c06cf79803f84b4a913ee6ebe1","number":40868,"subject":"\u003c48265\u003e\u003c1_1\u003e\u003cfengpengbo\u003e\u003caoki\u003e\u003cbsp\u003e\u003cNA\u003eporting changes of scw to aoki","owner":{"email":"church.zhong@hmdglobal.com","username":"church.zhong"},"url":"http://hmdgerritserver.southeastasia.cloudapp.azure.com/c/HMD/platform/vendor/opensource/audio-kernel/+/40868","commitMessage":"\u003c48265\u003e\u003c1_1\u003e\u003cfengpengbo\u003e\u003caoki\u003e\u003cbsp\u003e\u003cNA\u003eporting changes of scw to aoki\n\nSigned-off-by: admin \u003cadmin@ontim.cn\u003e\nChange-Id: I3dcca06f99ce25c06cf79803f84b4a913ee6ebe1\n","createdOn":1615187068,"lastUpdated":1615187088,"open":false,"status":"MERGED"}
{"project":"HMD/amss/trustzone_images","branch":"sm4350_TF_AKT_DEV","id":"Id2df901e1d304421441b4ebd0a2ee476553c9e21","number":40867,"subject":"\u003c48427\u003e\u003c2_2\u003e\u003cshangfei\u003e\u003caoki\u003e\u003cbsp\u003e\u003cNA\u003eupdate ca \u0026 ta for new holitech fingerprint module","owner":{"email":"church.zhong@hmdglobal.com","username":"church.zhong"},"url":"http://hmdgerritserver.southeastasia.cloudapp.azure.com/c/HMD/amss/trustzone_images/+/40867","commitMessage":"\u003c48427\u003e\u003c2_2\u003e\u003cshangfei\u003e\u003caoki\u003e\u003cbsp\u003e\u003cNA\u003eupdate ca \u0026 ta for new holitech fingerprint module\n\nSigned-off-by: admin \u003cadmin@ontim.cn\u003e\nChange-Id: Id2df901e1d304421441b4ebd0a2ee476553c9e21\n","createdOn":1615186844,"lastUpdated":1615186857,"open":false,"status":"MERGED"}
{"project":"HMD/device/nokia","branch":"sm4350_TF_AKT_DEV","id":"I134f5fdf4aa030539e6190ab2a727352b349dd70","number":40864,"subject":"\u003c48427\u003e\u003c2_1\u003e\u003cshangfei\u003e\u003caoki\u003e\u003cbsp\u003e\u003cNA\u003eupdate ca \u0026 ta for new holitech fingerprint module","owner":{"email":"church.zhong@hmdglobal.com","username":"church.zhong"},"url":"http://hmdgerritserver.southeastasia.cloudapp.azure.com/c/HMD/device/nokia/+/40864","commitMessage":"\u003c48427\u003e\u003c2_1\u003e\u003cshangfei\u003e\u003caoki\u003e\u003cbsp\u003e\u003cNA\u003eupdate ca \u0026 ta for new holitech fingerprint module\n\nSigned-off-by: admin \u003cadmin@ontim.cn\u003e\nChange-Id: I134f5fdf4aa030539e6190ab2a727352b349dd70\n","createdOn":1615186391,"lastUpdated":1615186430,"open":false,"status":"MERGED"}
{"type":"stats","rowCount":3,"runTimeMilliseconds":9,"moreChanges":false}
```


## 4.在Windows的 cmd 或  powershell 终端运行脚本

