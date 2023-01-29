package main

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-faker/faker/v4"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/weplanx/openapi/client"
	"time"
)

func invoke(ctx context.Context, event map[string]interface{}) (_ string, err error) {
	data := [][]interface{}{
		{"Name", "CCType", "CCNumber", "Century", "Currency", "Date", "Email", "URL"},
	}
	for n := 0; n < 1000000; n++ {
		data = append(data, []interface{}{
			faker.Name(), faker.CCType(), faker.CCNumber(), faker.Century(), faker.Currency(), faker.Date(), faker.Email(), faker.URL(),
		})
	}
	start := time.Now()
	if err = openapi.Excel(ctx, "xxx", map[string][][]interface{}{
		"Sheet1": data,
	}); err != nil {
		return
	}
	return fmt.Sprintf("Cost: %d!", time.Since(start).Seconds()), nil
}

var e Env

type Env struct {
	Url        string `env:"URL"`
	Key        string `env:"KEY"`
	Secret     string `env:"SECRET"`
	client.Cos `envPrefix:"COS_"`
}

var openapi *client.Client

func main() {
	err := env.Parse(&e)
	if err != nil {
		panic(err)
	}
	if openapi, err = client.New(e.Url,
		client.SetApiGateway(e.Key, e.Secret),
		client.SetCos(e.Cos.Url, e.Cos.SecretID, e.Cos.SecretKey),
	); err != nil {
		panic(err)
	}
	// Make the handler available for Remote Procedure Call by Cloud Function
	cloudfunction.Start(invoke)
}
