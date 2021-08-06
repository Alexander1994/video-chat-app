import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ChatSideBarComponent } from './components/chat-side-bar/chat-side-bar.component';
import { LoginLandingComponent } from './components/login-landing/login-landing.component';
import { LobbyCreatorComponent } from './components/lobby-creator/lobby-creator.component';
import { RoomFinderComponent } from './components/room-finder/room-finder.component';
import { MessageRoomComponent } from './components/message-room/message-room.component';
import { SocketIoModule, SocketIoConfig } from 'ngx-socket-io';

const config: SocketIoConfig = { url: 'http://localhost:8988', options: {} };


@NgModule({
  declarations: [
    AppComponent,
    ChatSideBarComponent,
    LoginLandingComponent,
    LobbyCreatorComponent,
    RoomFinderComponent,
    MessageRoomComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    AppRoutingModule,
    ReactiveFormsModule,
    FormsModule,
    SocketIoModule.forRoot(config)
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
