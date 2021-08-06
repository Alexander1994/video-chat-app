import { Component, OnInit } from '@angular/core';
import { ChatService, Message } from 'src/app/services/chat-service/chat.service';
import { WebSocketService } from 'src/app/services/web-socket/web-socket-service.service';
import { FormBuilder } from '@angular/forms';
import { LoginService } from 'src/app/services/login-service/login.service';
import { shareReplay } from 'rxjs/operators';

@Component({
  selector: 'app-chat-side-bar',
  templateUrl: './chat-side-bar.component.html',
  styleUrls: ['./chat-side-bar.component.less'],
  providers: [ WebSocketService, ChatService ]
})

export class ChatSideBarComponent implements OnInit {

  messages: Array<Message> = [];
  msgForm = this.formBuilder.group({
    message: ['']
  });

  constructor(private formBuilder: FormBuilder, private chatService: ChatService, private loginService:LoginService) {
    console.log("ChatSideBarComponent");
    this.chatService.fromMessage().subscribe((msg)=> {
      this.messages.push(msg);
    });
    //this.chatService.fromConnection().subscribe((o:Object) => {
    //  let msg :Message = {
    //    message: this.loginService.authenticationData()!.userToken
    //  };
    //  this.chatService.messages!.next(msg);
    //});
  }

  
  timeConverter(unixTimeStamp: number) :string {
    const a = new Date(unixTimeStamp);
    const months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
    const year = a.getFullYear();
    const month = months[a.getMonth()];
    const date = a.getDate();shareReplay
    const hour = a.getHours();
    const min = a.getMinutes() < 10 ? '0' + a.getMinutes() : a.getMinutes();
    const sec = a.getSeconds() < 10 ? '0' + a.getSeconds() : a.getSeconds();
    const time = date + ' ' + month + ' ' + year + ' ' + hour + ':' + min + ':' + sec ;
    return time;
  }
    
  ngOnInit(): void {

  }

  onSend() {
    const authData = this.loginService.authenticationData();
    if (authData) {
      let msg: Message = {
        message: this.msgForm.value.message
      };
      this.chatService.message(msg);
      this.msgForm.reset();
    } else {
      console.log("unauthenticated user attempting to send data");
    }
  }


}
