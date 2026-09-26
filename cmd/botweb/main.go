package main
import("flag";"log";"net/http";"os";"github.com/Ploos-AS/BotWeb/internal/server")
func main(){listen:=flag.String("listen","127.0.0.1:8080","HTTP listen address");registry:=flag.String("registry","bots.json","bot registry JSON");flag.Parse();auth:=server.AuthConfig{ViewerToken:os.Getenv("BOTWEB_VIEWER_TOKEN"),OperatorToken:os.Getenv("BOTWEB_OPERATOR_TOKEN")};s,err:=server.NewWithAuth(*registry,auth);if err!=nil{log.Fatal(err)};if !authEnabled(auth)&&!isLoopbackListen(*listen){log.Fatal("refusing non-loopback listen without BOTWEB_VIEWER_TOKEN or BOTWEB_OPERATOR_TOKEN")};log.Printf("BotWeb listening on %s",*listen);log.Fatal(http.ListenAndServe(*listen,s.Handler()))}
func authEnabled(a server.AuthConfig)bool{return a.ViewerToken!=""||a.OperatorToken!=""}
func isLoopbackListen(s string)bool{return len(s)>=9&&(s[:9]=="127.0.0.1"||s=="localhost:8080")}
