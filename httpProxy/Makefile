tls:
	openssl genrsa -out server.key 2048
	openssl req -new -x509 -key server.key -out server.crt -days 365

clean:
	rm *.crt *.key

qs:
	curl -Lv --proxy https://127.0.0.1:8888 --proxy-cacert server.crt https://www.baidu.com
q:
	curl -Lv --proxy 127.0.0.1:8888 https://www.baidu.com
