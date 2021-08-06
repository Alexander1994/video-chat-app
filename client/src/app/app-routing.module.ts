import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginLandingComponent } from "./components/login-landing/login-landing.component";

// Authentication
import { LoginAuthGuard } from './authentication/login/login.guard'
import { RoomAuthGuard } from './authentication/room/room.guard';

import { RoomFinderComponent } from "./components/room-finder/room-finder.component";
import { MessageRoomComponent } from "./components/message-room/message-room.component"

const routes: Routes = [
  { path: '', component: LoginLandingComponent },
  { path: 'main', component: RoomFinderComponent, canActivate: [LoginAuthGuard]},
  { path: 'room', component: MessageRoomComponent, canActivate: [LoginAuthGuard, RoomAuthGuard] }
];


@NgModule({
  imports: [
    RouterModule.forRoot(routes)
  ],
  exports: [ RouterModule ]
})
export class AppRoutingModule { }
