package main

import (
	"log"
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/emicklei/go-restful"
	//"github.com/emicklei/go-restful/swagger"
)

type MasterYaml struct {
        DbType string
	MasterHost string
        Params interface{}
}

type Mysql struct {
        LogPos string
        LogFile string
	MasterUser string
 	MasterPass string
	MysqlUser string
	MysqlPass string
	
}

// PUT http://localhost:8080/daemon/CreateProdFollow -d '{"DbType": "mysql", "MasterHost":"209.205.208.122"}'
//

func CreateProdFollow(request *restful.Request, response *restful.Response) {

	var param json.RawMessage
	JsonParam := MasterYaml{
                Params: &param,
        }
	fmt.Println("Enter CreateProdFollow function..\n ")
        
	if (request.HeaderParameter("Content-Type") == "application/x-www-form-urlencoded"){

		request.Request.Header.Set("Content-Type", "application/json")

	}
	fmt.Printf("\n\n",request.HeaderParameter("Content-Type"))

	err := request.ReadEntity(&JsonParam)

	if err != nil{
                fmt.Printf("\nfailed in readng json\n")
                response.AddHeader("Content-Type", "text/plain")
                response.WriteErrorString(http.StatusInternalServerError, err.Error())
                return
        }

	
	fmt.Println("\n\nDbType is :", JsonParam.DbType)

        /* Create service for master
	   TBD
	*/
	
        /*Fetch DB type and decode params*/

 	switch JsonParam.DbType {
        case "mysql":
		var res Mysql = DecodeJsonMysql(param)
		fmt.Println("Switch ::" , res)
		
		/* Populate mysql template*/
		err := PopulateMysqlYaml(res, JsonParam.MasterHost)	
		
		if err != nil{
			fmt.Printf("\nfailed in reading SQL yaml\n")
                	response.AddHeader("Content-Type", "text/plain")
                	response.WriteErrorString(http.StatusInternalServerError, err.Error())
                	return
		}
		/*Connect to Controller for deploying app using dt*/
		
	
        default:
		log.Printf("unknown DB type: %q", JsonParam.DbType)
		response.AddHeader("Content-Type", "text/plain")
                response.WriteErrorString(http.StatusInternalServerError, "Unknown DB type")
        	return
	}
	response.WriteHeader(200)

}

