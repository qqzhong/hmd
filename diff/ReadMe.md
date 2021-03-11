## 1.Modify 3 variables in file `diff_vendor_bsp.py`.


### bspPath = '/local/int/weikun.yan/code/Punisher_ROW_BSP_0227'
obvious, `bspPath`, it means android source path of `BSP manufacturer`.

### qcomPath = '/local/int/weikun.yan/code/Qcom_baseline_r1.0_r00015'
obvious, `qcomPath`, it means android source path of `QUALCOMM vendor`.

### androidSrc = ['build', 'device', 'kernel']
obvious, `androidSrc`, it means android source directories that `BSP manufacturer` modified.


## 2. run `diff_vendor_bsp.py`  or `run diff_vendor_bsp.py > diff.log`.


## 3.Output example.


`Modified: shortcut-fe/README`
`Added: shortcut-fe/xx.txt`

```c
android@Church:/data/tmp$ tree extract_2021_0304_101022/
extract_2021_0304_101022/
|-- bspModified
|   `-- shortcut-fe
|       |-- README
|       `-- xx.txt
`-- qcomBaseline
    `-- shortcut-fe
        `-- README

4 directories, 3 files
```

## 4. bug report to `church.zhong@hmdglobal.com`.

