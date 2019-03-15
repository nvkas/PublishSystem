package tool

/*
root:静态文件地址	例如:/home/vue/dist/test
port:监听端口 例如:9527
 */
func GetNginxTemplate(root string,port string,locations []string) (template string) {
	template =  "server {"+
					"listen "+port+";"+
					"server_name  localhost;"+
					"root "+root+";"

	if len(locations)>0{
		for _,location := range locations{
			template += "location "+location+" {"+
				"try_files $uri $uri/ /index.html last;"+
				"index  index.html index.htm;"+
				"}"
		}
	}else {
		template += "location / {"+
			"try_files $uri $uri/ /index.html last;"+
			"index  index.html index.htm;"+
			"}"
	}

	template += "error_page   500 502 503 504  /50x.html;"+
					"location = /50x.html {"+
						"root   html;"+
					"}"+
				"}"

	return
}
