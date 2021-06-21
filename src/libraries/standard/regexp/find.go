package main

import (
	"fmt"
	"os"
	"regexp"
	"unicode"
)

const parse = `
    getServoStatus = "getServoStatus"  # 1.1 获取机械臂伺服状态
    getMotorStatus = "getMotorStatus"  # 1.2 获取机械臂上下电状态
    setServoStatus = "set_servo_status"  # 1.3 设置机械臂伺服状态
    syncMotorStatus = "syncMotorStatus"  # 1.4 同步伺服编码器数据
    clearAlarm = "clearAlarm"  # 1.5 清除报警
    getMotorStatus = "getMotorStatus"  # 1.6 获取同步状态
    getRobotState = "getRobotState"  # 2.1 获取机器人状态
    getRobotMode = "getRobotMode"  # 2.2 获取机器人模式
    getRobotPos = "getRobotPos"  # 2.3 获取机器人当前位置信息
    getRobotPose = "getRobotPose"  # 2.4 获取机器人当前位姿信息
    getMotorSpeed = "getMotorSpeed"  # 2.5 获取机器人马达速度
    getCurrentCoord = "getCurrentCoord"  # 2.6 获取机器人当前坐标
    getCycleMode = "getServoStatus"  # 2.7 获取机器人循环模式
    getCurrentJobLine = "getCurrentJobLine"  # 2.8 获取机器人当前作业运行行号
    getCurrentEncode = "getCurrentEncode"  # 2.9 获取机器人当前编码器值列表
    getToolNumber = "getToolNumber"  # 2.10 获取机器人当前工具号
    getUserNumber = "getUserNumber"  # 2.11 获取机器人当前用户工具号
    getRobotTorques = "getRobotTorques"  # 2.12 获取机器人当前力矩信息
    getAnalogInput = "getAnalogInput"  # 2.13 获取模拟量输入
    setAnalogOutput = "setAnalogOutput"  # 2.14 设置模拟量输出
    dragTeachSwitch = "drag_teach_switch"  # 2.15 拖动示教开关
    setPayload = "cmd_set_payload"  # 2.16 设置机械臂负载和重心
    setToolCenterPoint = "cmd_set_tcp"  # 2.17 设置机械臂工具中心
    setToolNumber = "setToolNumber"  # 2.18 切换机器人当前工具号
    setUserNumber = "setUserNumber"  # 2.19 切换机器人用户工具号
    setCurrentCoord = "setCurrentCoord"  # 2.20 指定坐标系
    getCollisionState = "getCollisionState"  # 2.21 获取碰撞状态
    moveByJoint = "moveByJoint"  # 3.1 关节运动
    moveByLine = "moveByLine"  # 3.2 直线运动
    moveByArc = "moveByArc"  # 3.3 圆弧运动
    moveByRotate = "moveByRotate"  # 3.4 旋转运动
    stop = "stop"  # 3.5 停止机器人运行
    run = "run"  # 3.6 机器人自动运行
    pause = "pause"  # 3.7 机器人暂停
    addPathPoint = "addPathPoint"  # 3.8 添加路点信息2.0
    clearPathPoint = "clearPathPoint"  # 3.9 清除路点信息2.0
    moveByPath = "moveByPath"  # 3.10 轨迹运动2.0
    checkJbiExist = "checkJbiExist"  # 3.11 检查jbi文件是否存在
    runJbi = "runJbi"  # 3.12 运行jbi文件
    getJbiState = "getJbiState"  # 3.13 获取jbi文件运行状态
    getPathPointIndex = "getPathPointIndex"  # 3.14 获取jbi文件运行状态
    jog = "jog"  # 3.15 jog运动
    inverseKinematic = "inverseKinematic"  # 4.1 逆解函数
    positiveKinematic = "positiveKinematic"  # 4.2 正解函数
    convertPoseFromCartToUser = "convertPoseFromCartToUser"  # 4.3 基坐标到用户坐标位姿转化
    convertPoseFromUserToCart = "convertPoseFromUserToCart"  # 4.4 用户坐标到基坐标位姿转化
    inverseKinematicV2 = "inverseKinematic"  # 4.5 逆解函数2.0 带参考点位置逆解
    poseMul = "poseMul"  # 4.6 位姿相乘
    poseInv = "poseInv"  # 4.7 位姿求逆
    getInput = "getInput"  # 5.1 获取输入IO状态
    getOutput = "getOutput"  # 5.2 获取输出IO状态
    setOutput = "setOutput"  # 5.3 设置输出IO状态
    getVirtualInput = "getVirtualInput"  # 5.4 获取虚拟输入IO状态
    getVirtualOutput = "getVirtualOutput"  # 5.5 获取虚拟输出IO状态
    setVirtualOutput = "setVirtualOutput"  # 5.6 设置虚拟输出IO状态
    getSysVarB = "getSysVarB"  # 6.1 获取系统B变量值
    setSysVarB = "setSysVarB"  # 6.2 设置系统B变量值
    getSysVarI = "getSysVarI"  # 6.3 获取系统I变量值
    setSysVarI = "setSysVarI"  # 6.4 设置系统I变量值
    getSysVarD = "getSysVarD"  # 6.5 获取系统D变量值
    setSysVarD = "setSysVarD"  # 6.6 设置系统D变量值
    getSysVarPState = "getSysVarPState"  # 6.7 获取系统P变量是否启用
    getSysVarP = "getSysVarP"  # 6.8 获取系统P变量值
    getSysVarV = "getSysVarV"  # 6.9 获取系统V变量值
    ttInit = "transparent_transmission_init"  # 7.1 初始化透传服务
    ttSetCurrentServoJoint = "tt_set_current_servo_joint"  # 7.2 设置当前透传伺服目标关节点
    ttPutServoJointToBuf = "tt_put_servo_joint_to_buf"  # 7.3 添加透传伺服目标关节点信息到缓存中
    ttClearServoJointBuf = "tt_clear_servo_joint_buf"  # 7.4 清空透传缓存
    ttGetState = "get_transparent_transmission_state"  # 7.5 获取当前机器人是否处于透传状态
    getSoftVersion = "getSoftVersion"  # 10.2 获取控制器软件版本号
`

var template = `
    def %s(self) -> Any:
        '''%s'''
        # 注释
        return self.send_cmd(Method.%s)
`

func main() {
	dst := "dst.py"
	fw, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	reg := regexp.MustCompile(`(\w+) = ("\w+")\s*#\s*(.+)\n`)
	for _, val := range reg.FindAllStringSubmatch(parse, -1) {
		_, _ = fmt.Fprintf(fw, template, CamelCaseToUnderscore(val[1]), val[3], val[1])
	}
}

func CamelCaseToUnderscore(s string) string {
	var output []rune
	var previous rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) && !unicode.IsUpper(previous) {
			output = append(output, '_')
		}
		previous = r // 处理连续大写缩略词
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}
