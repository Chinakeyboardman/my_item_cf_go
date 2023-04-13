package user

//controller/user/user.go

import (
	"fmt"
	"my_item_cf_go/component/lock"
	"my_item_cf_go/context"
	myjwt "my_item_cf_go/middleware"
	"my_item_cf_go/response"
	"net/http"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

/*
  示例而已，因此字段只有几个
*/
type UserCreateReq struct {
	UserName    string
	Email       string
	Phone       string
	LockedState bool
	//etc ..
}

type LoginReq struct {
	UserName string `json:"userName"`
	//实际当中不会以明文传输密码，本工程是示例工程，为简单起见使用明文
	Passwd string `json:"passwd"`
}

// LoginResult 登录结果结构
type LoginResp struct {
	Token string `json:"token"`
}

/*
  示例而已，因此字段只有几个
*/
type User struct {
	UserId int64
	UserCreateReq
	//etc ..
}

func init() {
	user := &User{UserId: 0,
		UserCreateReq: UserCreateReq{
			UserName: "tom"}}
	fmt.Printf("user:%+v\n", user)
}

func TestItemCf(context *context.Context) *response.Response {

	l := lock.NewLock("test", 1*time.Second)

	defer l.Release()

	if l.Block(4 * time.Second) {

		data := map[string]interface{}{
			"msg":    "拿锁成功",
			"status": 200,
		}

		return response.Resp().Json(data)
	}

	return response.Resp().String("拿锁失败")
}

func Login(context *context.Context) *response.Response {
	// klog.Infof("login to get a token")
	var loginReq LoginReq
	if err := context.ShouldBindJSON(&loginReq); err == nil {
		//实际当中需要检查用户名和密码的正确性，这里为了简单起见，hardcode，只要和用户是tom，密码是123456就允许通过
		// check whether username exists and passwd is matched
		if loginReq.UserName == "tom" && loginReq.Passwd == "123456" {
			user := User{}
			user.UserName = loginReq.UserName
			user.UserId = 0
			generateToken(context, user, "admin", 30)
		} else {
			data := map[string]interface{}{
				"status": -1,
				"msg":    "验证失败, 用户不存在或者密码不正确",
			}
			return response.Resp().Json(data)
		}
	} else {
		data := map[string]interface{}{
			"status": -1,
			"msg":    "json 解析失败." + err.Error(),
		}
		return response.Resp().Json(data)
	}
}

/*
  此工程为了简单，直接将生成token放在controller中
  有效时间长度，单位是分钟
*/
func generateToken(c *context.Context, user User, roleId string, expiredTimeByMinute int64) {
	j := &myjwt.JWT{
		[]byte(myjwt.SignKey),
	}
	claims := myjwt.CustomClaims{
		user.UserId,
		user.UserName,
		roleId,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),                   // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + expiredTimeByMinute*60), // 过期时间 一小时
			Issuer:    "ginjwtdemo",                                      //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	data := LoginResp{
		Token: token,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "登录成功！",
		"data":   data,
	})
	return
}
