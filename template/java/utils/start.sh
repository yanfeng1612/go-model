#!/bin/bash

#----------以下变量可以在脚本中直接使用----------------
# $def_app_id           应用id
# $def_app_name         应用名称
# $def_app_domain       应用域名
# $def_app_deploy_path  兼容老部署,废弃
# $def_path_app_log     应用日志路径 如：/export/Logs/
# $def_path_app_data    如：/export/Data
# $def_group_id         分组 id
# $def_instance_id      实例id
# $def_instance_name    实例名称 server1
# $def_instance_path    实例完成路径 /export/Instances/jone/server1/
# $def_host_ip
#--------------------------

#set -o errexit
readonly APP_NAME="demo"          #定义当前应用的名称
#获取当前应用的进程 id
function get_running_pid {
  pgrep -f "$EXE_JAR"
}

BASEDIR=$(cd $(dirname $0) && pwd)/../lib        # 获取运行脚本的上级目录

readonly JAVA_HOME="/export/servers/jdk1.8.0_60" # java home
readonly JAVA="$JAVA_HOME/bin/java"

EXE_JAR="demo.jar"
echo "$EXE_JAR"
OPTS_MEMORY="-server -Xms4G -Xmx4G -Xmn1G -Xss256K -XX:+UseG1GC -Xloggc:/export/Logs/demo/gc.log -XX:+PrintGCDetails -XX:+PrintGCTimeStamps"       #定义启动 jvm 的参数信息。
CLASSPATH="$BASEDIR/"

[[ -z $(get_pid) ]] || {
    echo "ERROR:  $APP_NAME already running" >&2
    exit 1
}

echo "Starting $APP_NAME ...."
[[ -x $JAVA ]] || {
    echo "ERROR: no executable java found at $JAVA" >&2
    exit 1
}
cd $BASEDIR
#setsid "$JAVA" $OPTS_MEMORY -jar "$EXE_JAR" "$@" > /dev/null 2>&1 &
java $OPTS_MEMORY -jar "$EXE_JAR" "$@" > /dev/null 2>&1 &

sleep 0.5
[[ -n $(get_pid) ]] || {
    echo "ERROR: $APP_NAME failed to start" >&2
    exit 1
}

echo "$APP_NAME is up runnig :)"
