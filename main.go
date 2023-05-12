package main

import (
	"context"
	"fmt"
	"github.com/KubeOperator/kobe/api"
	utilansible "github.com/rubinus/go-check-k8s/util/ansible"
	"github.com/rubinus/go-check-k8s/util/check"
	"github.com/rubinus/go-check-k8s/util/kobe"
	"github.com/rubinus/go-check-k8s/util/phases"
	"log"
	"os"
	"sigs.k8s.io/yaml"
)

const inventoryName = "demo.yml"
const IP = "10.10.12.121"

func main() {
	writer, err := utilansible.CreateAnsibleLogWriterWithId("cluster.Name", "123")
	if err != nil {
		fmt.Println(err)
		return
	}
	//inventory
	var hosts []*api.Host
	apiHost := &api.Host{
		Name:       IP,
		Ip:         IP,
		Port:       int32(22),
		User:       "root",
		//Password:   "1qazXCDE#@",
		//121的privateKey
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAwxE7l+ggf/Lg297+BuJHiZMe0ljptAni0Q1Q0vhQ45/bXc7q\nP5g9icQNwVDead7MTEwL0E6iYOF3OGkK2MSVLH/JRV1H7iPaW5vf8EdsJ80z4R48\ngm7pER/WWJS7Gk4pkuqPJofyYooUnZM/RJgNLNFRefpBElbJr408kFn2SoOaXAYk\nBGmnBWK+hq1+uHtI0onP7WGllWz7P2juOXjB6xsSAR2Q6QPQXBGFgR0IyM4ljF05\neb5aXGHBic9QmWZyIlLytbQBcV3AZegbbsLXv7roBaHpolLCzMuAAo6m3ar36GsZ\nJ9cxGsHGH0llWzLXGiJ5k2alK5loJdMjV+tPJwIDAQABAoIBAFAsbQnqZjEwair0\nZAQATNbGmQxrbuKIjajOiEGtvdFQiqDrwmuQ7voIkn659jAdWmqhdtmO+D5JbO7K\nfaLKaWV4wAi6Zi4CnmS4lDn1oQZa2M/V1ZjmbPSU2UxfFOei6INx1JRJm93UUtTR\nCFfxBrk87vfrW4NmGE8HBbVuxEOrS7HfPl1hyf0QchQ/9r7HKzllkZt/OsYhtMIb\ne3HBmlRz3v/rIQNm5rGuTvcpqsNJBbeX1mFEVvRPtHgtSIg+zIdFE7ha+QQyymuw\nBHAcSRxepmhd5jxIHN3xza6AlVzUQcCPgWkyS3rJpF/I1FZs8BRi684eLR3maj9I\n602ZlwECgYEA5ETeIxOOrvZLXYOuat0B2llxiGsVO8Fx9UqKr099nmqh2ZNudCkS\niyNwc2BV1hemIadleXyZCRZ7J/A8/gVBGGgqiO0Gj40R3VKlrddBKmcNdROhkLPH\nIeX32sjGv2Ct05YQvWDgeQZ1g17V0jvIQ1F8Ad4xNnkEr/fNCaLwzzECgYEA2sPM\nFCIL8IrsqDWJzsGRNeZYa/6G5gN1QWCacstKyA3CGjUsBi6lXtam/sOOiSthQgG7\nK/Gjmso1lwQTaxWih6s9hY9cnBjzUqD1yVx/sE4v+D0n+S/L8NRwiokZTEFdQolj\nYzi/oJz0c6B38gSQm4bHKNpLJZeOVNPlTlXQ3dcCgYB/35h1G9qdZrm3bDIECUSl\ndd+k5R/i9q4JFDX2mVgsq115jh6dEfkiWrr+1yOeqGbXiXfOA5+TOLXLHMh+IKFj\ns87IH8fCGOu+CTNo3CHUSCCAynuCnUNbWQFs3XaA9P7LfdBo1mFJSvX/ntu3Rugb\n1gTa4wa8ljSrAu0ojc/KsQKBgCTbzGyv99cFcS4+JwPg9ThhoRBBCDWE66KiRiOF\nQQpH1yZXQx2fillaTTSrej5+Qpq+c+zJf8k6vKC/HQ5zzLiTD4CLUQ0z3vtTB1Zv\n8UuhQM/QbgW8Gd5vzK5qvwpsEOx+/XHgQ9kp2L4KkWsDfeHWaYPmk7a3vFFqij4S\nk2htAoGBAMn6tQNwjSJpJXdD9l8r7u/CSjYw0DjFKt0NrlTI6Xg1i1VpeGBQcx4I\nQBrFhLK040YE/Gq2CqUKOCt4dGEgMunE9rjS9XxFQh13EZNRTBU0T6FLQdc2Vx1t\nZ1uxK9PXWPEEvIJ0QdSonikHFAkD3HccodibjGNNiggM2R3co8//\n-----END RSA PRIVATE KEY-----",

	}
	hosts = append(hosts, apiHost)
	group := &api.Group{
		Name:  "master",
		Hosts: []string{IP},
	}
	var groups []*api.Group
	groups = append(groups, group)
	inventory := &api.Inventory{
		Hosts:  hosts,
		Groups: groups,
	}
	k := kobe.NewAnsible(&kobe.Config{
		Inventory: inventory,
	})

	//construct inventory array yaml
	var inventTasks []check.InventoryTask
	task101 := check.InventoryTask{
		Name:  "node_kernel",
		Shell: "uname -a",
	}
	task102 := check.InventoryTask{
		Name:  "node_os_version",
		Shell: "cat /etc/redhat-release",
	}
	task103 := check.InventoryTask{
		Name:  "node_cpu",
		Shell: "vmstat 1 3 | sed -n '$p' | awk '{print ($13+$14)\"%\"}'",
	}
	task104 := check.InventoryTask{
		Name:  "node_mem",
		Shell: "free -m | grep Mem | awk '{print $3/$2*100\"%\"}'",
	}
	inventTasks = append(inventTasks, task101, task102,task103,task104)
	proInvent := check.Inventory{
		Name:        "check_node",
		Hosts:       "all",
		GatherFacts: false,
		Tasks:       inventTasks,
	}
	var invents []check.Inventory
	invents = append(invents,proInvent)
	bytes, err := yaml.Marshal(invents)
	if err != nil {
		log.Println(err)
		return
	}

	//use the yml file to do delete this func
	bytes, err = os.ReadFile(inventoryName)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = k.CreateProject("check", inventoryName, bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	stopCh := make(chan struct{})
	resCh := make(chan *kobe.AllResult)
	go func() {
		//处理结果
		for {
			select {
			case val := <-resCh:
				result := val.Result
				host0s := result.Plays[0].Tasks[0].Hosts
				var hostKey string
				for idx, _ := range host0s {
					hostKey = idx
				}
				fmt.Println(val.TaskId, result.Plays[0].Duration)
				for _, task := range result.Plays[0].Tasks {
					fmt.Println(hostKey, "查询:", task.Name, "执行结果:", task.Hosts[hostKey]["stdout"])
				}

				break
			}
		}
	}()

	err = phases.RunPlaybookAndGetResult(context.TODO(), resCh, k, inventoryName, "", writer)
	if err != nil {
		fmt.Printf("error %s", err.Error())
		return
	}

	<-stopCh
}
