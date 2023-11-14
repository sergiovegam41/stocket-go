package main

import (
    "log"
    "net/http"
    "github.com/googollee/go-socket.io"
)

func main() {
	server := socketio.NewServer(nil)

    server.OnConnect("/", func(s socketio.Conn) error {
        s.SetContext("")
        log.Println("connected:", s.ID())
        return nil
    })


    server.OnEvent("/", "client:iniEmit:conection", func(s socketio.Conn, msg string) {
        // Tu lógica de manejo aquí
        log.Println("client:iniEmit:conection", msg)
    })

    server.OnEvent("/", "client:send:message", func(s socketio.Conn, msg string) {
        // Tu lógica de manejo aquí
        log.Println("client:send:message", msg)
    })

    server.OnError("/", func(s socketio.Conn, e error) {
        log.Println("error:", e)
    })

    server.OnDisconnect("/", func(s socketio.Conn, reason string) {
        log.Println("disconnected:", reason)
    })

    go server.Serve()
    defer server.Close()

    http.Handle("/socket.io/", server)
    // Define la carpeta pública que quieres exponer, digamos que se llama "public"
    fs := http.FileServer(http.Dir("./public"))
    // Servir los archivos estáticos de la carpeta 'public' en la ruta '/'
    http.Handle("/", fs)

    log.Println("Serving at localhost:8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
