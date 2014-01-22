import std.stdio
import std.socket

int main() {
  
  const int PORT = 34933;
  
  InternetAddress addr = new InternetAddress(/*server ip*/, PORT);
  TcpSocket sock = new socket(AdressFamily.INET, SocketType.STREAM, ProtocolType.TCP);
  
  sock.connect(addr); //connectTCP() ??



}
