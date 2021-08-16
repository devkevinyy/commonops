package exception

import "errors"

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/11/21 4:52 PM
 * @Desc:
 */

var (
	ArgsException = errors.New("输入参数异常。")
	TokenInvalidException = errors.New("token失效，请重新登录")
	UserInvalidException = errors.New("用户不存在或被禁用，请联系管理员")
	DmsUserAuthNeedValidDatabaseException = errors.New("用户权限需要指定到库！")
)