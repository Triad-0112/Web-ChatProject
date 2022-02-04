package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"webchat/tracer"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar,
}

func main() {

	//Authkey access to call login
	gomniauth.SetSecurityKey("983njfnv90n90n490ngf9gjg59njg002-2me,,-d;--0j3-0jrn-2njf9uh93#^T$b8cdb94bnvf0n490nfnB*#*b3nf8549nv9d9239f9mc9*G#*f94fn94ht5939dj9c0kld-0jg0tj03mdp0m9Enf49fn39ur9jgfmgn90gn940hnt0490ty49-5t940kf0fk0ck904n9gn405n9y50gj0j95y805j0ghmkb0imdo0ld-30d3-03,mog50j50kg-,40j0mf0mgo94h940t-hky-5kt-lyh-u-ko4-k-f,k=-3efk-04j-0yk=5kygko-g,koregmij40-jyh34-m-fw-fm3n0tyh40tj")
	gomniauth.WithProviders(
		google.New("60064740283-rlqg8cd2r5o3vgbjhhcph768e4js2g8o.apps.googleusercontent.com", "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", "http://localhost:8080/auth/callback/google"),
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
	)

	//using DefaultServeMux
	//Which inside DefaultServerMux using ServeHTTP to stream the data on Handle
	//and write a response after read the the pattern which was "/" in my case
	//and request it
	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			//ResponseWriter have a method to write
			w.Write([]byte(`
			<html>
				<head>
					<title>Chat</title>
				</head>
				<body>
					Let's chat!
				</body>
			</html>
			`))
		})
	*/

	// root access, whenever try to access root addr/ it will redirect to addr/chat
	r := newRoom(UseFileSystemAvatar)
	r.tracer = tracer.New(os.Stdout)

	//Not used since accessing /chat will able to chat sinc eno redirect for chat section
	//http.Handle("/", &templateHandler{filename: "chat.html"})

	http.Handle("/avatars/", http.StripPrefix("/avatars/", http.FileServer(http.Dir("./avatars/"))))

	http.Handle("/", MustAuth(&templateHandler{filename: "index.html"}))

	http.Handle("/index", MustAuth(&templateHandler{filename: "index.html"}))

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./templates/css/"))))

	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./templates/js/"))))

	//redirect chat so it need to be an authorization user / client
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))

	//redirect to login
	http.Handle("/login", &templateHandler{filename: "login.html"})

	//redirect logout handler along with delete cookies when request to log out invoked
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	//redirect loginHandler for login function
	http.HandleFunc("/auth/", loginHandler)

	//redirect to room
	http.Handle("/room", r)

	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)

	go r.run()
	//start the web server
}
