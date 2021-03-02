## 1.Use `get.sh hmd` to pull the fresh codes of `work branch` which will be merged with BSP patches.


Run get.sh on your local PC, then you will get some codes from the `work branch` that its skeleton is same as BSP patches.

### diffDir="JRD_PATCHS_20210301"
obviously, `diffDir`, it means patches path of `BSP manufacturer`.

### remotePath = '/local/int/milan.zhong/SM4350_ROW_SCW_DEV_0302/'
obviously, `remotePath`, it means android source path of `work branch`.

For example, you get folder "merge_2021_03_02_120429" from the `work branch`.


## 2. Use `Beyond Compare`, merge codes on your local PC.


**THIS IS tough work.**

![Beyond Compare](/data/data/picture/bc.png)


## 3. Use `set.sh password` to apply BSP patches to our `work branch` on Server.

**password is your PC password, so I don't know it.**
**compile and verify.**


For example, you apply the codes of folder "merge_2021_03_02_120429" to the `work branch`.


## 4. This is work flow.


![work flow](/data/data/picture/merge_workflow.png)

