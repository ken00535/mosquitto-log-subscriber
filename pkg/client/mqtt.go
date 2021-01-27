package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Host struct {
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	CaPath   string `json:"caPath"`
	CertPath string `json:"certPath"`
	KeyPath  string `json:"keyPath"`
}

func newClient(host Host) (mqtt.Client, error) {
	certpool := x509.NewCertPool()
	// config.CAPath = filepath.Join(path, path.CaPath)
	// config.CertPath = filepath.Join(path, path.CertPath)
	// config.KeyPath = filepath.Join(path, path.KeyPath)
	pemCerts, err := ioutil.ReadFile(host.CaPath)
	if err != nil {
		panic(err)
	}
	ok := certpool.AppendCertsFromPEM(pemCerts)
	if !ok {
		panic("failed to parse root ca")
	}

	cert, err := tls.LoadX509KeyPair(host.CertPath, host.KeyPath)
	if err != nil {
		panic(err)
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}

	url := "tcps://" + host.IP + ":" + strconv.Itoa(host.Port)
	opts := mqtt.NewClientOptions().AddBroker(url).
		SetUsername(host.Username).SetPassword(host.Password)
	tlsConfig := tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
	opts.SetTLSConfig(&tlsConfig)
	opts.SetClientID("log-subscriber")
	opts.SetOnConnectHandler(subscribeTopic)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}

	return client, nil
}
