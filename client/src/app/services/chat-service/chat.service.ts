import { Injectable, OnInit } from '@angular/core';
import { Observable, Subject } from "rxjs"
import { share } from 'rxjs/operators';
import { LoginService } from '../login-service/login.service';
import { WebSocketService } from '../web-socket/web-socket-service.service';

//const URL:string = "ws://127.0.0.1:8080/room/messages";
//
//type MessageType = "leave" | "error" | "message";

export interface Message {
  message: string;
  name?: string; // name not used in out going events
  date?: number; // date not used in out going events
}

//type SocketConfig = Partial<ManagerOptions & SocketOptions>;
@Injectable({
  providedIn: 'root'
})
export class ChatService implements OnInit {

  private readonly url: string = 'ws://localhost:8000/messages';

  private messageSubject:Subject<Message> | null = null;
  private leaveSubject:Subject<Message> | null = null;

  constructor(private websocketService: WebSocketService, private loginService:LoginService) {
    console.log("ChatService");
  }

  public connect() {
    const authHeaders:Record<string, string> = <Record<string, string>>this.loginService.authentictionHeader()!
    const reqHeader:string = '?' + (new URLSearchParams(authHeaders).toString());
    this.websocketService.connect(this.url + reqHeader);
  }

  public leave(): Subject<Message> {
    if (this.leaveSubject === null) this.leaveSubject = this.websocketService.fromEvent<Message>('leave');
    return this.leaveSubject!;
  }

  public message(): Subject<Message> {
    if (this.messageSubject === null) this.messageSubject = this.websocketService.fromEvent<Message>('message');
    return this.messageSubject!;
  }

  public fromConnection(): Observable<Object> {
    return this.websocketService.onConnect();
  }
  //public fromDisconnection(): Observable<Object> {
  //  //return this.websocketService.fromEvent<Object>('disconnect');
  //}
  //public message(msg:Message) {
  //  this.websocketService.emit('message', msg);
  //}
  //
  //public leave(msg:Message) {
  //  this.websocketService.emit('leave', msg);
  //}
  
  ngOnInit() {
  }
}

//class IoSocket {
//  private subscribersCounter: Record<string, number> = {};
//  private eventObservables$: Record<string, Observable<any>> = {};
//  private socket:Socket
//
//  constructor(url:string | undefined, options:SocketConfig) {
//    //const ioFunc = (io as any).default ? (io as any).default : io;
//    this.socket = (<any>io)(url, options);
//  }
//
//  public connect() {
//    this.socket.connect();
//  }
//
//  public emit(eventName:string, data:any) {
//    this.socket.emit(eventName, data);
//  }
//
//  public fromEvent<T>(eventName:string): Observable<T> {
//    if (!this.subscribersCounter[eventName]) {
//      this.subscribersCounter[eventName] = 0;
//    }
//    this.subscribersCounter[eventName]++;
////
//    if (!this.eventObservables$[eventName]) {
//      this.eventObservables$[eventName] = new Observable((observer: any) => {
//        const listener = (data: T) => {
//          observer.next(data);
//        };
//        this.socket.on(eventName, listener);
//        return () => {
//          this.subscribersCounter[eventName]--;
//          if (this.subscribersCounter[eventName] === 0) {
//            this.socket.off(eventName, listener);
//            delete this.eventObservables$[eventName];
//          }
//        };
//      }).pipe(share());
//    }
//    return this.eventObservables$[eventName];
//  }
//}