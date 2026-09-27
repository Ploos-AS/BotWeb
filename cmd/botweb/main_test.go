package main
import("net/http";"testing";"time")
func TestHTTPServerPolicy(t *testing.T){s:=newHTTPServer("127.0.0.1:0",http.NewServeMux());if s.ReadHeaderTimeout!=5*time.Second{t.Fatalf("ReadHeaderTimeout=%s",s.ReadHeaderTimeout)};if s.ReadTimeout!=15*time.Second{t.Fatalf("ReadTimeout=%s",s.ReadTimeout)};if s.WriteTimeout!=30*time.Second{t.Fatalf("WriteTimeout=%s",s.WriteTimeout)};if s.IdleTimeout!=60*time.Second{t.Fatalf("IdleTimeout=%s",s.IdleTimeout)}}
