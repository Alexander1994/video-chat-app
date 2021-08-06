import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { RoomService } from '../../services/room-service/room.service';

@Injectable({
    providedIn: 'root'
})
export class RoomAuthGuard implements CanActivate {

    constructor(private roomService: RoomService, private router:Router) {

    }

    canActivate(next: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
        this.authenticated = this.roomService.roomAuthenticationData() !== undefined;
        if (!this.authenticated) {
            this.router.navigate(['/']);
        }
        console.log("RoomAuthGuard canActivate");

        return this.authenticated;
    }

    private authenticated:boolean = false;
}