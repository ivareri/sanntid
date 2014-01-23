import std.stdio;
import std.socket;
import std.string;

void main() {
  
  const int PORT = 34933;
  string ip = "129.241.187.161";
  
  InternetAddress addr = new InternetAddress(ip, PORT);
  Socket sock = new TcpSocket(addr);
  scope(exit) sock.close();
  
  
  
  char [1024] buffer;
  auto recv = sock.receive(buffer);
  writefln("received: %s", buffer);
  
  for(int i = 0; i < 3; i++){
    char [] message;
    writeln("Write message:");
    readln(message);
    message[message.length-1] = '\0';
    sock.send(message);
    sock.receive(buffer);
    writeln(buffer,"\n");
    
  }
}
