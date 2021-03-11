#!/bin/bash
# @filename           :  wifi_log.sh
# @author             :  Copyright (C) Church.Zhong
# @date               :  Thu Mar 11 14:56:46 HKT 2021
# @function           :  parse Wi-Fi logcat&kernel.
# @see                :  
# @require            :  GNU bash, version 4.4.20
SECONDS=0
EX_OK=0
EX_USAGE=64

#get absolute path of a path name
abspath() {
	[[ -n  "$1" ]] && ( cd "$1" 2>/dev/null && pwd;)
}

err() {
	echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $@" >&2
}

OS_DATE_DAY=$(date +%Y%m%d)
OS_DATE_SECOND=$(date +%Y%m%d%H%M%S)

function main() {
	dir=$@
	if [ ! -d "${dir}" ];then
		err "${dir}: No such file or directory"
		exit ${EX_USAGE}
	fi

	DIR=church_log_${OS_DATE_SECOND}
	mkdir -p ${DIR}

	for f in $(find ${dir} -maxdepth 1 -type f -printf "%f\n")
	do
		echo ${f}
		if [[ ${f} == kernel* ]];then
			kernel=${DIR}/kernel.log
			sed -n "/wlan\|wifi/Iw ${kernel}" ${dir}/${f}
		fi

		if [[ ${f} == logcat* ]];then
			android_hardware_wifi=${DIR}/android_hardware_wifi.log
			sed -n "/android.hardware.wifi/Iw ${android_hardware_wifi}" ${dir}/${f}

			wificond=${DIR}/wificond.log
			sed -n "/wificond/Iw ${wificond}" ${dir}/${f}

			connectivityService=${DIR}/connectivityService.log
			sed -n "/ConnectivityService/Iw ${connectivityService}" ${dir}/${f}

			cnss_daemon=${DIR}/cnss_daemon.log
			sed -n "/cnss-daemon/Iw ${cnss_daemon}" ${dir}/${f}

			wpa_supplicant=${DIR}/wpa_supplicant.log
			sed -n "/wpa_supplicant/Iw ${wpa_supplicant}" ${dir}/${f}

			lowi=${DIR}/LOWI-9.0.0.75.log
			sed -n "/LOWI-9.0.0.75/Iw ${lowi}" ${dir}/${f}

			ethernetTracker=${DIR}/EthernetTracker.log
			sed -n "/EthernetTracker/Iw ${ethernetTracker}" ${dir}/${f}

		fi
	done

	sync;sync;sync
}

main $@

# do some work( or time yourscript.sh)
duration=$SECONDS
echo "$(($duration / 60)) minutes and $(($duration % 60)) seconds elapsed."

exit ${EX_OK}
