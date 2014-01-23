import std.stdio;
import std.getopt;
import std.string;
import std.socket;
import std.stream;
import std.socketstream;

void main(string[] args) {
	debug(TCP)
		writefln("Debuging enabled");
	string ip = "127.0.0.1";
	ushort port = 124;
	getopt(args, "host|h", &ip, "port|p", &port);
	writefln("ip: %s port: %s", ip, port);
	tcp_connect(ip, port);
}



void tcp_connect(string ip, ushort port) {
	debug(TCP)
		writefln("Connecting to host %s:%s", ip, port);
	Socket sock = new TcpSocket(new InternetAddress(ip, port));
	scope(exit) sock.close();
	
	debug(TCP)
		writefln("reading line");

	char[1024] line;
	sock.receive(line);
	writefln("string: %s", line);
	
	for(int i=0; i < 5; i++) {
		char[] message;
		writefln("Write message to send:");
		readln(message);
		message[message.length-1] = '\0';
		debug(TCP)
			writefln("sending message");
		sock.send(message);
		debug(TCP)
			writefln("recivev\n");
		sock.receive(line);
		writefln("Recived line: %s", line);
	}
}
	

