import { Injectable, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { shareReplay } from 'rxjs/operators';

type Header = { [header: string]: string | string[]; };

export interface LoginResponse {
  name: string
  token: string
  success: boolean
}

export interface LoginRequest {
  name: string
}

@Injectable({
  providedIn: 'root'
})
export class LoginBackendService implements OnInit  {

  constructor(private http: HttpClient) {
  }

  ngOnInit() {
    
  }

  login(name: string): Observable<LoginResponse> {
    const loginHeaders: Header = {};
    const headers = {...loginHeaders, ...this.headers};
    const params:LoginRequest = {name};
    return this.http.post<LoginResponse>(`${this.url}/login`, params, headers)
      .pipe(shareReplay());
  }

  logout(): Observable<Object>  {
    const logoutHeaders : Header = {};
    return this.http.post(`${this.url}/logout`, {...logoutHeaders, ...this.headers})
      .pipe(shareReplay());

  }

  private readonly url = 'http://localhost:8080';
  private headers : Header = {}
}
