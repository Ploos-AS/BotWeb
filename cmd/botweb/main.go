package main
import("flag";"log";"net/http";"github.com/Ploos-AS/BotWeb/internal/server")
func main(){listen:=flag.String("listen","127.0.0.1:8080","HTTP listen address");registry:=flag.String("registry","bots.json","bot registry JSON");flag.Parse();s,err:=server.New(*registry);if err!=nil{log.Fatal(err)};log.Printf("BotWeb listening on %s",*listen);log.Fatal(http.ListenAndServe(*listen,s.Handler()))}
