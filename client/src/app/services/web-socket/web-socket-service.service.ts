import { Injectable } from '@angular/core';
import { Subject, Observer, Observable }  from "rxjs";
import { AnonymousSubject } from 'rxjs/internal/Subject';
import { LoginService } from '../login-service/login.service';

import { filter, map } from 'rxjs/operators'


export interface BasicMessage<T> {
  type: string
  data: T
}

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {

  private connectionObserver: Observable<Event> | null = null;

  private ws:WebSocket | null = null;

  constructor(private loginService:LoginService) {
    console.log("WebSocketService")

   }

  public connect(url:string) {
    if (!this.ws) {
      this.ws = new WebSocket(url);

      this.connectionObserver = this.createConnectedSubject(this.ws);
      
      console.log("Successfully connected: " + url);
    }
  }

  public onConnect(): Observable<Event> {
    return this.connectionObserver!;
  }

  public fromEvent<T>(msgType:string): Subject<T> {
    return this.createMessageObserver<T>(msgType, this.ws!);
  }

  private createConnectedSubject(ws:WebSocket):Observable<Event> {
    return new Observable((obs: Observer<Event>) => {
      ws.onopen = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);
      return ws.close.bind(ws);
    });
  }
  
  private createMessageObserver<T>(msgType:string, ws:WebSocket): Subject<T> {
    let observable = new Observable((obs: Observer<MessageEvent<any>>) => {
      ws.onmessage =  obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs)
      ws.onclose = obs.complete.bind(obs);
      return ws.close.bind(ws);
    }).pipe(
      
      map((val:MessageEvent<any>) => { console.log(val, msgType); return <BasicMessage<T>>JSON.parse(val.data)}),
      filter((msg:BasicMessage<T>) =>  msg.type === msgType ),
      map((data:BasicMessage<T>)=> data.data)
    );
    let observer = {
      next: (data: T) => {
        const basicMsg:BasicMessage<T> = {
          type: msgType,
          data
        }
        console.log(basicMsg);
        console.log(ws.readyState === WebSocket.OPEN);

        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify(basicMsg));
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