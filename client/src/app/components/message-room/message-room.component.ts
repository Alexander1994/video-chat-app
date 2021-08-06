import { Component, OnInit } from '@angular/core';
import { RoomService } from 'src/app/services/room-service/room.service';

@Component({
  selector: 'app-message-room',
  templateUrl: './message-room.component.html',
  styleUrls: ['./message-room.component.less']
})
export class MessageRoomComponent implements OnInit {

  roomName:string = "";

  constructor(private roomService:RoomService) {

  }

  onExit(): void {
    // roomExit
  }

  ngOnInit(): void {

    this.roomName = this.roomService.roomAuthenticationData()!.roomName;

  }

}
