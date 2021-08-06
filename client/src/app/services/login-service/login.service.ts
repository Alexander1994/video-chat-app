import { Injectable, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { LoginResponse, LoginBackendService } from 'src/app/services/loginbackend-service/loginbackend.service'

export interface AuthenticationData {
  userToken:string;
  userName:string;
};

export type HeaderType = {
  [key: string]: string | null;
}

@Injectable({
  providedIn: 'root'
})
export class LoginService implements OnInit  {

  constructor(private loginbackend: LoginBackendService) {
  }

  authentictionHeader() : HeaderType | undefined {
    return (this.authData !== undefined) ? {"Authorization": "Bearer " +  this.authData!.userToken} : undefined;
  }

  authenticationData() : AuthenticationData | undefined {
    return this.authData;
  }

  ngOnInit() {
    
  }

  login(name: string): Observable<LoginResponse> {
    const obs = this.loginbackend.login(name);
    obs.subscribe({
      next: (resp: LoginResponse)=> {
        this.authData = { userName: resp.name, userToken: resp.token };
      }
    });
    return obs;
  }

  logout(): Observable<Object>  {
    const obs = this.loginbackend.logout();
    obs.subscribe({
      next: (resp: Object)=> {
        
      }
    });
    return obs;

  }
  
  private authData: AuthenticationData | undefined = undefined;
}
