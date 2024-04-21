/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init folder structure",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		// Get the current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}

		folders := []string{"api", "internal", "configs", "internal/server", "ent", "build", "schema"}
		for _, folder := range folders {
			err := os.Mkdir(folder, os.ModePerm)
			if err != nil {
				fmt.Printf("Error creating folder %s: %v\n", folder, err)
				return
			}
		}

		os.Chdir(currentDir)

		// Initialize Go module
		initCmd := exec.Command("go", "mod", "init", args[0])
		_, err = initCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error initializing Go module:", err)
			return
		}

		// Create args[0].proto file in api folder
		re := regexp.MustCompile(`([^/]+)$`)
		name := re.FindString(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		protoFileName := fmt.Sprintf("%s.proto", strings.ToLower(name))
		protoFilePath := filepath.Join(currentDir, "api", protoFileName)
		protoFile, err := os.Create(protoFilePath)
		if err != nil {
			fmt.Printf("Error creating %s.proto file: %v\n", args[0], err)
			return
		}
		defer protoFile.Close()
		// Define the content of the proto file
		protoContent := fmt.Sprintf(`
syntax = "proto3";

package %v;

option go_package =" %v";

service %v {

}`, name, fmt.Sprintf("%v/api;%v", args[0], name), name)

		// Write the content to the proto file
		_, err = protoFile.WriteString(protoContent)
		if err != nil {
			fmt.Printf("Error writing to %s.proto file: %v\n", name, err)
			return
		}

		//config file
		configFile, err := os.Create(filepath.Join(currentDir, "configs", "config.yml"))
		if err != nil {
			fmt.Printf("Error creating config file: %v", err)
			return
		}
		defer configFile.Close()

		configGoPath := filepath.Join(currentDir, "configs", "config.go")
		configGoFile, err := os.Create(configGoPath)
		if err != nil {
			fmt.Printf("Error creating config.go file: %v", err)
			return
		}
		defer configGoFile.Close()

		configContent := `
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Enviroment string

	Server struct {
		Host string
		Port int
	}

	Db struct {
		Username string
		Password string
		Host     string
		Port     int
		Database string
	}
}

func LoadConfig() (*Config, error) {
	var config *Config

	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
`
		_, err = configGoFile.WriteString(configContent)
		if err != nil {
			fmt.Printf("Error writing to config.go file: %v\n", err)
			return
		}

		//Go dependency
		dependency := []string{
			"entgo.io/ent/cmd/ent",
			"google.golang.org/protobuf",
			"google.golang.org/grpc",
			"golang.org/x/text",
			"github.com/spf13/viper",
			"github.com/go-sql-driver/mysql",
		}
		for _, v := range dependency {
			getCmd := exec.Command("go", "get", v)
			_, err := getCmd.CombinedOutput()
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		//schema
		baseSchemaPath := filepath.Join(currentDir, "schema", "base.go")
		baseSchemaFile, err := os.Create(baseSchemaPath)
		if err != nil {
			fmt.Printf("Error creating config.go file: %v", err)
			return
		}
		defer configGoFile.Close()

		baseSchemaContent := `
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type Base struct {
	mixin.Schema
}

func (Base) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id"),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}
`
		_, err = baseSchemaFile.WriteString(baseSchemaContent)
		if err != nil {
			fmt.Printf("Error writing to base.go file: %v\n", err)
			return
		}

		sampleSchemaPath := filepath.Join(currentDir, "schema", "sample.go")
		sampleSchemaFile, err := os.Create(sampleSchemaPath)
		if err != nil {
			fmt.Printf("Error creating sample.go file: %v", err)
			return
		}
		defer configGoFile.Close()

		sampleSchemaContent := `
package schema

import (
	"entgo.io/ent"
)

type Sample struct {
	ent.Schema
}

func (Sample) Mixin() []ent.Mixin {
	return []ent.Mixin{
		Base{},
	}
}
`
		_, err = sampleSchemaFile.WriteString(sampleSchemaContent)
		if err != nil {
			fmt.Printf("Error writing to base.go file: %v\n", err)
			return
		}

		generatePath := filepath.Join(currentDir, "ent", "generate.go")
		generateSchemaFile, err := os.Create(generatePath)
		if err != nil {
			fmt.Printf("Error creating config.go file: %v", err)
			return
		}
		defer configGoFile.Close()

		generateContent := `
package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --target ../ent --feature sql/lock,sql/modifier,sql/upsert,sql/execquery ../schema

`
		_, err = generateSchemaFile.WriteString(generateContent)
		if err != nil {
			fmt.Printf("Error writing to generate.go file: %v\n", err)
			return
		}

		
		// UpdateGoMod()
		entCmd := exec.Command("go", "run", "-mod=mod", "entgo.io/ent/cmd/ent", "generate", "--target", "./ent", "--feature", "sql/lock,sql/modifier,sql/upsert,sql/execquery", "./schema")
		_, err = entCmd.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			return
		}
		UpdateGoMod()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func UpdateGoMod() {
	tidy := exec.Command("go", "mod", "tidy")
	_, err := tidy.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return
	}

	vendorCmd := exec.Command("go", "mod", "vendor")
	_, err = vendorCmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
