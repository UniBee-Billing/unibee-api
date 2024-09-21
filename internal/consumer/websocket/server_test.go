package websocket

//func TestMerchantWebSocketMessageEntry(t *testing.T) {
//	client := gclient.NewWebSocket()
//	client.HandshakeTimeout = time.Second
//	client.Proxy = http.ProxyFromEnvironment
//	client.TLSClientConfig = &tls.Config{}
//
//	conn, _, err := client.Dial("ws://127.0.0.1:8088/merchant_ws/"+"EUXAgwv3Vcr1PFWt2SgBumMHXn3ImBqM", nil)
//	if err != nil {
//		panic(err)
//	}
//	defer func(conn *websocket.Conn) {
//		err := conn.Close()
//		if err != nil {
//
//		}
//	}(conn)
//
//	for {
//		mt, data, err := conn.ReadMessage()
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(mt, string(data))
//	}
//}
