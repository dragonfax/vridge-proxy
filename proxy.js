var lodash = require('lodash');
var udpProxy = require('udp-proxy');
var tcpProxy = require("node-tcp-proxy");

var UDP_PORTS = [18481, 38217];
var TCP_PORTS_LOW = 38216;
var TCP_PORTS_HIGH = 38230;
var TCP_PORTS = lodash.range(TCP_PORTS_LOW, TCP_PORTS_HIGH + 1);
console.log("tcp ports " + TCP_PORTS);

var vridge_host = '64.62.255.87';
var local_host = '0.0.0.0';

var udp_servers = [];
UDP_PORTS.forEach(function (port) {

	var options = {
		address: vridge_host,
		port: port,
		localaddress: local_host,
		localport: port
	};

	var server = udpProxy.createServer(options);
	udp_servers.push(server);

	// 'message' is emitted when the server gets a message
	/*
	server.on('message', function (message, sender) {
		console.log('message from ' + sender.address + ':' + sender.port);
	});
	*/
});


var tcp_servers = [];
TCP_PORTS.forEach(function(port) {

	var server = tcpProxy.createProxy(port, vridge_host, port, { hostname: '0.0.0.0'});
	tcp_servers.push(server);
});

