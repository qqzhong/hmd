#!/bin/bash
# @filename           :  query_gerrit.sh
# @author             :  Copyright (C) Church.Zhong
# @date               :  Thu Mar  8 10:00:26 HKT 2018
# @function           :  query gerrit database by ssh.
# @see                :  
# @require            :  OpenSSH_7.6p1 Ubuntu-4ubuntu0.3, OpenSSL 1.0.2n  7 Dec 2017
SECONDS=0
EX_OK=0
EX_USAGE=64

err() {
	echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $@" >&2
}

#==================Config=======================================
date="2021-03-05";
username="church.zhong";
bspBranch="sm4350_TF_AKT_BSP";
devBranch="sm4350_TF_AKT_DEV";
repoPath="/data/church/";
#===============================================================

echo "date=${date}";
echo "username=${username}";
echo "bspBranch=${bspBranch}";
echo "devBranch=${devBranch}";
echo "repoPath=${repoPath}";

#==================Main=========================================
bspBranchJson=${bspBranch}".json";
devBranchJson=${devBranch}".json";
echo "bspBranchJson=${bspBranchJson}";
echo "devBranchJson=${devBranchJson}";

query="ssh -p 29418 ${username}@hmdgerritserver.hmdglobal.com gerrit query --format=json --current-patch-set"
${query} "(status:merged AND branch:"${bspBranch}") AND after:"${date} > ${bspBranchJson};
sync;
${query} "(status:merged AND branch:"${devBranch}") AND after:"${date} > ${devBranchJson};
sync;

./diff_gerrit ${date} -username ${username} -bsp ${bspBranch} -dev ${devBranch} -repo ${repoPath};

echo "Done";
# do some work( or time yourscript.sh)
duration=$SECONDS
echo "$(($duration / 60)) minutes and $(($duration % 60)) seconds elapsed."
exit ${EX_OK}
#================End============================================