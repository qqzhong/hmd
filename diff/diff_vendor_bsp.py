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
# Watch out:global variables
bspPath = '/local/int/weikun.yan/code/Punisher_ROW_BSP_0303/'
bspPathLength = len(bspPath)
qcomPath = '/local/int/weikun.yan/code/Qcom_baseline_r015/'
qcomPathLength = len(qcomPath)
# we're happy if android source are given.
androidSrc = ['bootable','build','device','frameworks','hardware','kernel','vendor']
bspDict = {}
# Output:what are BSP modified.
bspModifiedOutput = ''
# Output:what are Qcom baseline codes.
qcomBaselineOutput = ''

def copy_file(srcPath, destPath, name):
    target = destPath + os.sep + name
    targetDir = os.path.dirname(target)
    if not os.path.exists(targetDir):
        os.makedirs(targetDir);
    shutil.copyfile(srcPath + os.sep + name, target)
    print('copy:{}'.format(target))

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
    global bspModifiedOutput
    global qcomBaselineOutput

    fullName = str(entry.path)
    filePath = os.path.dirname(fullName)
    fileName = fullName[qcomPathLength:]
    md5sum = hashlib.md5(open(entry.path,'rb').read()).hexdigest()
    #size = entry.stat(follow_symlinks=False).st_size
    #row = '{:0>32d},{},{},{},{}'.format(size, fileName, size, 0, md5sum)
    #print('{0},{1}'.format(fileName, md5sum))

    if fileName in bspDict:
        if bspDict[fileName]==md5sum:
            pass
        else:
            copy_file(bspPath, bspModifiedOutput, fileName)
            copy_file(qcomPath, qcomBaselineOutput, fileName)
        bspDict[fileName] = None

def scan_qcom(path):
    scan_android(path, qcom_func)

def work(path):
    global bspPath
    global qcomPath
    global bspModifiedOutput
    global qcomBaselineOutput

    print('bsp={}'.format(bspPath))
    print('qcom={}'.format(qcomPath))

    moment = datetime.now().strftime("%Y_%m%d_%H%M%S")
    bspModifiedOutput = '{0}{1}{2}{3}{4}{5}'.format(os.getcwd(),os.sep,'extract_',moment,os.sep,'bspModified')
    os.makedirs(bspModifiedOutput)
    qcomBaselineOutput = '{0}{1}{2}{3}{4}{5}'.format(os.getcwd(),os.sep,'extract_',moment,os.sep,'qcomBaseline')
    os.makedirs(qcomBaselineOutput)

    for f in androidSrc:
        scan_bsp(bspPath + os.sep + f)

    for f in androidSrc:
        scan_qcom(qcomPath + os.sep + f)

    # BSP added something
    for key, value in bspDict.items():
        if value:
            copy_file(bspPath, bspModifiedOutput, key)
            print('BSP added={}'.format(key))

def main():
    start = time.time()
    work(os.getcwd())
    print('running time:%s' % (time.time() - start))

if __name__ == '__main__':
    main()
