package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris"
	"publish_client_user/datamodels"
	"time"
)

//注册jwt中间件
func GetJWT() *jwtmiddleware.Middleware {
	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		//这个方法将验证jwt的token
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			fmt.Println("----------token:",token.Raw)

			//自己加密的秘钥或者说盐值
			return []byte("My spacej"), nil
		},
		//加密的方式
		SigningMethod: jwt.SigningMethodHS256,
		//验证未通过错误处理方式
		//ErrorHandler: func(context.Context, string)
		ErrorHandler: func(ctx iris.Context, s string) {
			result := datamodels.Result{Status:false,Code:"403", Msg:"认证失败，请重新登录认证"}
			ctx.JSON(result)
		},
	})
	return jwtHandler
}
//生成token
func GenerateToken(user *datamodels.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nick_name":  user.LoginName, //用户信息
		"session":  user.Session, //session
		"id":         user.ID,   //用户信息
		"iss":        "Iris",                                                   //签发者
		"iat":        time.Now().Unix(),                                        //签发时间
		"jti":        "9527",                                                   //jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击。
		"exp":        time.Now().Add(10 * time.Hour * time.Duration(1)).Unix(), //过期时间
	})
	tokenString, _ := token.SignedString([]byte("My spacej"))
	fmt.Println("签发时间：",time.Now().Unix())
	fmt.Println("到期时间：", time.Now().Add(10 * time.Hour * time.Duration(1)).Unix())
	return tokenString
}

func ParseToken(tokenString string, key string) (interface{}, bool){
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return "", false
	}
}
func CheckSession(ctx iris.Context) {
	fmt.Println("请求接口:",ctx.Request().URL)

	token := ctx.GetHeader("Authorization")
	if token != "" && len(token)>7 {
		token = token[7 : len(token)]
	}
	fmt.Println(token)
	var sessionID = ""
	if token != "" && token != "undefined" {
		v,_ := ParseToken(token, "My spacej")
		if v != ""{
			sessionID = v.(jwt.MapClaims)["session"].(string)
		}
	}

	fmt.Println("---------sessionID:----------------",sessionID)
	var onlineSessionIDList = SMgr.GetSessionIDList()
	for _, onlineSessionID := range onlineSessionIDList {
		fmt.Println("在线用户SessionID:",onlineSessionID)
		if onlineSessionID == sessionID {
			ctx.Next()
			return
		}
	}
	fmt.Println("接口",ctx.Request().URL,"-------------------被拦截了")

	result := datamodels.Result{Status:false,Code:"512", Msg:"该账号已在其它设备登陆!"}
	ctx.JSON(result)
	//ctx.Next()
}