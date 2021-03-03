#!/usr/bin/python3
# -*- coding: utf-8 -*-
# @filename           :  count_files.py
# @author             :  church.zhong@hmdglobal.com
# @date               :  Wed Mar  3 09:14:01 HKT 2021
# @function           :  count the type of all files in diff folder.
# @see                :  PEP 471
# @require            :  python3

import os
import time
from datetime import datetime

# Watch out:global variables
diffPath = '/local/int/church/diff_20210302195056'
diffPathLength = len(diffPath)
mimes = {}

def _scan(path, func):
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
                _scan(entry.path, func)
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

def _count(entry):
    global diffPath
    global mimes
    name = str(entry.path)
    #name = name[diffPathLength:]
    #md5sum = hashlib.md5(open(entry.path,'rb').read()).hexdigest()
    #size = entry.stat(follow_symlinks=False).st_size
    #row = '{:0>32d},{},{},{},{}'.format(size, name, size, 0, md5sum)
    xtuple = os.path.splitext(name)

    if 0==len(xtuple[1]):
        filename = os.path.basename(xtuple[0])
        if filename in mimes:
            mimes[filename] += 1
        else:
            mimes[filename] = 1
    else:
        if xtuple[1] in mimes:
            mimes[xtuple[1]] += 1
        else:
            mimes[xtuple[1]] = 1

def work(path):
    global diffPath
    global mimes
    #print('diffPath={}'.format(diffPath))
    _scan(diffPath, _count)
    for key, value in mimes.items():
        print('{}\t{}'.format(key, value))

def main():
    start = time.time()
    work(os.getcwd())
    print('running time:%s' % (time.time() - start))

if __name__ == '__main__':
    main()
