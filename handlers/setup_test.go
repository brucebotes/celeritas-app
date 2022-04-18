package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/brucebotes/celeritas"
	"github.com/brucebotes/celeritas/filesystems/miniofilesystem"
	"github.com/brucebotes/celeritas/filesystems/s3filesystem"
	"github.com/brucebotes/celeritas/filesystems/sftpfilesystem"
	"github.com/brucebotes/celeritas/filesystems/webdavfilesystem"
	"github.com/brucebotes/celeritas/mailer"
	"github.com/brucebotes/celeritas/render"
	"github.com/go-chi/chi/v5"
)

var cel celeritas.Celeritas
var testSession *scs.SessionManager
var testHandlers Handlers

func TestMain(m *testing.M) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.Secure = false

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader("../views"),
		jet.InDevelopmentMode(),
	)

	myRenderer := render.Render{
		Renderer: "jet",
		RootPath: "../",
		Port:     "4000",
		JetViews: views,
		Session:  testSession,
	}

	cel = celeritas.Celeritas{
		AppName:       "myapp",
		Debug:         true,
		Version:       "1.0.0",
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		RootPath:      "../",
		Routes:        nil,
		Render:        &myRenderer,
		Session:       testSession,
		DB:            celeritas.Database{},
		JetViews:      views,
		EncryptionKey: cel.RandomString(32),
		Cache:         nil,
		Scheduler:     nil,
		Mail:          mailer.Mail{},
		Server:        celeritas.Server{},
		FileSystems:   map[string]interface{}{},
		S3:            s3filesystem.S3{},
		SFTP:          sftpfilesystem.SFTP{},
		WebDAV:        webdavfilesystem.WebDAV{},
		Minio:         miniofilesystem.Minio{},
	}

	testHandlers.App = &cel

	os.Exit(m.Run()) // this will run m.Run() - which runs the tests in this directory - after these vars have been set
}

// This function replaces the routes in celeritas with our test version
func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cel.SessionLoad) // load middleware
	mux.Get("/", testHandlers.Home)

	fileServer := http.FileServer(http.Dir("./../public"))
	mux.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return mux
}

// This function creates a context for Sessions for our local tests
func getCtx(req *http.Request) context.Context {
	ctx, err := testSession.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
