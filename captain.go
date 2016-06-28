package main

import (
        "log"
        "net/http"
        "os"
        "github.com/emicklei/go-restful"
        "github.com/emicklei/go-restful/swagger"
)


func main() {
        restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))

        ws := new(restful.WebService)
        ws.Path("/daemon")
        ws.Consumes(restful.MIME_JSON)
        ws.Produces(restful.MIME_JSON)

        ws.Route(ws.PUT("/CreateProdFollow").To(CreateProdFollow).
                // docs
                Doc("Creates production follow ").
                Operation("CreateProdFollow").
                Consumes(restful.MIME_JSON, restful.MIME_XML, "application/x-www-form-urlencoded").
                Produces(restful.MIME_JSON).
                Param(ws.BodyParameter("MasterYaml", "Information from Master").DataType("main.MasterYaml")))

        restful.Add(ws)

        config := swagger.Config{
                WebServices:    restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
                WebServicesUrl: "http://localhost:8080",
                ApiPath:        "/apidocs.json",

                SwaggerPath:     "/apidocs/",
                SwaggerFilePath: "/root/workspace/src/swagger-ui/dist"}
        swagger.RegisterSwaggerService(config, restful.DefaultContainer)

        log.Printf("start listening on localhost:8080")
        server := &http.Server{Addr: ":8080", Handler: restful.DefaultContainer}
        log.Fatal(server.ListenAndServe())
}
