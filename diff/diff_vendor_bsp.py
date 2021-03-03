#!/usr/bin/python3
# -*- coding: utf-8 -*-
# @filename           :  diff_vendor_bsp.py
# @author             :  church.zhong@hmdglobal.com
# @date               :  Mon Mar  1 09:58:22 HKT 2021
# @function           :  diff codes of vendor and bsp then output.
# @see                :  PEP 471
# @require            :  python3

import os
import re
import shutil
import time
import hashlib
from datetime import datetime

#python3 -m cProfile church/diff_vendor_bsp.py
#Punisher BSP分支0227代码:
#\\10.231.13.158\local\int\weikun.yan\code\Punisher_ROW_BSP_0227
#高通15基线代码:
#\\10.231.13.158\local\int\weikun.yan\code\Qcom_baseline_r1.0_r00015

#@Church Zhong
# Kernel\ 下除msm-5.4\techpack\camera以外的其它目录
# Device\
# Build\

# Watch out:global variables
bspPath = '/local/int/weikun.yan/code/Punisher_ROW_BSP_0227'
#bspPath = '/data/new'
bspPathLength = len(bspPath)
qcomPath = '/local/int/weikun.yan/code/Qcom_baseline_r1.0_r00015'
#qcomPath = '/data/church'
qcomPathLength = len(qcomPath)
# we're happy if android source are given.
androidSrc = ['build', 'device', 'kernel']
#androidSrc = ['shortcut-fe']
bspDict = {}
bspCompared = []
output = ''

def copy_file(name):
    target = output + os.sep + name
    targetDir = os.path.dirname(target)
    if not os.path.exists(targetDir):
        os.makedirs(targetDir);
    shutil.copyfile(bspPath + os.sep + name, target)

def scan_android(path, func):
    """Scan android"""
    for entry in os.scandir(path):
        if entry.is_dir(follow_symlinks=False):
            if entry.name.startswith('.git'):
                #print('Ignore={}'.format(entry.path))
                pass
            elif entry.name.startswith('.svn'):
                #print('Ignore={}'.format(entry.path))
                pass
            elif entry.name=='CVS':
                #print('Ignore={}'.format(entry.path))
                pass
            else:
                scan_android(entry.path, func)
        elif entry.is_dir(follow_symlinks=True):
            #print('Ignore={}'.format(entry.path))
            pass
        elif entry.is_file():
            if not os.path.exists(entry.path):
                #print('File not exist={}'.format(entry.path))
                pass
            elif ''==entry.name:
                #print('File name empty={}'.format(entry.path))
                pass
            else:
                func(entry)
        else:
            print('FIXME={}'.format(entry.path))
            pass

def bsp_func(entry):
    global bspPathLength
    global bspDict
    name = str(entry.path)
    name = name[bspPathLength:]
    md5sum = hashlib.md5(open(entry.path,'rb').read()).hexdigest()
    #size = entry.stat(follow_symlinks=False).st_size
    #row = '{:0>32d},{},{},{},{}'.format(size, name, size, 0, md5sum)
    #print('{0},{1}'.format(name, md5sum))
    bspDict[name] = md5sum

def scan_bsp(path):
    scan_android(path, bsp_func)

def qcom_func(entry):
    global bspPath
    global qcomPathLength
    global bspDict
    global bspCompared
    name = str(entry.path)
    name = name[qcomPathLength:]
    md5sum = hashlib.md5(open(entry.path,'rb').read()).hexdigest()
    #size = entry.stat(follow_symlinks=False).st_size
    #row = '{:0>32d},{},{},{},{}'.format(size, name, size, 0, md5sum)
    #print('{0},{1}'.format(name, md5sum))

    if name in bspDict:
        bspCompared.append(name)
        if bspDict[name]==md5sum:
            pass
        else:
            copy_file(name)
            print(name)

def scan_qcom(path):
    scan_android(path, qcom_func)

def work(path):
    global bspPath
    global qcomPath
    global output

    print('bsp={}'.format(bspPath))
    print('qcom={}'.format(qcomPath))

    moment = datetime.now().strftime("%Y%m%d%H%M%S")
    output = os.getcwd() + os.sep + 'diff_' + moment
    os.mkdir(output)

    for f in androidSrc:
        scan_bsp(bspPath + os.sep + f)

    for f in androidSrc:
        scan_qcom(qcomPath + os.sep + f)

    # BSP added something
    for key in bspCompared:
        bspDict[key] = 'compared'
        pass
    for key, value in bspDict.items():
        if value=='compared':
            pass
        else:
            copy_file(key)
            print('BSP added={}'.format(key))

def main():
    start = time.time()
    work(os.getcwd())
    print('running time:%s' % (time.time() - start))

if __name__ == '__main__':
    main()
