import std.stdio;
import std.socket;
import std.c.time;


void main() {
  
  const int PORT = 34933;
  const int localPort = 50007;
  
  auto addr = new parseAddress("129.241.187.161", PORT);
  auto sock = new TcpSocket();
  char [1024] buffer;
  sock.setOption(SocketOptionLevel.SOCKET, SocketOption.SNDBUF,1024);
  sock.setOption(SocketOptionLevel.SOCKET, SocketOption.RCVBUF,1024);
  scope(exit) acceptSock.close();
  
  writeln("Connecting...");
  sock.receive(buffer);
  writeln(buffer);
  
  char message[34] = "Connect to: 129.241.187.155:50007\n"; 
  
  
  
/*  acceptSock.bind(addr);
  acceptSock.listen(5);
  
  
  
  
  while(true){
    acceptSock.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, true);
    Socket server = acceptSock.accept();
    char[1024] buffer;
    auto received = server.receive(buffer);
    
    writeln("The server said:\n%s", buffer[0..received];
*/  
}