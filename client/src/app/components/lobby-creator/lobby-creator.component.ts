import { Component, OnInit, EventEmitter, Output } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { RoomService, } from 'src/app/services/room-service/room.service';

class LobbyInfo {
  roomId:string
  name:string;

  private static currentRoomId:number = 0;

  constructor(name:string) {
    this.name = name;
    this.roomId = LobbyInfo.generateRoomId();
  }
  
  private static generateRoomId():string {
    LobbyInfo.currentRoomId++;
    return LobbyInfo.currentRoomId.toString();
  }
}

@Component({
  selector: 'app-lobby-creator',
  templateUrl: './lobby-creator.component.html',
  styleUrls: ['./lobby-creator.component.less']
})
export class LobbyCreatorComponent implements OnInit {
  lobbyCreator = this.formBuilder.group({
    name: ['', Validators.required],
  });

  @Output() login = new EventEmitter<LobbyInfo>();

  constructor(private formBuilder: FormBuilder, private roomService: RoomService) { }
  
  onSubmit() {
    const name = this.lobbyCreator.value.name;
    this.roomService.create(name);

    this.lobbyCreator.reset();
  }

  ngOnInit(): void {
  }


}
