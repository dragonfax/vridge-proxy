var lodash = require('lodash');
var proxy = require('udp-proxy');

var UDP_PORTS = [18481, 38217];
var TCP_PORTS_LOW = 38216;
var TCP_PORTS_HIGH = 38230;
var TCP_PORTS = lodash.range(TCP_PORTS_LOW, TCP_PORTS_HIGH + 1);

var options = {
	address: '127.0.0.1',
	port: 53,
	ipv6: true,
	localaddress: '0.0.0.0',
	localport: 53535,
	localipv6: false,
	proxyaddress: '::0',
	timeOutTime: 10000
};

var server = proxy.createServer(options);

// 'message' is emitted when the server gets a message
server.on('message', function (message, sender) {
	console.log('message from ' + sender.address + ':' + sender.port);
});

