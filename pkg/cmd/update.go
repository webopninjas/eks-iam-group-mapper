// Copyright © 2019 WEBOP NINJAS, LLC <shawn@webop.ninja>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/endpoints"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"github.com/kataras/golog"
	"os"
	"strings"
)

var groupMaps []string
var finalUserMap map[string]MapUserConfig
var kubeconfig string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update your clusters (local) aws-auth configmap",
	Long:  `Update your clusters (local) aws-auth configmap`,
	Run: func(cmd *cobra.Command, args []string) {
		mapGroupstoNiceList()

	},
}

type MapUserConfig struct {
	UserArn  string   `yaml:"userarn"`
	Username string   `yaml:"username"`
	Groups   []string `yaml:"groups"`
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringArrayVar(&groupMaps, "map", []string{}, "Mapping of IAM Group to Kubernetes Group (ex: 'AdminGroup:KubernetesAdminGroup') ")
	updateCmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Optional Kubeconfig file for testing")
	updateCmd.Flags().BoolP("verbose", "v", false, "")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func cleanupList(groupMap []string) map[string][]string {
	finalMap := map[string][]string{}
	for _, tempMap := range groupMap {
		sliced := strings.Split(tempMap, ":")
		finalMap[sliced[0]] = append(finalMap[sliced[0]], sliced[1])
	}
	return finalMap

}
func getIamClient() *iam.IAM {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config " + err.Error())
	}
	cfg.Region = endpoints.UsEast2RegionID
	return iam.New(cfg)
}

func mapGroupstoNiceList() { // map[string]MapUserConfig {
	svc := getIamClient()
	groups := cleanupList(groupMaps)

	for group, kgroups := range groups {
		fmt.Println(group)
		fmt.Println(kgroups)
		data := getAwsIamGroup(svc, &group)
		fmt.Println(data)
		for _, user := range data.Users {
			fmt.Println(user)
		}
	}
	//	data := getAwsIamGroup(svc, &group)
	//	fmt.Println(kgroups)
	//	for _, user := range data.Users {
	//		fmt.Println(user)
	//	}
	//}

}

func getAwsIamGroup(iamClient *iam.IAM, groupName *string) *iam.GetGroupOutput {
	groupInput := &iam.GetGroupInput{GroupName: groupName}
	groupReq := iamClient.GetGroupRequest(groupInput)
	group, err := groupReq.Send()
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				golog.Error(iam.ErrCodeNoSuchEntityException, aerr.Error())
			case iam.ErrCodeServiceFailureException:
				golog.Error(iam.ErrCodeServiceFailureException, aerr.Error())
			default:
				golog.Error(aerr.Error())
			}
		}
	}
	return group
}

func getKubernetesClient() *kubernetes.Clientset {
	var kconfig *string
	kconfig = getKubeconfig()
	if len(*kconfig) > 0 {
		config, err := clientcmd.BuildConfigFromFlags("", *kconfig)
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		return clientset
	} else {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientset, err := kubernetes.NewForConfig(config)
		return clientset
	}
}

func setUsersToConfigMap(k8sClient *kubernetes.Clientset, groupMap *map[string][]string) {
	cf, err := k8sClient.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	//var newConfig []MapUserConfig
	//Build Yaml
	cf.Data["mapUsers"] = string("")

	newCF, err := k8sClient.CoreV1().ConfigMaps("kube-system").Update(cf)
	if err != nil {
		golog.Error(err)
	} else {
		golog.Info("Successfully updated user roles")
		golog.Info(newCF)
	}
}

func getKubeconfig() *string {
	if len(kubeconfig) > 0 {
		return &kubeconfig
	}
	h := os.Getenv("HOME")
	if h == "" {
		h = os.Getenv("USERPROFILE")
	}
	return &h
}
