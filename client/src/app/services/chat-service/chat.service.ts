import { Injectable, OnInit } from '@angular/core';
import { Observable } from "rxjs"
import { share } from 'rxjs/operators';
import { LoginService } from '../login-service/login.service';
//import { WebSocketService } from '../web-socket/web-socket-service.service';
import { io, Socket, ManagerOptions, SocketOptions } from "socket.io-client";

const URL:string = "ws://127.0.0.1:8080/room/messages";

type MessageType = "leave" | "error" | "message";

export interface Message {
  message: string;
  name?: string; // name not used in out going events
  date?: number; // date not used in out going events
}

type SocketConfig = Partial<ManagerOptions & SocketOptions>;

@Injectable({
  providedIn: 'root'
})
export class ChatService implements OnInit {

  private readonly url: string = 'http://localhost:8080/socket.io/';
  private socket:IoSocket

  constructor(private loginService: LoginService) {
    //this.loginService.authentictionHeader()
    const Authorization = this.loginService.authentictionHeader()!.Authorization!;
    const cfg:SocketConfig = {
      withCredentials: true,
      extraHeaders: {Authorization}
    }; 
    this.socket = new IoSocket(this.url, {});
    console.log("ChatService");
    this.socket.connect();//
  }

  public connect() {
    //return this.socket.connect();
  }

  public fromMessage(): Observable<Message> {
    return this.socket.fromEvent<Message>('message');
  }
  public fromLeave(): Observable<Message> {
    return this.socket.fromEvent<Message>('message');
  }
  public fromConnection(): Observable<Object> {
    return this.socket.fromEvent<Object>('connection');
  }
  public fromDisconnection(): Observable<Object> {
    return this.socket.fromEvent<Object>('disconnect');
  }
  public message(msg:Message) {
    this.socket.emit('message', msg);
  }
  
  public leave(msg:Message) {
    this.socket.emit('leave', msg);
  }
  
  ngOnInit() {
  }
}

class IoSocket {
  private subscribersCounter: Record<string, number> = {};
  private eventObservables$: Record<string, Observable<any>> = {};
  private socket:Socket

  constructor(url:string, options:SocketConfig) {
    //const ioFunc = (io as any).default ? (io as any).default : io;
    this.socket = io(url, options);
  }

  public connect() {
    this.socket.connect();
  }

  public emit(eventName:string, data:any) {
    this.socket.emit(eventName, data);
  }

  public fromEvent<T>(eventName:string): Observable<T> {
    if (!this.subscribersCounter[eventName]) {
      this.subscribersCounter[eventName] = 0;
    }
    this.subscribersCounter[eventName]++;

    if (!this.eventObservables$[eventName]) {
      this.eventObservables$[eventName] = new Observable((observer: any) => {
        const listener = (data: T) => {
          observer.next(data);
        };
        this.socket.on(eventName, listener);
        return () => {
          this.subscribersCounter[eventName]--;
          if (this.subscribersCounter[eventName] === 0) {
            this.socket.off(eventName, listener);
            delete this.eventObservables$[eventName];
          }
        };
      }).pipe(share());
    }
    return this.eventObservables$[eventName];
  }
}