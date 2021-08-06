import { Injectable, OnInit, OnDestroy } from '@angular/core';

import { HttpClient } from '@angular/common/http';
import { Observable, Subject } from 'rxjs';
import { AuthenticationData, LoginService, HeaderType } from '../login-service/login.service';
import { shareReplay } from 'rxjs/operators';

type Header = { [header: string]: string | string[] | null; };

export interface JoinRequest {
  roomName: string
}

export interface JoinResponse {
  roomName: string
  success: boolean
  previousMessages: string[]
}

export interface CreateRequest {
  roomName: string
  name: string
}

export interface CreateResponse {
  roomId: string
  roomName: string
  success: boolean
  previousMessages: string[]
}

@Injectable({
  providedIn: 'root'
})
export class RoomBackendService implements OnInit, OnDestroy {

  constructor(private http: HttpClient, private loginService: LoginService) {

  }
  ngOnDestroy() {

  }

  ngOnInit() {
  }

  create(roomName:string): Observable<CreateResponse> | undefined {
    const authData:AuthenticationData | undefined = this.loginService.authenticationData();
    const authHeader: HeaderType | undefined = this.loginService.authentictionHeader();

    if (authData !== undefined) {
      const params:CreateRequest = {
        roomName,
        name: authData.userName,
      };
      const createHeaders : Header = {};
      const auth: string = authHeader!.Authorization as string;

      const headers = {...createHeaders, ...this.headers, ...authHeader};
      return this.http.post<CreateResponse>(`${this.url}/create`, params, {headers: {'Authorization': auth}})
        .pipe(shareReplay());
    }
    return undefined;
  }

  join(roomName:string): Observable<JoinResponse> | undefined {
    const authHeader: HeaderType | undefined = this.loginService.authentictionHeader();
    if (authHeader !== undefined) {
      const params: JoinRequest = {
        roomName,
      };
      const joinHeaders : Header = { };
      const headers : Header = {...joinHeaders, ...this.headers, ...authHeader };
      const auth: string = authHeader!.Authorization as string;
      return this.http.post<JoinResponse>(`${this.url}/join`, params, {headers: {'Authorization': auth}})
        .pipe(shareReplay());
    }
    return undefined;
  }

  delete(roomId:string): Observable<Object> {
    const deleteHeaders : Header = {};
    return this.http.delete(`${this.url}`, {...deleteHeaders, ...this.headers});
  }
  
  private readonly url = 'http://localhost:8080/room';
  private headers : Header = {};
}
