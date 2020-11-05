#!/bin/bash
# AgOcean.InfraInfo/resources/infra-app/ag-ocean/base-app/cloud/accountsystem
K8sBaseKustomizationPath=$1 
# 1.1.1.1
IMAGEID=$2
# agocean.accountsystem.1.1.1.1
newRepoPath=$4.$3.$2
# tmp/agocean
gameName=$5/$4
# agocean.accountsystem.1.1.1.1.zip
zipName=$newRepoPath$6
# 更新infra儲存庫
git pull -C AgOcean.InfraInfo pull
# 建立暫存區
mkdir -p $gameName
# 複製infra目標yaml到暫存區
cp -Rf $K8sBaseKustomizationPath $gameName/$newRepoPath
# 替換VersionID
ls $gameName/$newRepoPath | while read -r filename; do sed -i "s/VERSION_ID/${IMAGEID}/g" $gameName/$newRepoPath/$filename; done
cd $gameName
# 打包yaml檔
zip $zipName $newRepoPath/ -r -m -j
