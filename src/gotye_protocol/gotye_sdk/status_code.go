package gotye_sdk

const (
	//system
	API_SUCCESS                    = 200
	API_INVALID_TOKEN_ERROR        = 401 //认证失败，无效token
	API_INVALID_LICENSE_ERROR      = 411 //无效许可
	API_NO_PERMISSION_ACCESS_ERROR = 420 //没有权限访问此资源
	API_SYSTEM_ERROR               = 500 //系统处理异常

	//live room
	API_INVALID_LIVEROOM_ID_ERROR       = 2001 //无效的主播室ID
	API_NOT_EXISTS_LIVEROOM_ID_ERROR    = 2002 //主播室不存在
	API_REPECT_PASSWORD_LIVEROOM_ERROR  = 2003 //主播室密码
	API_INVALID_PASSWORD_LIVEROOM_ERROR = 2004 //主播室密码非法
	API_INVALID_LIVEROOM_NAME_ERROR     = 2005 //主播室名称非法
	API_NULL_LIVEROOM_ID_ERROR          = 2007 //主播室ID为空
	API_EXISTS_THIRD_LIVEROOM_ID_ERROR  = 2008 //thirdRoomId已经存在

	//app
	API_APP_ACCOUNT_OVERDUE_ERROR = 3001 //APP账号已过期
	API_STOP_APP_SERVICE_ERROR    = 3003 //暂停服务

	//attachment
	API_REPECT_ATTACHMENT_ERROR = 4001 //附件重复
)

//
var HttpHeaders = map[string]string{
	"Accept":        "application/json",
	"Content-Type":  "application/json",
	"Authorization": "", //验证token
}

//通用调用参数:

//index //索引, 默认0
//count //数量, 默认10，最大100
//通用返回值:

//accessPath //访问接口路径
//runtime //接口执行时间
//systime //服务器调用时间戳
//affectedRows //受影响的行,一般添加或修改时会返回
//status //状态码
//errorDesc //错误详细信息
