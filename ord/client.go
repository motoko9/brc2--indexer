package ord

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	url string
}

func New(url string) *Client {
	c := &Client{
		url: url,
	}
	return c
}

func (c *Client) BlockHeight() (uint64, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/blockheight", c.url)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return 0, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return 0, err
	}
	var rsp uint64
	err = json.Unmarshal(rspJson, &rsp)
	if err != nil {
		return 0, err
	}
	return rsp, nil
}

func (c *Client) BlockHash() (string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/blockhash", c.url)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return "", fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return "", err
	}
	return string(rspJson), nil
}

func (c *Client) BlockTime() (uint64, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/blocktime", c.url)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return 0, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return 0, err
	}
	var rsp uint64
	err = json.Unmarshal(rspJson, &rsp)
	if err != nil {
		return 0, err
	}
	return rsp, nil
}

func (c *Client) InscriptionById(id string) (*Inscription, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/inscription/%s", c.url, id)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	var inscription Inscription
	err = json.Unmarshal(rspJson, &inscription)
	if err != nil {
		return nil, err
	}
	return &inscription, nil
}

func (c *Client) inscriptions() ([]string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/inscriptions", c.url)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	var rsp struct {
		Inscriptions []string `json:"inscriptions"`
		Prev         string   `json:"prev"`
		Next         string   `json:"next"`
		Lowest       uint64   `json:"lowest"`
		Highest      uint64   `json:"highest"`
	}
	err = json.Unmarshal(rspJson, &rsp)
	if err != nil {
		return nil, err
	}
	return rsp.Inscriptions, nil
}

func (c *Client) Inscriptions() ([]string, error) {
	return c.inscriptions()
}

func (c *Client) inscriptionsByBlock(height uint64) ([]string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/inscriptions/block/%d", c.url, height)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	var rsp struct {
		Inscriptions []string `json:"inscriptions"`
		Prev         string   `json:"prev"`
		Next         string   `json:"next"`
		Lowest       uint64   `json:"lowest"`
		Highest      uint64   `json:"highest"`
	}
	err = json.Unmarshal(rspJson, &rsp)
	if err != nil {
		return nil, err
	}
	return rsp.Inscriptions, nil
}

func (c *Client) InscriptionsByBlock(height uint64) ([]string, error) {
	return c.inscriptionsByBlock(height)
}

func (c *Client) InscriptionContent(id string) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/content/%s", c.url, id)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	return rspJson, nil
}

func (c *Client) Output(txhash string, n int) (*Output, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/output/%s:%d", c.url, txhash, n)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspJson, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	var output Output
	err = json.Unmarshal(rspJson, &output)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

/*
func (c *Client) Tx(hash string) (*Transaction, error) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/tx/%s", c.url, hash)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpRsp.Body.Close()

	if httpRsp.StatusCode != 200 {
		return nil, fmt.Errorf("http response code %d", httpRsp.StatusCode)
	}
	rspHtml, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, err
	}
	//
	//items := bytes.Split(rspHtml, []byte("\n"))
	//rspHtml = bytes.Join(items, []byte{})
	rspHtml = bytes.TrimPrefix(rspHtml, []byte(" "))
	rspHtml = bytes.ReplaceAll(rspHtml, []byte("\n        "), []byte(""))
	rspHtml = bytes.ReplaceAll(rspHtml, []byte("\n      "), []byte(""))
	rspHtml = bytes.ReplaceAll(rspHtml, []byte("\n    "), []byte(""))
	rspHtml = bytes.ReplaceAll(rspHtml, []byte("\n  "), []byte(""))
	rspHtml = bytes.ReplaceAll(rspHtml, []byte("\n"), []byte(""))
	doc, err := html.Parse(bytes.NewReader(rspHtml))
	if err != nil {
		return nil, err
	}
	transaction := c.parser(doc)
	if transaction == nil {
		return nil, fmt.Errorf("can not find transaction")
	}
	return transaction, nil
}

func (c *Client) parser(doc *html.Node) *Transaction {
	mainNode := findNode(doc, "main")
	if mainNode == nil {
		return nil
	}
	// get transaction
	transactionNode := findNode(mainNode, "Transaction")
	if transactionNode == nil {
		return nil
	}
	hash := transactionNode.NextSibling.FirstChild.Data
	// get inputs
	inputsNode := findNode(mainNode, "Input")
	if inputsNode == nil {
		return nil
	}
	inputListNode := inputsNode.Parent.NextSibling
	inputs := make([]Input, 0)
	inputNode := inputListNode.FirstChild
	for inputNode != nil {
		inputs = append(inputs, Input{Id: inputNode.LastChild.FirstChild.Data})
		inputNode = inputNode.NextSibling
	}
	// gget outputs
	outPutsNode := findNode(mainNode, "Output")
	if outPutsNode == nil {
		return nil
	}
	outputListNode := outPutsNode.Parent.NextSibling
	outputs := make([]Output, 0)
	elements := make(map[string]string, 0)
	outputNode := outputListNode.FirstChild
	for outputNode != nil {
		dataNode := outputNode.FirstChild.NextSibling.FirstChild
		for dataNode != nil {
			k := dataNode.FirstChild.Data
			dataNode = dataNode.NextSibling
			v := dataNode.FirstChild.Data
			dataNode = dataNode.NextSibling
			elements[k] = v
		}
		outputs = append(outputs, Output{
			Id:      outputNode.FirstChild.FirstChild.Data,
			Value:   elements["value"],
			Address: elements["address"],
		})
		outputNode = outputNode.NextSibling
	}
	//
	return &Transaction{
		Hash:    hash,
		Inputs:  inputs,
		Outputs: outputs,
	}
}

func findNode(node *html.Node, data string) *html.Node {
	child := node.FirstChild
	for child != nil {
		if strings.Contains(child.Data, data) {
			return child
		}
		found := findNode(child, data)
		if found != nil {
			return found
		}
		child = child.NextSibling
	}
	return nil
}
*/
