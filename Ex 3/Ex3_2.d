import std.stdio;
import std.socket;
import core.thread;


void main() {
  
  int ServerPORT = 34933;
  int localPORT = 50007;
  string serverIp = "129.241.187.161"
  
  auto serverAddr = new InternetAddress(serverIp, ServerPORT);
  auto sock = new TcpSocket(serverAddr);
  char [1024] buffer;
  sock.setOption(SocketOptionLevel.SOCKET, SocketOption.SNDBUF,1024);
  sock.setOption(SocketOptionLevel.SOCKET, SocketOption.RCVBUF,1024);
  scope(exit) sock.close();
  
  writeln("Connecting...");
  sock.receive(buffer);
  writeln(buffer);
  
  char message[34] = "Connect to: 129.241.187.155:50007\n"; 
  
  auto localAddr = new InternetAddress(localPORT)
  Socket acceptSock = new TcpSocket(localAddr);
  acceptSock.setOption(SocketOptionLevel.SOCKET, SocketOption.REUSEADDR, true);
  acceptSock.bind(localAddr);
  
  newSock = acceptSock.listen(3);
  
  //Threads for handeling send()/recv() on newSock
  
  
  
  

    Socket server = acceptSock.accept();
    auto received = server.receive(buffer);
    writeln("The server said:\n%s", buffer[0.. received]);
}
  

