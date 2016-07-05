package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	//"bytes"
	"encoding/json"
	"github.com/hoisie/mustache"
	"io/ioutil"
)

func DecodeJsonMysql(param json.RawMessage) Mysql {
	var s Mysql
	if err := json.Unmarshal(param, &s); err != nil {
		log.Fatal(err)
	}
	fmt.Println("DecodeJsonMysql:::", s)
	return s
}

func ReplaceString(fileName string, res Mysql, MasterHost string) error {
	fmt.Println("\nres", res)
	context := map[string]string{
		"HOST_IP":     MasterHost,
		"LOGFILE":     res.LogFile,
		"LOGPOS":      res.LogPos,
		"REPL_USER":   res.MasterUser,
		"REPL_PASSWD": res.MasterPass,
		"RUSER":       res.MysqlUser,
		"RPASSWORD":   res.MysqlPass}

	output := mustache.RenderFile(fileName, context)
	fmt.Println("testfile output", output)
	if err := ioutil.WriteFile(fileName, []byte(output), 0666); err != nil {
		fmt.Println(err)
	}
	return nil
}

func PopulateMysqlYaml(res Mysql, MasterHost string) (error, string) {
	srcTemplate := "./templates/mysql_production_follow.yaml"
	dstYaml := "./templates/mysql_production_follow"
	fmt.Println("PopulateMysqlYaml:::", res)
	/* Check if sql yaml template exists*/
	if _, err := os.Stat(srcTemplate); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("\nfile does not exist\n")
			return err, ""
		} else {
			fmt.Println("\n", err.Error())
			return err, ""
		}
	} else {
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
				if err != nil {
					fmt.Printf("\nCopy command failed\n")
					return err, ""
				}
				/*Write values to this file*/
				err = ReplaceString(dstYaml, res, MasterHost)
				if err != nil {
					fmt.Printf("\nReplacing RES failed\n")
					return err, ""
				}

			} else {
				return err, ""
			}
		} else {
			fmt.Printf("\nFile already exists\n")
		}
	}

	return nil, dstYaml
}
