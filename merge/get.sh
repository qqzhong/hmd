#!/bin/bash

function err() {
	echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $@" >&2
}

function log() {
	echo "[church]: $@"
}

password=$1
#Input
diffDir="JRD_PATCHS_20210301"
DATE=$(date "+%Y_%m_%d_%H%M%S")
mergeDir=merge_${DATE}

# !!!Note:the remotePath must be the full path of branch which ends with '/'!
remotePath="/local/int/milan.zhong/SM4350_ROW_SCW_DEV_0302/"
remotePathLen=${#remotePath}
#log "remotePathLen=${remotePathLen}"

#tree
find ${diffDir} -type f -printf "%P\n" > ${diffDir}.list
while read line;
do
	len=${remotePathLen}
	filePath=${remotePath}/${line}
	fileName=${filePath##*/}
	#log ${fileName}
	#log ${filePath}
	# ${filePath:0:${remotePathLen}}:it is remotePath.
	#log ${filePath:0:${len}}
	let len+=1
	target=${filePath:${len}}
	#log ${target}
	targetDir=${target%/*}
	#log ${targetDir}

	# !!!Note:adjust path here
	# But you'd better adjust ${diffDir} manually!
	if [ x"kernel/msm-5.4/drivers/input/touchscreen/focaltech_touch" = x"${targetDir}" ];then
		targetDir="kernel/msm-5.4/drivers/input/touchscreen/focaltech_thething"
	fi
	#log ${targetDir}

	# update the filePath
	filePath=${remotePath}/${targetDir}/${fileName}
	localPath=${mergeDir}/${targetDir}
	#log ${localPath}
	if [ ! -d ${localPath} ]; then
		mkdir -p ${localPath}
	fi

	sshpass -p ${password} scp android@10.231.12.110:${filePath}  ${localPath}
done < ${diffDir}.list

echo "${mergeDir} Done!"

sync;sync;sync
exit 0
