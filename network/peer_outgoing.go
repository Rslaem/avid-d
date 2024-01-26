package network

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type outgoingPeer struct {
	ourID      int
	ourAddr    string
	num        int
	peers      []peer
	mutex      *sync.RWMutex
	ready      bool
	totalBytes int
	maxBytes   int
}

func NewOutgoingPeer(id, num int, addr string, peers []peer, mutex *sync.RWMutex) *outgoingPeer {
	p := &outgoingPeer{
		ourID:      id,
		ourAddr:    addr,
		num:        num,
		peers:      peers,
		mutex:      mutex,
		ready:      false,
		totalBytes: 0,
		maxBytes:   0,
	}
	return p
}

func (p *outgoingPeer) init() {
	var wg sync.WaitGroup
	wg.Add(p.num - 1)
	for _, peer := range p.peers {
		if p.ourID == peer.Id {
			continue
		}
		go func(ip string, id int) {
			apiUrl := fmt.Sprintf("http://%s/ready", ip)
			data := url.Values{}
			data.Set("id", fmt.Sprint(p.ourID))
			data.Set("addr", p.ourAddr)
			u, err := url.ParseRequestURI(apiUrl)
			if err != nil {
				fmt.Printf("parse url requestUrl failed,err:%v\n", err)
			}
			u.RawQuery = data.Encode() // URL encode
			for {
				//log.Printf("[node %d] trying connect to node %d on %s", p.ourID, id, ip)
				_, err = http.Get(u.String())
				if err != nil {
					//log.Printf("[node %d] failed connect to node %d on %s\n", p.ourID, id, ip)
					data = nil
					time.Sleep(3 * time.Second)
					continue
				} else {
					data = nil
					//log.Printf("[node %d] connect to node %d on %s: %s\n", p.ourID, id, ip, resp.Status)
					wg.Done()
					break
				}
			}
		}(peer.Addr, peer.Id)
	}
	wg.Wait()
	p.ready = true
	time.Sleep(1 * time.Second)
	log.Printf("[node %d] dial ready\n", p.ourID)
}

func (p *outgoingPeer) SendPost(id int, dataType, api string, data []byte) ([]byte, error) {
	var ip string
	Client := &http.Client{
		Timeout: 30 * time.Second,
	}
	for _, peer := range p.peers {
		if id == peer.Id {
			ip = peer.Addr
			break
		}
	}
	apiUrl := fmt.Sprintf("http://%s/%s", ip, api)

	p.mutex.Lock()
	if len(data) > p.maxBytes {
		p.maxBytes = len(data)
	}
	p.totalBytes += len(data)
	p.mutex.Unlock()

	postData := url.Values{}
	postData.Add("From", fmt.Sprintf("%d", p.ourID))
	postData.Add("Lenth", fmt.Sprintf("%d", len(data)))
	postData.Add("Type", dataType)
	postData.Add("Content", string(data))
	//log.Printf("send message %v", postData)
	reader := strings.NewReader(postData.Encode())
	request, err := http.NewRequest("POST", apiUrl, reader)
	request.Close = true
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Idempotency-Key", fmt.Sprintf("%d%s%d", p.ourID, dataType, id))
	resp, err := Client.Do(request)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return respBytes, nil
}
