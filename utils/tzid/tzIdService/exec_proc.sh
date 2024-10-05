#!/bin/bash

action=$1

date_Ymd=`date +%F`

stop_proc(){
    local cpwd=$(cd $(dirname $0); pwd)
    local proc_name=`echo ${cpwd} |awk -F/ '{print $NF}'`
    local pid=`ps -ef |grep /${proc_name}/|grep -v grep |grep -v bash |awk '{print $2}'`
    if [ "${pid}" == "" ];then
        echo "${proc_name}进程不存在"
        exit 1
    else
        ps -ef |grep /${proc_name}/|grep -v grep |grep -v bash |awk '{print $2}' |xargs kill
        if [ $? -eq 0 ];then
            echo "${proc_name}进程已停止"
        else
            echo "${proc_name}进程停止失败"
            exit 1
        fi

    fi
}

start_proc(){
    local cpwd=$(cd $(dirname $0); pwd)
    local proc_name=`echo ${cpwd} |awk -F/ '{print $NF}'`
    
    
	chmod +x ${proc_name}
    local pid=`ps -ef |grep /${proc_name}/|grep -v grep |grep -v bash |awk '{print $2}'`
    if [ "${pid}" != "" ];then
        echo "${proc_name}进程已存在"
        exit 1
    else
        logs_dir="${cpwd}/logs"
        cd ${cpwd}
        logs_file="${logs_dir}/${proc_name}_${date_Ymd}.log"
        if [ ! -d ${logs_dir} ]; then
            mkdir ${logs_dir}
        fi
        nohup ${cpwd}/${proc_name} >>${logs_file} 2>&1 &
        if [ $? -eq 0 ];then
            echo "${proc_name}进程已启动"
        else
            echo "${proc_name}进程启动失败"
            exit 1
        fi
    fi
}

case $action in
    "stop")
        stop_proc
        ;;
    "start")
        start_proc
        ;;
    "restart")
        stop_proc
        start_proc
        ;;
    *)
        echo "未知参数"
        ;;
esac