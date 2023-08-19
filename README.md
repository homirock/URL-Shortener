# URL-Shortener
Build a simple URL shortener service that will accept a URL as an argument over a REST API and return a shortened URL as a result.
Run the Application
=======================
step1- build docker image
docker build -t url-shortener .
step2- Run the docker image 
docker run -p 8084:8084 url-shortener

APIs
=========
1.shorten url https://example.com
Request-
curl -X POST -d '{"url": "https://example.com"}' http://localhost:8084/shorten
response-
{"short_url":"JBkBVoA"}
2.get Actual url of shortenurl
Request-curl 'http://localhost:8084/r/JBkBVoA'
Response-
<a href="https://example.com/">See Other</a>.
3.get 3 top domain name
Request-
curl 'http://localhost:8084/metrics'
response-
{"top_domains":["example.com","play.golang.com"]}

