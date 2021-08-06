import { Injectable } from '@angular/core';
import { Subject, Observer, Observable }  from "rxjs";
import { AnonymousSubject } from 'rxjs/internal/Subject';

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {
  private webSocketSubject: Subject<MessageEvent> | null = null;
  private connectionObserver: Observable<Event> | null = null;

  constructor() {
    console.log("WebSocketService")

   }

  public connect(url:string): [Subject<MessageEvent>, Observable<Event>] {
    if (!this.webSocketSubject) {
      let ws = new WebSocket(url);

      this.webSocketSubject = this.createMessageObserver(ws);
      this.connectionObserver = this.createConnectedSubject(ws);
      
      console.log("Successfully connected: " + url);
    }
    return [this.webSocketSubject!, this.connectionObserver!];
  }

  private createConnectedSubject(ws:WebSocket):Observable<Event> {
    return new Observable((obs: Observer<Event>) => {
      ws.onopen = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);
      return ws.close.bind(ws);
    });
  }
  
  private createMessageObserver(ws:WebSocket): Subject<MessageEvent> {
    let observable = new Observable((obs: Observer<MessageEvent>) => {
      ws.onmessage = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);
      return ws.close.bind(ws);
    });
    let observer = {
      next: (data: Object) => {
        console.log(data);
        console.log(ws.readyState === WebSocket.OPEN);

        if (ws.readyState === WebSocket.OPEN) {

          ws.send(JSON.stringify(data));
        }
      },
      error: (err: any) => {},
      complete: () => {}
    };
    return new AnonymousSubject(observer, observable);
  }
}

export enum MessageType {
  Join = "JOIN",
  CreateRoom = "CREATE_ROOM",
  Message = "MESSAGE",
}