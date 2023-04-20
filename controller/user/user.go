package userController

//controller/user/user.go

import (
	"fmt"
	"my_item_cf_go/context"
	myjwt "my_item_cf_go/middleware"
	"my_item_cf_go/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	jwtgo "github.com/dgrijalva/jwt-go"
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

// 初始化
func init() {
	user := &User{UserId: 0,
		UserCreateReq: UserCreateReq{
			UserName: "tom"}}
	fmt.Printf("user:%+v\n", user)
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
			// generateToken(context, user, "admin", 30)

			res := generateToken(context, user, "admin", 30)
			fmt.Print(res)
			return res
		} else {
			data := map[string]interface{}{
				"status": -1,
				"msg":    "验证失败, 用户不存在或者密码不正确",
			}
			return response.Resp().Json(data)
		}
	} else {
		// context.JSON(http.StatusOK, gin.H{
		// 	"status": -1,
		// 	"msg":    "json 解析失败." + err.Error(),
		// })

		data := map[string]interface{}{
			"status": -1,
			"msg":    "json 解析失败." + err.Error(),
		}
		return response.Resp().Json(data)
	}

	data := map[string]interface{}{
		"status": 200,
		"msg":    "",
	}
	return response.Resp().Json(data)
}

/*
  此工程为了简单，直接将生成token放在controller中
  有效时间长度，单位是分钟
*/
func generateToken(c *context.Context, user User, roleId string, expiredTimeByMinute int64) *response.Response {
	j := &myjwt.JWT{
		SigningKey: []byte(myjwt.SignKey),
	}
	claims := myjwt.CustomClaims{
		UserId:   user.UserId,
		UserName: user.UserName,
		RoleId:   roleId,
		StandardClaims: jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000),                   // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + expiredTimeByMinute*60), // 过期时间 一小时
			Issuer:    "ginjwtdemo",                                      //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)

	if err != nil {
		// c.JSON(http.StatusOK, gin.H{
		// 	"status": -1,
		// 	"msg":    err.Error(),
		// 	"data":   "{}",
		// })
		// return

		res := map[string]interface{}{
			"status": -1,
			"msg":    err.Error(),
			"data":   "",
		}
		return response.Resp().Json(res)
	}

	data := LoginResp{
		Token: token,
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": 0,
	// 	"msg":    "登录成功！",
	// 	"data":   data,
	// })
	// return

	res := map[string]interface{}{
		"status": 200,
		"msg":    "登录成功！",
		"data":   data,
	}
	return response.Resp().Json(res)
}

func GetAllUsers(c *context.Context) *response.Response {
	// func (ctl *UserController) GetAllUsers(c *gin.Context) {
	claimsFromContext, err := c.Get(myjwt.Gin_Context_Key)
	if err != false {
		panic("Failed to get all users")
	}

	fmt.Print(claimsFromContext)
	fmt.Print(claimsFromContext.(*myjwt.CustomClaims))

	claims := claimsFromContext.(*myjwt.CustomClaims)
	currentUser := claims.UserName
	// klog.Infof("get all users, loginUser:%q", currentUser)
	fmt.Print("get all users, loginUser:%q", currentUser)
	var users []User
	for i := 0; i < 3; i++ {
		userName := fmt.Sprintf("tom%d", i)
		user := User{UserId: 1}
		user.UserName = userName
		users = append(users, user)
	}

	fmt.Print("users")
	fmt.Print(users)

	res := map[string]interface{}{
		"status": http.StatusOK,
		"msg":    "登录成功！",
		"data": gin.H{
			"result": users,
			"count":  len(users),
		},
	}
	return response.Resp().Json(res)

}
