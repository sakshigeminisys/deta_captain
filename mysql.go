package main

import (
        "log"
        "fmt"
	"os"
	"os/exec"
	"bytes"
	"io/ioutil"
        "encoding/json"
)

func DecodeJsonMysql(param json.RawMessage) Mysql {
	var s Mysql
	if err := json.Unmarshal(param, &s); err != nil {
		log.Fatal(err)
	}
	fmt.Println("DecodeJsonMysql:::", s)
	return s
}

func ReplaceString(fileName string, sText string, rText string) error {
	input, err := ioutil.ReadFile(fileName)
        if err != nil {
        	fmt.Println(err)
		return err
         }
	output := bytes.Replace(input, []byte(sText), []byte(rText), -1)
	if err = ioutil.WriteFile(fileName, output, 0666); err != nil {
                 fmt.Println(err)
                 return err
        }
	return nil
}


func PopulateMysqlYaml(res Mysql, MasterHost string) error {
	srcTemplate := "./templates/mysql_production_follow.yaml"
	dstYaml := "./templates/mysql_production_follow"
	fmt.Println("PopulateMysqlYaml:::", res)
	/* Check if sql yaml template exists*/ 
	if _, err := os.Stat(srcTemplate); err != nil {
    		if os.IsNotExist(err) {
        		fmt.Println("\nfile does not exist\n")
			return err
    		} else {
        	 	fmt.Println("\n", err.Error())
			return err
    		}
	}else{
		/* Copy the template to a new file and write values*/
		dstYaml += "_"
        	dstYaml += MasterHost
		dstYaml += ".yaml"
		fmt.Println("\nMasterHost:", MasterHost)
		fmt.Println("\ndstYaml:", dstYaml)
		if _, err := os.Stat(dstYaml); err != nil {
			if os.IsNotExist(err) {
  				fmt.Printf("\ncreating file\n")
				cpCmd := exec.Command("cp", "-rf", srcTemplate, dstYaml)
                        	err := cpCmd.Run()
                        	if err != nil{
                                	fmt.Printf("\nCopy command failed\n")
                                	return err
                        	}
				/*Write values to this file*/
				err = ReplaceString(dstYaml, "HOST_IP", MasterHost)
				if err != nil{
                                	fmt.Printf("\nReplacing HOST_IP failed\n")
                                	return err
                        	}
				err = ReplaceString(dstYaml, "LOGFILE", res.LogFile)
                        	if err != nil{
                                	fmt.Printf("\nReplacing LogFile failed\n")
                                	return err
                        	}
				err = ReplaceString(dstYaml, "LOGPOS", res.LogPos)
                        	if err != nil{
                                	fmt.Printf("\nReplacing LOGPOS failed\n")
                                	return err
                        	}
				err = ReplaceString(dstYaml, "REPL_USER", res.MasterUser)
                        	if err != nil{
                                	fmt.Printf("\nReplacing MasterUser failed\n")
                                	return err
                        	}
				err = ReplaceString(dstYaml, "REPL_PASSWD", res.MasterPass)
                        	if err != nil{
                                	fmt.Printf("\nReplacing MasterPass failed\n")
                                	return err
                        	}
				err = ReplaceString(dstYaml, "RUSER", res.MysqlUser)
                        	if err != nil{
                                	fmt.Printf("\nReplacing MysqlUser failed\n")
                                	return err
                        	}
				err = ReplaceString(dstYaml, "RPASSWORD", res.MysqlPass)
                        	if err != nil{
                                	fmt.Printf("\nReplacing MysqlPass failed\n")
                                	return err
                        	}
			}else{
				return err
			}
		}else{
			fmt.Printf("\nFile already exists\n")
		}
	}	

	return nil
}
