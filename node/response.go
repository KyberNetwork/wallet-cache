package node

// {
//     "jsonrpc": "2.0",
//     "id": 5,
//     "error": {
//         "code": -32000,
//         "message": "Origin \"None\" is not on whitelist."
//     }
// }

// {
//     "jsonrpc": "2.0",
//     "result": "0x296",
//     "id": 5
// }

type RequestRPC struct {
	Method string `json:"method"`
}
