// Copyright © 2021 Lucidworks
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
	// "fmt"
	"fmt"
	"log"
	"os"
	"strings"

	// "github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

type ManifestValues struct {
	CloudProvider string `json:"cloud_provider"`
	CustomerName  string `json:"customer_name"`
	ClusterType   string `json:"cluster_type"`
	ClusterName   string `json:"cluster_name"`
	Env           string `json:"env"`
	Domain        string `json:"domain"`
	Namespace     string `json:"namespace"`
	Version       string `json:"version"`
	OutputDir     string `json:"output_dir"`
	ClusterSize   string `json:"clusterSize"`

}

var (
	cloud     string
	env       string
	namespace string
	domain    string
	version   string
)

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Generates k8s templates to deploy the managed fusion services.",
	Long:  `Generates k8s templates to deploy the managed fusion services.`,
	Run: func(cmd *cobra.Command, args []string) {
		manifestExecutor()
	},
}

func init() {
	rootCmd.AddCommand(manifestCmd)
	manifestCmd.Flags().StringVarP(&clusterType, "clusterType", "t", "", "The type of cluster to generate the templates for. eg: prod, non-prod, prod-<xyz>, non-prod-<xyz> (required)")
	manifestCmd.Flags().StringVarP(&env, "env", "e", "", "The environment type to generate the manifest for. eg: dev, stg, prd etc (required)")
	manifestCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "The namespace to deploy the manifest into (required)")
	manifestCmd.Flags().StringVarP(&cloud, "cloud", "p", "", "The cloud to generate the template for. Options: gke, eks, or aks (required)")
	manifestCmd.Flags().StringVarP(&domain, "domain", "d", "", "The domain to use when generating the values for any ingress host. (default '<cloudPrefix>.lucidworks.cloud')")
	manifestCmd.Flags().StringVarP(&version, "version", "v", "", "The version of fusion to deploy. (required)")
	manifestCmd.Flags().StringVarP(&clusterSize, "clusterSize", "s", "large", "The siZe of the cluster large or small for prod(required)")
	manifestCmd.MarkFlagRequired("clusterType")
	manifestCmd.MarkFlagRequired("env")
	manifestCmd.MarkFlagRequired("namespace")
	manifestCmd.MarkFlagRequired("cloud")
	manifestCmd.MarkFlagRequired("version")
}

func manifestExecutor() {
	manifestValues := new(ManifestValues)
	manifestValues.CustomerName = customerName
	manifestValues.ClusterType = clusterType
	manifestValues.ClusterName = fmt.Sprintf("lw-%s-%s", strings.ToLower(customerName), strings.ToLower(clusterType))
	manifestValues.Env = genereateNormalizeddEnv(env)
	manifestValues.Namespace = namespace
	manifestValues.CloudProvider = getCloudProvider(cloud)
	manifestValues.Version = version
	manifestValues.ClusterSize = clusterSize

	if outputDir == "" {
		outputDir = strings.ToLower(clusterType)
	}
	manifestValues.OutputDir = outputDir

	if domain == "" {
		cloudPrefix := getCloudPrefix(getCloudProvider(cloud))
		manifestValues.Domain = fmt.Sprintf("%s.lucidworks.cloud", cloudPrefix)
	} else {
		manifestValues.Domain = domain
	}

	if isProdLike(env) {
		if clusterSize == "small" {
			log.Printf("Processing small manifests (%s)...", namespace)
			templateDirectory(manifestValues, "templates/fusion/small/prod", fmt.Sprintf("%s/%s", outputDir, namespace))
		} else if clusterSize == "large" {
			log.Printf("Processing standard manifests (%s)...", namespace)
			templateDirectory(manifestValues, "templates/fusion/standard/prod", fmt.Sprintf("%s/%s", outputDir, namespace))
		}
	} else {
		log.Printf("Processing standard manifests (%s)...", namespace)
		templateDirectory(manifestValues, "templates/fusion/standard/dev", fmt.Sprintf("%s/%s", outputDir, namespace))
	}

	log.Printf("Processing shared manifests (%s)...", namespace)
	templateDirectory(manifestValues, "templates/shared/manifests/fusion", fmt.Sprintf("%s/%s", outputDir, namespace))
	templateDirectory(manifestValues, "templates/shared/manifests/infra", fmt.Sprintf("%s/infra", outputDir))
	templateDirectory(manifestValues, "templates/shared/manifests/ingress-nginx", fmt.Sprintf("%s/ingress-nginx", outputDir))

	log.Printf("Processing TLS manifests (%s)...", namespace)
	if getCloudProvider(cloud) == "gcp" {
		templateDirectory(manifestValues, "templates/shared/manifests/tls/gcp", fmt.Sprintf("%s/infra/templates", outputDir))
		templateDirectory(manifestValues, "templates/shared/manifests/tls/gcp", fmt.Sprintf("%s/%s/templates", outputDir, namespace))
	} else if getCloudProvider(cloud) == "aws" {
		templateDirectory(manifestValues, "templates/shared/manifests/tls/aws", fmt.Sprintf("%s/infra/templates", outputDir))
		templateDirectory(manifestValues, "templates/shared/manifests/tls/aws", fmt.Sprintf("%s/%s/templates", outputDir, namespace))
	}else if getCloudProvider(cloud) == "azure" {
		templateDirectory(manifestValues, "templates/shared/manifests/tls/azure", fmt.Sprintf("%s/infra/templates", outputDir))
		templateDirectory(manifestValues, "templates/shared/manifests/tls/azure", fmt.Sprintf("%s/%s/templates", outputDir, namespace))
	}

	log.Printf("Processing argocd application manifests (%s)...", clusterType)
	templateDirectory(manifestValues, "templates/shared/applications", fmt.Sprintf("%s/applications", outputDir))

	log.Printf("Commiting files to the customer's github repo...")
	filesToCommit := getFilesInDir(outputDir)
	githubOrgName := "lucidworks-managed-fusion"
	customerRepo := fmt.Sprintf("%s-config", customerName)
	commitBranch := "init"
	commitMessage := fmt.Sprintf("inital commit for %s - %s", customerName, namespace)
	authorName := "lw-cloud-infra-service-account"
	authorEmail := "lw-cloud-infra-service-account@lucidworks.com"
	prSubject := fmt.Sprintf("inital commit for %s - %s", customerName, namespace)
	prBranch := "main"
	prDescription := fmt.Sprintf("inital commit for %s - %s", customerName, namespace)
	commitFilesToGithub(githubOrgName, customerRepo, &commitBranch, &commitMessage, filesToCommit, &authorName, &authorEmail, &prSubject, &prBranch, &prDescription)

	os.RemoveAll(outputDir)
}

func genereateNormalizeddEnv(env string) string {
	switch strings.ToLower(env) {
	case "dev", "development":
		return "dev"
	case "qa":
		return "qa"
	case "stg", "stage":
		return "stg"
	case "perf", "performance":
		return "perf"
	case "prd", "prod", "production":
		return "prd"
	default:
		log.Fatalf("Unknow env value: %s", env)
		return ""
	}
}

func isProdLike(env string) bool {
	switch strings.ToLower(env) {
	case "dev", "development":
		return false
	case "qa":
		return false
	case "stg", "stage":
		return true
	case "perf", "performance":
		return true
	case "prd", "prod", "production":
		return true
	default:
		log.Fatalf("Unknow env value: %s", env)
		return false
	}
}

func getCloudPrefix(cloud string) string {
	lowered := strings.ToLower(cloud)
	switch lowered {
	case "aws":
		return "c"
	case "eks":
		return "c"
	case "gcp":
		return "b"
	case "gke":
		return "b"
	case "aks":
		return "a"
	case "azure":
		return "a"
	default:
		log.Fatalf("Unknow cloud value", cloud)
		return ""
	}
}

func getCloudProvider(cloud string) string {
	lowered := strings.ToLower(cloud)
	switch lowered {
	case "aws":
		return "aws"
	case "eks":
		return "aws"
	case "gcp":
		return "gcp"
	case "gke":
		return "gcp"
	case "aks":
		return "azure"
	case "azure":
		return "azure"
	default:
		log.Fatalf("Unknow cloud value", cloud)
		return ""
	}
}
