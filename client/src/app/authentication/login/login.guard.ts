import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { LoginService } from '../../services/login-service/login.service';

@Injectable({
    providedIn: 'root'
})
export class LoginAuthGuard implements CanActivate {

    constructor(private loginService: LoginService, private router:Router) {

    }

    canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
        this.authenticated = this.loginService.authenticationData() !== undefined;
        if (!this.authenticated) {
            this.router.navigate(['/']);
        }
        return this.authenticated;
    }

    private authenticated:boolean = false;
}